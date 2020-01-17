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
	"strconv"
)

// property
var gtmPropertyConfigP1 = fmt.Sprintf(`
resource "akamai_gtm_property" `)

// Process property resources
func processProperties(properties []*gtm.Property, pImportList map[string][]int, dcImportList map[int]string, resourceDomainName string) string {

	// Get Null values list
        var coreFieldsNullMap map[string]string
	var childFieldsNullMap  map[string]interface{}

        nullFieldsMap := getNullValuesList("Properties")       

	propertiesString := ""
	for _, property := range properties {
		if _, ok := pImportList[property.Name]; !ok {
			continue
		}
		// Retrieve Core null fields map
		if propNullFieldObjectMap, ok := nullFieldsMap[property.Name]; ok {
			coreFieldsNullMap = propNullFieldObjectMap.CoreObjectFields
			childFieldsNullMap = propNullFieldObjectMap.ChildObjectFields
		} else {
			coreFieldsNullMap = map[string]string{}
			childFieldsNullMap = map[string]interface{}{}
		} 
		propertyBody := ""
		name := ""
		propString := gtmPropertyConfigP1
		propElems := reflect.ValueOf(property).Elem()
		for i := 0; i < propElems.NumField(); i++ {
			varName := propElems.Type().Field(i).Name
			varType := propElems.Type().Field(i).Type
			varValue := propElems.Field(i).Interface()
 			if _, ok := coreFieldsNullMap[varName]; ok {
				continue
			}
			keyVal := fmt.Sprint(varValue)
			key := convertKey(varName, keyVal, varType.Kind())
			if key == "" {
				continue
			}
			switch varName {
			case "LivenessTests":
                                if _, ok := childFieldsNullMap[varName]; !ok {
                                        continue
                                }
				propertyBody += processLivenessTests(varValue.([]*gtm.LivenessTest), childFieldsNullMap[varName].(map[string]gtm.NullPerObjectAttributeStruct))
			case "TrafficTargets":
                                if _, ok := childFieldsNullMap[varName]; !ok {
                                        continue
                                }
				propertyBody += processTrafficTargets(varValue.([]*gtm.TrafficTarget), childFieldsNullMap[varName].(map[string]gtm.NullPerObjectAttributeStruct))
			case "StaticRRSets":
				if _, ok := childFieldsNullMap[varName]; !ok {
					continue
				}
				propertyBody += processStaticRRSets(varValue.([]*gtm.StaticRRSet), childFieldsNullMap[varName].(map[string]gtm.NullPerObjectAttributeStruct))
			case "MxRecords":
				continue // deprecated in schema 1.4+
			default:
				if key == "name" {
					name = keyVal
				}
				if varType.Kind() == reflect.String {
					propertyBody += tab4 + key + " = \"" + keyVal + "\"\n"
				} else {
					propertyBody += tab4 + key + " = " + keyVal + "\n"
				}
			}
		}
		propString += "\"" + name + "\" {\n"
		propString += gtmRConfigP2 + resourceDomainName + ".name}\"\n"
		propString += propertyBody
		propString += dependsClauseP1 + resourceDomainName + "\""
		// process dc dependencies (only one type in 1.4 schema)
		for _, dcDep := range pImportList[name] {
			propString += ",\n"
			propString += tab8 + "\"" + datacenterResource + "." + dcImportList[dcDep] + "\""
		}
		propString += "\n"
		propString += tab4 + "]\n"
		propString += "}\n"
		propertiesString += propString
	}

	return propertiesString

}

func processHttpHeaders(headers []*gtm.HttpHeader) string {

	if len(headers) == 0 {
		return ""
	}
	headerString := ""
	for _, header := range headers {
		headerString += tab8 + "http_header {\n"
		headElems := reflect.ValueOf(header).Elem()
		for i := 0; i < headElems.NumField(); i++ {
			varName := headElems.Type().Field(i).Name
			varType := headElems.Type().Field(i).Type
			varValue := headElems.Field(i).Interface()
			keyVal := fmt.Sprint(varValue)
			key := convertKey(varName, keyVal, varType.Kind())
			if varType.Kind() == reflect.String {
				headerString += tab12 + key + " = \"" + keyVal + "\"\n"
			} else {
				headerString += tab12 + key + " = " + keyVal + "\n"
			}
		}
		headerString += tab8 + "}\n"
	}

	return headerString
}

func processTrafficTargets(targets []*gtm.TrafficTarget, childObjectList map[string]gtm.NullPerObjectAttributeStruct) string {

	if len(targets) == 0 {
		return ""
	}
	targetString := ""
	for _, target := range targets {
		targetName := strconv.Itoa(target.DatacenterId)
		trgNullFields := childObjectList[targetName]
		targetString += tab4 + "traffic_target {\n"
		targElems := reflect.ValueOf(target).Elem()
		for i := 0; i < targElems.NumField(); i++ {
			varName := targElems.Type().Field(i).Name
			varType := targElems.Type().Field(i).Type
			varValue := targElems.Field(i).Interface()
			keyVal := fmt.Sprint(varValue)
			if _, ok := trgNullFields.CoreObjectFields[varName]; ok {
                                continue
                        }
			key := convertKey(varName, keyVal, varType.Kind())
			if varName == "Servers" {
				keyVal = processStringList(target.Servers)
			}
			if varType.Kind() == reflect.String {
				targetString += tab8 + key + " = \"" + keyVal + "\"\n"
			} else {
				targetString += tab8 + key + " = " + keyVal + "\n"
			}
		}
		targetString += tab4 + "}\n"
	}

	return targetString

}

func processLivenessTests(tests []*gtm.LivenessTest, childObjectList map[string]gtm.NullPerObjectAttributeStruct) string {

	if len(tests) == 0 {
		return ""
	}
	testsString := ""
	for _, test := range tests {
		liveNullFields := childObjectList[test.Name]
		testsString += tab4 + "liveness_test {\n"
		testElems := reflect.ValueOf(test).Elem()
		for i := 0; i < testElems.NumField(); i++ {
			varName := testElems.Type().Field(i).Name
			varType := testElems.Type().Field(i).Type
			varValue := testElems.Field(i).Interface()
                        if _, ok := liveNullFields.CoreObjectFields[varName]; ok {
                                continue
                        }
			keyVal := fmt.Sprint(varValue)
			key := convertKey(varName, keyVal, varType.Kind())
			if key == "" {
				continue
			}
			if varName == "HttpHeaders" {
				testsString += processHttpHeaders(varValue.([]*gtm.HttpHeader))
				continue
			}
			if varType.Kind() == reflect.String {
				testsString += tab8 + key + "  = \"" + keyVal + "\"\n"
			} else {
				testsString += tab8 + key + " = " + keyVal + "\n"
			}
		}

		testsString += tab4 + "}\n"
	}
	return testsString

}

func processStaticRRSets(rrsets []*gtm.StaticRRSet, childObjectList map[string]gtm.NullPerObjectAttributeStruct) string {

	if len(rrsets) == 0 {
		return ""
	}
	setString := ""
	for _, rrset := range rrsets {
		// There are no default null values or children in rrset as of the 1.4 schema
		//rrsetNullFields := childObjectList[rrset.Name]
		setString += tab4 + "static_rr_set {\n"
		setElems := reflect.ValueOf(rrset).Elem()
		for i := 0; i < setElems.NumField(); i++ {
			varName := setElems.Type().Field(i).Name
			varType := setElems.Type().Field(i).Type
			varValue := setElems.Field(i).Interface()
			keyVal := fmt.Sprint(varValue)
			key := convertKey(varName, keyVal, varType.Kind())
			if varName == "Rdata" {
				keyVal = processStringList(rrset.Rdata)
			}
			if varType.Kind() == reflect.String {
				setString += tab8 + key + " = \"" + keyVal + "\"\n"
			} else {
				setString += tab8 + key + " = " + keyVal + "\n"
			}
		}

		setString += tab4 + "}\n"
	}
	return setString

}
