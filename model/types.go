// Copyright 2024 JC-Lab
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package model

type ObjectType uint32
type ObjectSubType uint32
type ApplicationType uint32

const (
	ObjectApplication ObjectType = 0x10000000
	ObjectInherit     ObjectType = 0x20000000
	ObjectDevice      ObjectType = 0x30000000

	FirmwareApplication     ObjectSubType = 0x100000
	WindowsBootApplication  ObjectSubType = 0x200000
	LegacyLoaderApplication ObjectSubType = 0x300000
	RealModeApplication     ObjectSubType = 0x400000

	InheritableByAnyObject          ObjectSubType = 0x100000
	InheritableByApplicationObjects ObjectSubType = 0x200000
	InheritableByDeviceObjects      ObjectSubType = 0x300000

	ApplicationFwbootmgr  ApplicationType = 1
	ApplicationBootmgr    ApplicationType = 2
	ApplicationOsloader   ApplicationType = 3
	ApplicationResume     ApplicationType = 4
	ApplicationMemdiag    ApplicationType = 5
	ApplicationNtldr      ApplicationType = 6
	ApplicationSetupldr   ApplicationType = 7
	ApplicationBootsector ApplicationType = 8
	ApplicationStartup    ApplicationType = 9
	ApplicationBootapp    ApplicationType = 10
)

// ApplicationType의 String 메서드 구현
func (a ApplicationType) String() string {
	switch a {
	case ApplicationFwbootmgr:
		return "fwbootmgr"
	case ApplicationBootmgr:
		return "bootmgr"
	case ApplicationOsloader:
		return "osloader"
	case ApplicationResume:
		return "resume"
	case ApplicationMemdiag:
		return "memdiag"
	case ApplicationNtldr:
		return "ntldr"
	case ApplicationSetupldr:
		return "setupldr"
	case ApplicationBootsector:
		return "bootsector"
	case ApplicationStartup:
		return "startup"
	case ApplicationBootapp:
		return "bootapp"
	default:
		return ""
	}
}

type ValueType string

const (
	RegNone                     ValueType = "RegNone"
	RegSz                       ValueType = "RegSz"
	RegExpandSz                 ValueType = "RegExpandSz"
	RegBinary                   ValueType = "RegBinary"
	RegDword                    ValueType = "RegDword"
	RegDwordBigEndian           ValueType = "RegDwordBigEndian"
	RegLink                     ValueType = "RegLink"
	RegMultiSz                  ValueType = "RegMultiSz"
	RegResourceList             ValueType = "RegResourceList"
	RegFullResourceDescriptor   ValueType = "RegFullResourceDescriptor"
	RegResourceRequirementsList ValueType = "RegResourceRequirementsList"
	RegQword                    ValueType = "RegQword"
)

type BcdElementMeta struct {
	Name   string
	Format string
}

var GenericElementTypes = map[string]*BcdElementMeta{
	"11000001": {
		Name:   "Device",
		Format: "Device",
	},
	"12000002": {
		Name:   "Path",
		Format: "string",
	},
	"12000004": {
		Name:   "Description",
		Format: "string",
	},
	"12000005": {
		Name:   "Locale",
		Format: "string",
	},
	"14000006": {
		Name:   "Inherit",
		Format: "GUID list",
	},
	"14000008": {
		Name:   "RecoverySequence",
		Format: "GUID list",
	},
	"16000009": {
		Name:   "RecoveryEnabled",
		Format: "boolean",
	},
}

// BcdBootMgrElementTypes http://msdn.microsoft.com/en-us/library/windows/desktop/aa362641(v=vs.85).aspx
var BcdBootMgrElementTypes = map[string]*BcdElementMeta{
	"24000001": {
		Name:   "DisplayOrder",
		Format: "ObjectList",
	},
	"24000002": {
		Name:   "BootSequence",
		Format: "ObjectList",
	},
	"23000003": {
		Name:   "DefaultObject",
		Format: "Object",
	},
	"25000004": {
		Name:   "Timeout",
		Format: "Integer",
	},
	"26000005": {
		Name:   "AttemptResume",
		Format: "Boolean",
	},
	"23000006": {
		Name:   "ResumeObject",
		Format: "Object",
	},
	"24000010": {
		Name:   "ToolsDisplayOrder",
		Format: "ObjectList",
	},
	"26000020": {
		Name:   "DisplayBootMenu",
		Format: "Boolean",
	},
	"26000021": {
		Name:   "NoErrorDisplay",
		Format: "Boolean",
	},
	"21000022": {
		Name:   "BcdDevice",
		Format: "Device",
	},
	"22000023": {
		Name:   "BcdFilePath",
		Format: "String",
	},
	"26000028": {
		Name:   "ProcessCustomActionsFirst",
		Format: "Boolean",
	},
	"27000030": {
		Name:   "CustomActionsList",
		Format: "IntegerList",
	},
	"26000031": {
		Name:   "PersistBootSequence",
		Format: "Boolean",
	},
}

// BcdDeviceElementTypes https://learn.microsoft.com/ko-kr/previous-versions/windows/desktop/bcd/bcddeviceobjectelementtypes?redirectedfrom=MSDN
var BcdDeviceElementTypes = map[string]*BcdElementMeta{
	"35000001": {
		Name:   "RamdiskImageOffset",
		Format: "Integer",
	},
	"35000002": {
		Name:   "TftpClientPort",
		Format: "Integer",
	},
	"31000003": {
		Name:   "SdiDevice",
		Format: "Integer",
	},
	"32000004": {
		Name:   "SdiPath",
		Format: "Integer",
	},
	"35000005": {
		Name:   "RamdiskImageLength",
		Format: "Integer",
	},
	"36000006": {
		Name:   "RamdiskExportAsCd",
		Format: "Boolean",
	},
	"36000007": {
		Name:   "RamdiskTftpBlockSize",
		Format: "Integer",
	},
	"36000008": {
		Name:   "RamdiskTftpWindowSize",
		Format: "Integer",
	},
	"36000009": {
		Name:   "RamdiskMulticastEnabled",
		Format: "Boolean",
	},
	"3600000A": {
		Name:   "RamdiskMulticastTftpFallback",
		Format: "Boolean",
	},
	"3600000B": {
		Name:   "RamdiskTftpVarWindow",
		Format: "Boolean",
	},
}

// BcdOsLoaderElementTypes https://learn.microsoft.com/ko-kr/previous-versions/windows/desktop/bcd/bcdosloaderelementtypes?redirectedfrom=MSDN
var BcdOsLoaderElementTypes = map[string]*BcdElementMeta{
	"21000001": {
		Name:   "OSDevice",
		Format: "Device",
	},
	"22000002": {
		Name:   "SystemRoot",
		Format: "String",
	},
	"23000003": {
		Name:   "AssociatedResumeObject",
		Format: "Object",
	},
	"26000010": {
		Name:   "DetectKernelAndHal",
		Format: "Boolean",
	},
	"22000011": {
		Name:   "KernelPath",
		Format: "String",
	},
	"22000012": {
		Name:   "HalPath",
		Format: "String",
	},
	"22000013": {
		Name:   "DbgTransportPath",
		Format: "String",
	},
	"25000020": {
		Name:   "NxPolicy",
		Format: "Integer",
	},
	"25000021": {
		Name:   "PAEPolicy",
		Format: "Integer",
	},
	"26000022": {
		Name:   "WinPEMode",
		Format: "Boolean",
	},
	"26000024": {
		Name:   "DisableCrashAutoReboot",
		Format: "Boolean",
	},
	"26000025": {
		Name:   "UseLastGoodSettings",
		Format: "Boolean",
	},
	"26000027": {
		Name:   "AllowPrereleaseSignatures",
		Format: "Boolean",
	},
	"26000030": {
		Name:   "NoLowMemory",
		Format: "Boolean",
	},
	"25000031": {
		Name:   "RemoveMemory",
		Format: "Integer",
	},
	"25000032": {
		Name:   "IncreaseUserVa",
		Format: "Integer",
	},
	"26000040": {
		Name:   "UseVgaDriver",
		Format: "Boolean",
	},
	"26000041": {
		Name:   "DisableBootDisplay",
		Format: "Boolean",
	},
	"26000042": {
		Name:   "DisableVesaBios",
		Format: "Boolean",
	},
	"26000043": {
		Name:   "DisableVgaMode",
		Format: "Boolean",
	},
	"25000050": {
		Name:   "ClusterModeAddressing",
		Format: "Integer",
	},
	"26000051": {
		Name:   "UsePhysicalDestination",
		Format: "Boolean",
	},
	"25000052": {
		Name:   "RestrictApicCluster",
		Format: "Integer",
	},
	"26000054": {
		Name:   "UseLegacyApicMode",
		Format: "Boolean",
	},
	"25000055": {
		Name:   "X2ApicPolicy",
		Format: "Integer",
	},
	"26000060": {
		Name:   "UseBootProcessorOnly",
		Format: "Boolean",
	},
	"25000061": {
		Name:   "NumberOfProcessors",
		Format: "Integer",
	},
	"26000062": {
		Name:   "ForceMaximumProcessors",
		Format: "Boolean",
	},
	"25000063": {
		Name:   "ProcessorConfigurationFlags",
		Format: "Boolean",
	},
	"26000064": {
		Name:   "MaximizeGroupsCreated",
		Format: "Boolean",
	},
	"26000065": {
		Name:   "ForceGroupAwareness",
		Format: "Boolean",
	},
	"25000066": {
		Name:   "GroupSize",
		Format: "Integer",
	},
	"26000070": {
		Name:   "UseFirmwarePciSettings",
		Format: "Integer",
	},
	"25000071": {
		Name:   "MsiPolicy",
		Format: "Integer",
	},
	"25000080": {
		Name:   "SafeBoot",
		Format: "Integer",
	},
	"26000081": {
		Name:   "SafeBootAlternateShell",
		Format: "Boolean",
	},
	"26000090": {
		Name:   "BootLogInitialization",
		Format: "Boolean",
	},
	"26000091": {
		Name:   "VerboseObjectLoadMode",
		Format: "Boolean",
	},
	"260000a0": {
		Name:   "KernelDebuggerEnabled",
		Format: "Boolean",
	},
	"260000a1": {
		Name:   "DebuggerHalBreakpoint",
		Format: "Boolean",
	},
	"260000A2": {
		Name:   "UsePlatformClock",
		Format: "Boolean",
	},
	"260000A3": {
		Name:   "ForceLegacyPlatform",
		Format: "Boolean",
	},
	"250000A6": {
		Name:   "TscSyncPolicy",
		Format: "Integer",
	},
	"260000b0": {
		Name:   "EmsEnabled",
		Format: "Boolean",
	},
	"250000c1": {
		Name:   "DriverLoadFailurePolicy",
		Format: "Integer",
	},
	"250000C2": {
		Name:   "BootMenuPolicy",
		Format: "Integer",
	},
	"260000C3": {
		Name:   "AdvancedOptionsOneTime",
		Format: "Boolean",
	},
	"250000E0": {
		Name:   "BootStatusPolicy",
		Format: "Integer",
	},
	"260000E1": {
		Name:   "DisableElamDrivers",
		Format: "Boolean",
	},
	"250000F0": {
		Name:   "HypervisorLaunchType",
		Format: "Integer",
	},
	"260000F2": {
		Name:   "HypervisorDebuggerEnabled",
		Format: "Boolean",
	},
	"250000F3": {
		Name:   "HypervisorDebuggerType",
		Format: "Integer",
	},
	"250000F4": {
		Name:   "HypervisorDebuggerPortNumber",
		Format: "Integer",
	},
	"250000F5": {
		Name:   "HypervisorDebuggerBaudrate",
		Format: "Integer",
	},
	"250000F6": {
		Name:   "HypervisorDebugger1394Channel",
		Format: "Integer",
	},
	"250000F7": {
		Name:   "BootUxPolicy",
		Format: "Integer",
	},
	"220000F9": {
		Name:   "HypervisorDebuggerBusParams",
		Format: "String",
	},
	"250000FA": {
		Name:   "HypervisorNumProc",
		Format: "Integer",
	},
	"250000FB": {
		Name:   "HypervisorRootProcPerNode",
		Format: "Integer",
	},
	"260000FC": {
		Name:   "HypervisorUseLargeVTlb",
		Format: "Boolean",
	},
	"250000FD": {
		Name:   "HypervisorDebuggerNetHostIp",
		Format: "Integer",
	},
	"250000FE": {
		Name:   "HypervisorDebuggerNetHostPort",
		Format: "Integer",
	},
	"25000100": {
		Name:   "TpmBootEntropyPolicy",
		Format: "Integer",
	},
	"22000110": {
		Name:   "HypervisorDebuggerNetKey",
		Format: "String",
	},
	"26000114": {
		Name:   "HypervisorDebuggerNetDhcp",
		Format: "Boolean",
	},
	"25000115": {
		Name:   "HypervisorIommuPolicy",
		Format: "Integer",
	},
	"2500012b": {
		Name:   "XSaveDisable",
		Format: "Integer",
	},
}

var BcdApplicationElementTypes = map[ApplicationType]map[string]*BcdElementMeta{}

func init() {
	BcdApplicationElementTypes[ApplicationBootmgr] = concatBcdElementTypes(
		GenericElementTypes,
		BcdBootMgrElementTypes,
	)
	BcdApplicationElementTypes[ApplicationFwbootmgr] = concatBcdElementTypes(
		GenericElementTypes,
		BcdBootMgrElementTypes,
	)
	BcdApplicationElementTypes[ApplicationOsloader] = concatBcdElementTypes(
		GenericElementTypes,
		BcdOsLoaderElementTypes,
	)
}

func concatBcdElementTypes(inputs ...map[string]*BcdElementMeta) map[string]*BcdElementMeta {
	concated := make(map[string]*BcdElementMeta)
	for _, input := range inputs {
		for key, meta := range input {
			concated[key] = meta
		}
	}
	return concated
}
