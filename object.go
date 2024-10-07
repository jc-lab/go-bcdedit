package go_bcdedit

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/gabriel-samfira/go-hivex"
	"github.com/jc-lab/go-bcdedit/model"
	"github.com/jc-lab/go-bcdedit/pkg/hiveutil"
	"slices"
	"strings"
	"unicode/utf16"
)

type BcdElement interface {
	GetType() ValueType
	GetRaw() []byte
}

type BcdObject interface {
	GetId() string
	GetDescription() model.BcdDescription
	GetElements() map[string]BcdElement
	SortedElements() []BcdElement
	SetElement(key string, typ ValueType, raw []byte) (BcdElement, error)
	ToJson() *model.BcdObject
}

// KnownObjectIds BCD.docx, page 9: Standard application Objects
var KnownObjectIds = map[string]string{
	"{9dea862c-5cdd-4e70-acc1-f32b344d4795}": "{bootmgr}",   // 0x10100002
	"{a5a30fa2-3d06-4e9f-b5f4-a01df9d1fcba}": "{fwbootmgr}", // 0x10100001
	"{b2721d73-1db4-4c62-bf78-c548a880142d}": "{memdiag}",
	"{466f5a88-0af2-4f76-9038-095b170dc21c}": "{ntldr}", // 0x10300006
	"{fa926493-6f1c-4193-a414-58f0b2456d1e}": "{current}",
	"{5189b25c-5558-4bf2-bca4-289b11bd29e2}": "{badmemory}",
	"{6efb52bf-1766-41db-a6b3-0ee5eff72bd7}": "{bootloadersettings}",
	"{4636856e-540f-4170-a130-a84776f4c654}": "{dbgsettings}",
	"{0ce4991b-e6b3-4b16-b23c-5e0d9250e5d9}": "{emssettings}",
	"{7ea2e1ac-2e61-4728-aaa3-896d9d0a9f0e}": "{globalsettings}",
	"{1afa9c49-16ab-4a5c-901b-212802da9460}": "{resumeloadersettings}",
}

func (e *HiveBcdElement) Meta() *model.BcdElementMeta {
	var elementTypes map[string]*model.BcdElementMeta
	switch e.Parent.Description.ObjectType() {
	case model.ObjectApplication:
		elementTypes = model.BcdApplicationElementTypes[e.Parent.Description.ApplicationType()]
	case model.ObjectDevice:
		elementTypes = model.BcdDeviceElementTypes
	case model.ObjectInherit:
		switch e.Parent.Description.ObjectSubType() {
		case model.InheritableByApplicationObjects:
			elementTypes = model.BcdApplicationElementTypes[e.Parent.Description.ApplicationType()]
		case model.InheritableByDeviceObjects:
			elementTypes = model.BcdDeviceElementTypes
		}
	}
	if elementTypes != nil {
		return elementTypes[e.Key]
	}
	return nil
}

func (e *HiveBcdElement) Name() string {
	meta := e.Meta()
	if meta == nil {
		return ""
	}
	return meta.Name
}

func (e *HiveBcdElement) GetString() (string, error) {
	if e.Type != RegSz {
		return "", fmt.Errorf("no RegSz type: %d", e.Type)
	}
	_, s, err := Utf16LEToString(e.Raw)
	if err != nil {
		return "", err
	}
	return s, nil
}

func (e *HiveBcdElement) GetMultiStrings() ([]string, error) {
	if e.Type != RegMultiSz {
		return nil, fmt.Errorf("no RegMultiSz type: %d", e.Type)
	}
	var results []string
	remaining := e.Raw
	for len(remaining) > 0 {
		n, s, err := Utf16LEToString(remaining)
		if err != nil {
			return nil, err
		}
		if s == "" && n == 2 {
			break
		}
		results = append(results, s)
		remaining = remaining[n:]
	}
	return results, nil
}

func (e *HiveBcdElement) GetDword() (uint32, error) {
	if e.Type != RegDword {
		return 0, fmt.Errorf("no RegDword type: %d", e.Type)
	}
	return binary.LittleEndian.Uint32(e.Raw), nil
}

func (e *HiveBcdElement) String() string {
	var err error
	var results []string

	switch e.Type {
	case RegSz:
		s, err := e.GetString()
		if err != nil {
			return fmt.Sprintf("ERROR: %+v", err)
		}
		results = append(results, s)
	case RegMultiSz:
		results, err = e.GetMultiStrings()
		if err != nil {
			return fmt.Sprintf("ERROR: %+v", err)
		}
	default:
		return fmt.Sprintf("Type=%v, Raw=%s", e.Type, hex.EncodeToString(e.Raw))
	}
	return strings.Join(results, "\n")
}

func (o *HiveBcdObject) SortedElements() []BcdElement {
	var sortedElements []BcdElement
	for _, element := range o.Elements {
		sortedElements = append(sortedElements, element)
	}
	slices.SortFunc(sortedElements, func(a, b BcdElement) int {
		ha := a.(*HiveBcdElement)
		hb := b.(*HiveBcdElement)
		return strings.Compare(ha.Key, hb.Key)
	})
	return sortedElements
}

func (o *HiveBcdObject) SetElement(key string, typ ValueType, raw []byte) (BcdElement, error) {
	elementNode, err := hiveutil.UpsertNode(o.Bcd.Hive, o.ElementsNode, key)
	if err != nil {
		return nil, err
	}
	_, err = o.Bcd.Hive.NodeSetValue(elementNode, hivex.HiveValue{
		Type:  int(typ),
		Key:   "Element",
		Value: raw,
	})
	e := &HiveBcdElement{
		Parent: nil,
		Node:   elementNode,
		Key:    key,
		Type:   typ,
		Raw:    raw,
	}
	o.Elements[key] = e
	return e, nil
}

func Utf16LEToString(b []byte) (int, string, error) {
	if len(b)%2 != 0 {
		return 0, "", fmt.Errorf("invalid UTF-16 LE byte array length: %d", len(b))
	}

	u16 := make([]uint16, len(b)/2)
	err := binary.Read(bytes.NewReader(b), binary.LittleEndian, &u16)
	if err != nil {
		return 0, "", fmt.Errorf("failed to read UTF-16 LE bytes: %v", err)
	}

	n := 0
	for n = 0; n < len(u16); n++ {
		if u16[n] == 0 {
			u16 = u16[:n]
			n++
			break
		}
	}

	runes := utf16.Decode(u16)

	return n * 2, string(runes), nil
}

func stringToUtf16LE(buffer *bytes.Buffer, s string) error {
	utf16Encoded := utf16.Encode([]rune(s))
	return binary.Write(buffer, binary.LittleEndian, utf16Encoded)
}

func StringToUtf16LE(s string) ([]byte, error) {
	var buffer bytes.Buffer
	err := stringToUtf16LE(&buffer, s)
	if err != nil {
		return nil, err
	}
	return buffer.Bytes(), nil
}

func StringsToMultiUtf16LE(list []string) ([]byte, error) {
	var buffer bytes.Buffer
	for _, s := range list {
		err := stringToUtf16LE(&buffer, s)
		if err != nil {
			return nil, err
		}
		buffer.WriteByte(0)
		buffer.WriteByte(0)
	}
	buffer.WriteByte(0)
	buffer.WriteByte(0)
	return buffer.Bytes(), nil
}
