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

type BcdElement struct {
	Type         ValueType `json:"type"`
	Raw          string    `json:"raw"` // base64 encoded
	ValueSz      string    `json:"valueSz,omitempty"`
	ValueMultiSz []string  `json:"valueMultiSz,omitempty"`
	ValueDword   *uint32   `json:"valueDword,omitempty"`
}

type BcdObject struct {
	Description BcdDescription         `json:"description"`
	Elements    map[string]*BcdElement `json:"elements"` // e.g. key="11000001"
}

type EnumerateResponse struct {
	Objects map[string]*BcdObject `json:"objects"` // e.g. key="{b2721d73-1db4-4c62-bf78-c548a880142d}"
}
