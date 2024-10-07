package bcdedit_cmd

import (
	"encoding/base64"
	"encoding/binary"
	"encoding/json"
	"flag"
	"fmt"
	go_bcdedit "github.com/jc-lab/go-bcdedit"
	"github.com/jc-lab/go-bcdedit/model"
	"github.com/pkg/errors"
	"log"
	"os"
	"slices"
	"strconv"
	"strings"
)

type ArrayFlags []string

// String is an implementation of the flag.Value interface
func (i *ArrayFlags) String() string {
	return fmt.Sprintf("%v", *i)
}

// Set is an implementation of the flag.Value interface
func (i *ArrayFlags) Set(value string) error {
	*i = append(*i, value)
	return nil
}

type Flags struct {
	Json        bool
	CreateStore string
	Store       string

	Enum string

	CreateId          string
	ObjectDescription string
	CreateDescription uint

	Set          bool
	SetId        string
	SetKey       string
	SetValueType string
	SetValueRaw  string
	SetValue     ArrayFlags
}

type commandDefine struct {
	Usage    string
	Writable int
	Runner   func(flags *Flags, args []string, bcd go_bcdedit.Bcdedit) error
}

var commands = map[string]commandDefine{
	"createstore": {
		Usage:    "/createstore <bcd_file>\nCreates a new and empty boot configuration data store.",
		Writable: -1,
		Runner: func(flags *Flags, args []string, bcd go_bcdedit.Bcdedit) error {
			flags.Store = args[0]
			return doCreateStore(flags)
		},
	},
	"enum": {
		Usage:    "/enum all\nThis command lists entries in a store.",
		Writable: 0,
		Runner: func(flags *Flags, args []string, bcd go_bcdedit.Bcdedit) error {
			flags.Enum = args[0]
			return doEnum(flags, bcd)
		},
	},

	// bcdedit /store BCD --create {4e2a0b91-5003-47c6-b00b-e496f2d66a11} -object-type 0x10200002
	"create": {
		Usage:    "/create <id> --object-type <object type(e.g. 0x10200002)> [/d <description>]\nThis command creates a new entry in the boot configuration data store.",
		Writable: 1,
		Runner: func(flags *Flags, args []string, bcd go_bcdedit.Bcdedit) error {
			flags.CreateId = args[0]

			subFlagset := flag.NewFlagSet("", flag.ExitOnError)
			subFlagset.UintVar(&flags.CreateDescription, "object-type", 0, "")
			subFlagset.StringVar(&flags.ObjectDescription, "d", "", "description")
			subFlagset.Parse(args[1:])

			return doCreateObject(flags, bcd)
		},
	},

	// bcdedit /store BCD /set {ObjectId} --value-type RegSz --value-raw "AAAA"
	// bcdedit /store BCD /set {ObjectId} --value-type RegSz --value "Hello"
	// bcdedit /store BCD /set {ObjectId} --value-type RegMultiSz --value "First" --value "Second"
	"set": {
		Usage: "/set <id> --value-type <ValueType(e.g. RegSz)> --value-raw \"BASE64\"\n" +
			"/set <id> --value-type <ValueType(e.g. RegMultiSz)> --value \"first\" --value \"second\"\n" +
			"This command sets an entry option value in the boot configuration data store.",
		Writable: 1,
		Runner: func(flags *Flags, args []string, bcd go_bcdedit.Bcdedit) error {
			setFlagset := flag.NewFlagSet("", flag.ExitOnError)
			setFlagset.StringVar(&flags.SetValueType, "type", "", "")
			setFlagset.StringVar(&flags.SetValueRaw, "raw", "", "")
			setFlagset.Var(&flags.SetValue, "value", "")

			if strings.Contains(args[0], "help") || strings.Contains(args[0], "?") {
				setFlagset.PrintDefaults()
				return nil
			}

			flags.SetId = args[0]
			flags.SetKey = args[1]
			setFlagset.Parse(args[2:])
			return doSetRaw(flags, bcd)
		},
	},
}

func Main(args []string) {
	var flags Flags
	var err error

	flagset := flag.NewFlagSet(args[0], flag.ExitOnError)
	flagset.BoolVar(&flags.Json, "json", false, "Output result as JSON")
	flagset.StringVar(&flags.Store, "store", "", "Used to specify a BCD store.")

	appliedCommand := make(map[string]*bool)
	for s, def := range commands {
		appliedCommand[s] = flagset.Bool(s, false, def.Usage)
	}

	var fixedArgs []string
	for _, s := range args[1:] {
		if strings.HasPrefix(s, "/") {
			s = "--" + s[1:]
		}
		fixedArgs = append(fixedArgs, s)
	}
	flagset.Parse(fixedArgs)

	err = func() error {
		for s, define := range commands {
			if *appliedCommand[s] {
				var err error
				var bcd go_bcdedit.Bcdedit
				if define.Writable >= 0 && flags.Store != "" {
					bcd, err = go_bcdedit.OpenStore(flags.Store, define.Writable == 1)
				}
				if err != nil {
					return err
				}
				runErr := define.Runner(&flags, flagset.Args(), bcd)
				if bcd == nil {
					return runErr
				}
				closeErr := bcd.Close()
				if runErr != nil {
					return runErr
				}
				return closeErr
			}
		}
		return errors.New("no command")
	}()

	if err != nil {
		log.Panicln(err)
	} else {
		log.Println("The operation completed successfully")
	}
}

func doCreateStore(flags *Flags) error {
	bcd, err := go_bcdedit.CreateStore(flags.Store)
	if err != nil {
		return err
	}
	return bcd.Close()
}

func doEnum(flags *Flags, bcd go_bcdedit.Bcdedit) error {
	objectMap, err := bcd.Enumerate(flags.Enum)
	if err != nil {
		return err
	}

	if flags.Json {
		response := &model.EnumerateResponse{
			Objects: make(map[string]*model.BcdObject),
		}

		for id, object := range objectMap {
			response.Objects[id] = object.ToJson()
		}

		jsonResp, err := json.Marshal(response)
		if err != nil {
			return err
		}
		_, err = os.Stdout.Write(jsonResp)
		return err
	}

	var objectList []go_bcdedit.BcdObject
	for _, object := range objectMap {
		objectList = append(objectList, object)
	}

	slices.SortFunc(objectList, func(a, b go_bcdedit.BcdObject) int {
		if a.GetDescription() < b.GetDescription() {
			return -1
		} else if a.GetDescription() > b.GetDescription() {
			return 1
		}
		return 0
	})

	for _, object := range objectList {
		switch object.GetDescription().ObjectType() {
		case model.ObjectApplication:
			applicationTypeName := object.GetDescription().ApplicationType().String()
			if applicationTypeName == "" {
				applicationTypeName = fmt.Sprintf("unknown type (0x%08x)", object.GetDescription())
			}
			fmt.Printf("%s\n", applicationTypeName)
		case model.ObjectInherit:
			switch object.GetDescription().ObjectSubType() {
			case model.InheritableByApplicationObjects:
				fmt.Printf("%s (inherited)\n", object.GetDescription().ApplicationType().String())
			case model.InheritableByDeviceObjects:
				fmt.Printf("Device options (inherited)\n")
			default:
				fmt.Printf("Inherited\n")
				fmt.Printf("object_description: 0x%08x\n", object.GetDescription())
			}
		case model.ObjectDevice:
			fmt.Printf("%s\n", "Device options")
		default:
			fmt.Printf("OBJECT[ID: %s]:\n", object.GetId())
			fmt.Printf("object_description: 0x%08x\n", object.GetDescription())
		}
		fmt.Printf("%s\n", strings.Repeat("-", 25))
		fmt.Printf("%s %s\n", StringWithPad("Identifier"), ObjectIdToString(object.GetId()))

		for _, element := range object.SortedElements() {
			hiveElement := element.(*go_bcdedit.HiveBcdElement)
			meta := hiveElement.Meta()
			name := hiveElement.Key
			if meta != nil {
				name = meta.Name
			}
			if name == "Inherit" || name == "DisplayOrder" {
				var fixedIds []string
				for _, parentId := range strings.Split(hiveElement.String(), "\n") {
					fixedIds = append(fixedIds, ObjectIdToString(parentId))
				}
				fmt.Printf("%s %s\n", StringWithPad(name), StringWithNewLinePad(strings.Join(fixedIds, "\n")))
			} else {
				fmt.Printf("%s %s\n", StringWithPad(name), StringWithNewLinePad(hiveElement.String()))
			}
		}
		fmt.Printf("\n")
		_ = object
	}

	return nil
}

func doCreateObject(flags *Flags, bcd go_bcdedit.Bcdedit) error {
	if flags.CreateDescription == 0 {
		return errors.New("need object-type")
	}

	object, err := bcd.UpsertObject(flags.CreateId, model.BcdDescription(flags.CreateDescription))
	if err != nil {
		return err
	}

	if len(flags.ObjectDescription) > 0 {
		raw, err := go_bcdedit.StringToUtf16LE(flags.ObjectDescription)
		if err != nil {
			return err
		}
		_, err = object.SetElement("12000004", go_bcdedit.RegSz, raw)
		return err
	}
	return nil
}

func doSetRaw(flags *Flags, bcd go_bcdedit.Bcdedit) error {
	var err error
	var raw []byte
	if flags.SetValueRaw != "" {
		raw, err = base64.StdEncoding.DecodeString(flags.SetValueRaw)
	} else {
		raw, err = DecodeValueToRaw(model.ValueType(flags.SetValueType), flags.SetValue)
	}
	if err != nil {
		return err
	}

	object, err := bcd.GetObject(flags.SetId)
	if err != nil {
		return err
	}

	valueType := go_bcdedit.ValueTypeFromJson(model.ValueType(flags.SetValueType))
	_, err = object.SetElement(flags.SetKey, valueType, raw)
	return err
}

func ObjectIdToString(id string) string {
	known, ok := go_bcdedit.KnownObjectIds[strings.ToLower(id)]
	if ok {
		return known
	}
	return id
}

func StringWithPad(s string) string {
	pad := 24
	return s + strings.Repeat(" ", pad-len(s))
}

func StringWithNewLinePad(s string) string {
	out := ""
	for i, v := range strings.Split(s, "\n") {
		if i == 0 {
			out = v
		} else {
			out += "\n" + strings.Repeat(" ", 25) + v
		}
	}
	return out
}

// See http://www.mistyprojects.co.uk/documents/BCDEdit/files/object_element_codes.htm

func DecodeValueToRaw(valueType model.ValueType, input []string) ([]byte, error) {
	var err error
	switch valueType {
	case model.RegNone:
		return []byte{}, nil
	case model.RegSz:
		return go_bcdedit.StringToUtf16LE(input[0])
	case model.RegBinary:
		return base64.StdEncoding.DecodeString(input[0])
	case model.RegDword:
		var n uint64
		if strings.HasPrefix(input[0], "0x") {
			n, err = strconv.ParseUint(input[0][2:], 16, 32)
		} else {
			n, err = strconv.ParseUint(input[0], 10, 32)
		}
		if err != nil {
			return nil, err
		}
		return binary.LittleEndian.AppendUint32(nil, uint32(n)), nil
	case model.RegMultiSz:
		return go_bcdedit.StringsToMultiUtf16LE(input)
	case model.RegQword:
		var n uint64
		if strings.HasPrefix(input[0], "0x") {
			n, err = strconv.ParseUint(input[0][2:], 16, 64)
		} else {
			n, err = strconv.ParseUint(input[0], 10, 64)
		}
		if err != nil {
			return nil, err
		}
		return binary.LittleEndian.AppendUint64(nil, n), nil
	}
	return nil, fmt.Errorf("not supported type: %s", valueType)
}
