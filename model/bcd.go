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

type BcdDescription uint32

func BcdDescriptionFrom(ObjectType ObjectType, ObjectSubType ObjectSubType, ApplicationType ApplicationType) BcdDescription {
	return BcdDescription(uint32(ObjectType) | uint32(ObjectSubType) | uint32(ApplicationType))
}

func (d BcdDescription) ObjectType() ObjectType {
	return ObjectType(d & 0xf0000000)
}

func (d BcdDescription) ObjectSubType() ObjectSubType {
	return ObjectSubType(d & 0x00f00000)
}

func (d BcdDescription) ApplicationType() ApplicationType {
	return ApplicationType(d & 0x000fffff)
}
