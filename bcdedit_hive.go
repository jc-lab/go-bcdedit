package go_bcdedit

import (
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"github.com/gabriel-samfira/go-hivex"
	"github.com/jc-lab/go-bcdedit/model"
	"github.com/jc-lab/go-bcdedit/pkg/hiveutil"
)

const (
	RegNone                     ValueType = hivex.RegNone
	RegSz                       ValueType = hivex.RegSz
	RegExpandSz                 ValueType = hivex.RegExpandSz
	RegBinary                   ValueType = hivex.RegBinary
	RegDword                    ValueType = hivex.RegDword
	RegDwordBigEndian           ValueType = hivex.RegDwordBigEndian
	RegLink                     ValueType = hivex.RegLink
	RegMultiSz                  ValueType = hivex.RegMultiSz
	RegResourceList             ValueType = hivex.RegResourceList
	RegFullResourceDescriptor   ValueType = hivex.RegFullResourceDescriptor
	RegResourceRequirementsList ValueType = hivex.RegResourceRequirementsList
	RegQword                    ValueType = hivex.RegQword
)

func (t ValueType) ToJson() model.ValueType {
	switch t {
	case RegNone:
		return model.RegNone
	case RegSz:
		return model.RegSz
	case RegExpandSz:
		return model.RegExpandSz
	case RegBinary:
		return model.RegBinary
	case RegDword:
		return model.RegDword
	case RegDwordBigEndian:
		return model.RegDwordBigEndian
	case RegLink:
		return model.RegLink
	case RegMultiSz:
		return model.RegMultiSz
	case RegResourceList:
		return model.RegResourceList
	case RegFullResourceDescriptor:
		return model.RegFullResourceDescriptor
	case RegResourceRequirementsList:
		return model.RegResourceRequirementsList
	case RegQword:
		return model.RegQword
	}
	return ""
}

func ValueTypeFromJson(t model.ValueType) ValueType {
	switch t {
	case model.RegNone:
		return RegNone
	case model.RegSz:
		return RegSz
	case model.RegExpandSz:
		return RegExpandSz
	case model.RegBinary:
		return RegBinary
	case model.RegDword:
		return RegDword
	case model.RegDwordBigEndian:
		return RegDwordBigEndian
	case model.RegLink:
		return RegLink
	case model.RegMultiSz:
		return RegMultiSz
	case model.RegResourceList:
		return RegResourceList
	case model.RegFullResourceDescriptor:
		return RegFullResourceDescriptor
	case model.RegResourceRequirementsList:
		return RegResourceRequirementsList
	case model.RegQword:
		return RegQword
	}
	return RegNone
}

type HiveBcdedit struct {
	Hive     *hivex.Hivex
	Writable bool
}

func NewWithHive(hive *hivex.Hivex, writable bool) (Bcdedit, error) {
	return &HiveBcdedit{
		Hive:     hive,
		Writable: writable,
	}, nil
}

func (b *HiveBcdedit) Close() error {
	var err error
	if b.Writable {
		_, err = b.Hive.Commit()
	}
	closeErr := b.Hive.Close()
	if err != nil {
		return err
	}
	return closeErr
}

func (b *HiveBcdedit) Enumerate(targetObjectId string) (map[string]BcdObject, error) {
	objectMap := map[string]BcdObject{}

	root, err := hiveutil.GetObjectsNode(b.Hive)
	if err != nil {
		return nil, err
	}
	err = hiveutil.ReadNode(b.Hive, root, func(objectNode int64, objectId string, err error) error {
		if err != nil {
			return err
		}
		var object BcdObject
		if len(targetObjectId) > 0 && targetObjectId != "all" {
			if objectId != targetObjectId {
				return nil
			}
		}
		object, err = b.getObject(objectId, objectNode)
		if err != nil {
			return err
		}
		objectMap[objectId] = object
		return nil
	})
	if err != nil {
		return nil, err
	}

	return objectMap, nil
}

func (b *HiveBcdedit) GetObject(objectId string) (BcdObject, error) {
	root, err := hiveutil.GetObjectsNode(b.Hive)
	if err != nil {
		return nil, err
	}
	objectNode, err := hiveutil.FindChild(b.Hive, root, objectId)
	if err != nil {
		return nil, err
	}
	if objectNode == 0 {
		return nil, fmt.Errorf("not exists %s", objectId)
	}
	return b.getObject(objectId, objectNode)
}

func (b *HiveBcdedit) getObject(objectId string, objectNode int64) (BcdObject, error) {
	descriptionNode, err := hiveutil.FindChild(b.Hive, objectNode, "Description")
	if err != nil {
		return nil, err
	}
	if descriptionNode == 0 {
		return nil, fmt.Errorf("not exists %s\\Description", objectId)
	}
	typeValue, err := hiveutil.FindValue(b.Hive, descriptionNode, "Type")
	if err != nil {
		return nil, err
	}
	if typeValue == 0 {
		return nil, fmt.Errorf("not exists %s\\Description\\Type", objectId)
	}
	valType, valueBytes, err := b.Hive.ValueValue(typeValue)
	if err != nil {
		return nil, fmt.Errorf("%s\\Description\\Type type read failed: %+v", objectId, err)
	}
	if valType != hivex.RegDword {
		return nil, fmt.Errorf("%s\\Description\\Type type is %d not dword", objectId, valType)
	}
	elementsNode, err := hiveutil.FindChild(b.Hive, objectNode, "Elements")
	if err != nil {
		return nil, err
	}
	object := &HiveBcdObject{
		Bcd:          b,
		Node:         objectNode,
		Id:           objectId,
		ElementsNode: elementsNode,
		Description:  model.BcdDescription(binary.LittleEndian.Uint32(valueBytes)),
		Elements:     make(map[string]*HiveBcdElement),
	}
	err = object.readElements(b)
	if err != nil {
		return nil, err
	}
	return object, nil
}

func (b *HiveBcdedit) UpsertObject(objectId string, description model.BcdDescription) (BcdObject, error) {
	root, err := hiveutil.GetObjectsNode(b.Hive)
	if err != nil {
		return nil, err
	}
	objectNode, err := hiveutil.UpsertNode(b.Hive, root, objectId)
	if err != nil {
		return nil, err
	}
	descriptionNode, err := hiveutil.UpsertNode(b.Hive, objectNode, "Description")
	if err != nil {
		return nil, err
	}
	_, err = b.Hive.NodeSetValue(descriptionNode, hivex.HiveValue{
		Type:  hivex.RegDword,
		Key:   "Type",
		Value: binary.LittleEndian.AppendUint32(nil, uint32(description)),
	})
	if err != nil {
		return nil, err
	}
	elementsNode, err := hiveutil.UpsertNode(b.Hive, objectNode, "Elements")
	if err != nil {
		return nil, err
	}

	object := &HiveBcdObject{
		Bcd:          b,
		Node:         objectNode,
		ElementsNode: elementsNode,
		Id:           objectId,
		Description:  description,
		Elements:     make(map[string]*HiveBcdElement),
	}
	return object, nil
}

func (b *HiveBcdedit) getElement(parent *HiveBcdObject, node int64, key string, value int64) (*HiveBcdElement, error) {
	element := &HiveBcdElement{
		Parent: parent,
		Node:   node,
		Key:    key,
	}
	valType, valueBytes, err := b.Hive.ValueValue(value)
	if err != nil {
		return nil, err
	}
	element.Type = ValueType(valType)
	element.Raw = valueBytes
	return element, nil
}

type HiveBcdElement struct {
	Parent *HiveBcdObject
	Node   int64
	Key    string
	Type   ValueType
	Raw    []byte
}

func (e *HiveBcdElement) GetType() ValueType {
	return e.Type
}

func (e *HiveBcdElement) GetRaw() []byte {
	return e.Raw
}

type HiveBcdObject struct {
	Bcd *HiveBcdedit

	Node         int64
	ElementsNode int64

	Id string // e.g. "{b2721d73-1db4-4c62-bf78-c548a880142d}"

	Description model.BcdDescription

	Elements map[string]*HiveBcdElement // e.g. key="11000001"
}

func (o *HiveBcdObject) readElements(bcd *HiveBcdedit) error {
	if o.ElementsNode == 0 {
		return nil
	}

	elements := make(map[string]*HiveBcdElement)
	err := hiveutil.ReadNode(bcd.Hive, o.ElementsNode, func(node int64, name string, err error) error {
		if err != nil {
			return err
		}
		value, err := hiveutil.FindValue(bcd.Hive, node, "Element")
		if err != nil {
			return err
		}
		element, err := bcd.getElement(o, node, name, value)
		if err != nil {
			return err
		}
		elements[name] = element
		return nil
	})
	if err != nil {
		return err
	}
	o.Elements = elements
	return nil
}

func (o *HiveBcdObject) GetId() string {
	return o.Id
}

func (o *HiveBcdObject) GetDescription() model.BcdDescription {
	return o.Description
}

func (o *HiveBcdObject) GetElements() map[string]BcdElement {
	fixedMap := make(map[string]BcdElement)
	for key, element := range o.Elements {
		fixedMap[key] = element
	}
	return fixedMap
}

func (o *HiveBcdObject) ToJson() *model.BcdObject {
	elements := make(map[string]*model.BcdElement)
	for key, element := range o.Elements {
		jsonElement := &model.BcdElement{
			Type: element.GetType().ToJson(),
			Raw:  base64.StdEncoding.EncodeToString(element.GetRaw()),
		}

		switch element.GetType() {
		case RegSz:
			jsonElement.ValueSz, _ = element.GetString()

		case RegMultiSz:
			jsonElement.ValueMultiSz, _ = element.GetMultiStrings()

		case RegDword:
			valueDword, err := element.GetDword()
			if err == nil {
				jsonElement.ValueDword = &valueDword
			}
		}

		elements[key] = jsonElement
	}
	return &model.BcdObject{
		Description: o.Description,
		Elements:    elements,
	}
}
