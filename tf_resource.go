// Copyright 2020. Akamai Technologies, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"fmt"
	gtm "github.com/akamai/AkamaiOPEN-edgegrid-golang/configgtm-v1_4"
	"reflect"
)

// resource config
var gtmResourceConfigP1 = fmt.Sprintf(`
resource "akamai_gtm_resource" `)

// Process resource resources
func processResources(resources []*gtm.Resource, rImportList map[string][]int, dcIL map[int]string, resourceDomainName string) string {

	// Get Null values list
	var coreFieldsNullMap map[string]string
	nullFieldsMap := getNullValuesList("Resources")

	resourcesString := ""
	for _, resource := range resources {
		if _, ok := rImportList[resource.Name]; !ok {
			continue
		}
		// Retrieve Core null fields map
		if rsrcNullFieldObjectMap, ok := nullFieldsMap[resource.Name]; ok {
			coreFieldsNullMap = rsrcNullFieldObjectMap.CoreObjectFields
		} else {
			coreFieldsNullMap = map[string]string{}
		}
		resourceBody := ""
		name := ""
		rString := gtmResourceConfigP1
		rElems := reflect.ValueOf(resource).Elem()
		for i := 0; i < rElems.NumField(); i++ {
			varName := rElems.Type().Field(i).Name
			varType := rElems.Type().Field(i).Type
			varValue := rElems.Field(i).Interface()
			if _, ok := coreFieldsNullMap[varName]; ok {
				continue
			}
			keyVal := fmt.Sprint(varValue)
			key := convertKey(varName, keyVal, varType.Kind())
			if key == "" {
				continue
			}
			if key == "name" {
				name = keyVal
			}
			if varName == "ResourceInstances" {
				resourceBody += processResourceInstances(varValue.([]*gtm.ResourceInstance), dcIL)
				continue
			}
			resourceBody += tab4 + key + " = "
			if varType.Kind() == reflect.String {
				resourceBody += "\"" + keyVal + "\"\n"
			} else {
				resourceBody += keyVal + "\n"
			}
		}
		rString += "\"" + normalizeResourceName(name) + "\" {\n"
		rString += gtmRConfigP2 + resourceDomainName + ".name\n"
		rString += resourceBody
		rString += dependsClauseP1 + resourceDomainName
		// process dc dependencies (only one type in 1.4 schema)
		for _, dcDep := range rImportList[name] {
			rString += ",\n"
			rString += tab8 + datacenterResource + "." + normalizeResourceName(dcIL[dcDep])
		}
		rString += "\n"
		rString += tab4 + "]\n"
		rString += "}\n"
		resourcesString += rString
	}

	return resourcesString

}

func processResourceInstances(instances []*gtm.ResourceInstance, dcIDs map[int]string) string {

	if len(instances) == 0 {
		return ""
	}
	instanceString := ""
	for _, instance := range instances {
		instanceString += tab4 + "resource_instance {\n"
		instElems := reflect.ValueOf(instance).Elem()
		for i := 0; i < instElems.NumField(); i++ {
			varName := instElems.Type().Field(i).Name
			varType := instElems.Type().Field(i).Type
			varValue := instElems.Field(i).Interface()
			keyVal := fmt.Sprint(varValue)
			key := convertKey(varName, keyVal, varType.Kind())
			if varName == "LoadObject" {
				keyVal = processLoadObject(&instance.LoadObject)
				instanceString += keyVal
				continue
			}
			if varType.Kind() == reflect.String {
				instanceString += tab8 + key + " = \"" + keyVal + "\"\n"
			} else {
				// check for datacenter_id
				if varName == "DatacenterId" {
					instanceString += tab8 + key + " = " + datacenterResource + "." + normalizeResourceName(dcIDs[varValue.(int)]) + ".datacenter_id\n"
				} else {
					instanceString += tab8 + key + " = " + keyVal + "\n"
				}
			}
		}
		instanceString += tab4 + "}\n"
	}

	return instanceString

}
