// Copyright 2018 The casbin Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package server

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/casbin/casbin/model"
)

type ABACModel struct {
	V0     string
	V1     string
	V2     string
	V3     string
	V4     string
	V5     string
	V6     string
	V7     string
	V8     string
	V9     string
	V10    string
	source map[string]string
}

func toUpperFirstChar(str string) string {
	for i, v := range str {
		return string(unicode.ToUpper(v)) + str[i+1:]
	}
	return ""
}

func MakeABAC(obj interface{}) (string, error) {
	data, err := json.Marshal(&obj)
	if err != nil {
		return "", err
	}
	return "ABAC::" + string(data), nil
}

func resolveABAC(obj string) (ABACModel, error) {
	var jsonMap map[string]interface{}
	model := ABACModel{source: map[string]string{}}

	err := json.Unmarshal([]byte(obj[len("ABAC::"):]), &jsonMap)
	if err != nil {
		return model, err
	}

	i := 0
	for k, v := range jsonMap {
		key := toUpperFirstChar(k)
		value := fmt.Sprintf("%v", v)
		model.source[key] = "V" + strconv.Itoa(i)
		switch i {
		case 0:
			model.V0 = value
		case 1:
			model.V1 = value
		case 2:
			model.V2 = value
		case 3:
			model.V3 = value
		case 4:
			model.V4 = value
		case 5:
			model.V5 = value
		case 6:
			model.V6 = value
		case 7:
			model.V7 = value
		case 8:
			model.V8 = value
		case 9:
			model.V9 = value
		case 10:
			model.V10 = value
		}
		i++
	}

	return model, nil
}

func parseAbacParam(param string, m *model.Assertion) interface{} {
	if strings.HasPrefix(param, "ABAC::") == true {
		model, err := resolveABAC(param)
		if err != nil {
			panic(err)
		}
		for k, v := range model.source {
			old := "." + k
			if strings.Contains(m.Value, old) {
				m.Value = strings.Replace(m.Value, old, "."+v, -1)
			}
		}

		return model
	} else {
		return param
	}
}
