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
	//"io/ioutil"
	//"path/filepath"
	//"encoding/json"
	//"os"
	"reflect"
	gtm "github.com/akamai/AkamaiOPEN-edgegrid-golang/configgtm-v1_4"
	"fmt"

)

// property
var gtmPropertyConfigP1 = fmt.Sprintf(`
resource "akamai_gtm_property" `)

// Process property resources
func processProperties(properties []*gtm.Property, pImportList map[string][]int, dcImportList map [int]string, resourceDomainName string) (string) {

	propertiesString := ""
	for _, property := range properties {
		if _, ok := pImportList[property.Name]; !ok {
			continue
		}
        	propertyBody := ""
        	name := ""
        	propString := gtmPropertyConfigP1
        	propElems := reflect.ValueOf(property).Elem()
        	for i := 0; i < propElems.NumField(); i++ {
                	varName := propElems.Type().Field(i).Name
                	varType := propElems.Type().Field(i).Type
                	varValue := propElems.Field(i).Interface()
                	key := convertKey(varName)
                	if key == "" {
				continue
			}
                	keyVal := fmt.Sprint(varValue)
                        if keyVal == "" && varType.Kind() == reflect.String {
                                continue
                        }
			switch varName {
			case "LivenessTests":
				keyVal = processLivenessTests(varValue.([]*gtm.LivenessTest)) 
                        case "TrafficTargets":
                                keyVal = processTrafficTargets(varValue.([]*gtm.TrafficTarget))
                        case "StaticRRSets":
                                keyVal = processStaticRRSets(varValue.([]*gtm.StaticRRSet))
                        case "MxRecords":
                                continue 	// deprecated in schema 1.4+
			default: 
                		if key == "name" { name = keyVal }
			}
                	if varType.Kind() == reflect.String {
                        	propertyBody += tab4 + key + " = \"" + keyVal + "\"\n"
                	} else {
                        	propertyBody += tab4 + key + " = " + keyVal + "\n"
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

        headerString := "[]\n"			// assume MT
        for ii, header := range headers {
                headerString = "[{\n"		// at least one
                headElems := reflect.ValueOf(header).Elem()
                for i := 0; i < headElems.NumField(); i++ {
                        varName := headElems.Type().Field(i).Name
                        varType := headElems.Type().Field(i).Type
                        varValue := headElems.Field(i).Interface()
                        key := convertKey(varName)
                        keyVal := fmt.Sprint(varValue)
                        if varType.Kind() == reflect.String {
                                headerString += tab12 + key + " = \"" + keyVal + "\"\n"
                        } else {
                                headerString += tab12 + key + " = " + keyVal + "\n"
                        }
                }
		if ii < len(headers)-1 {
			headerString += tab12 + "},\n" + tab12 + "{\n"
		} else {
                        headerString += tab12 + "}\n"
			headerString += tab8 + "]"
		}
	}		
	return headerString
}

func processTrafficTargets(targets []*gtm.TrafficTarget) string {

        targetString := "[]\n"                  // assume MT
        for ii, target := range targets {
                targetString = "[{\n"           // at least one
                targElems := reflect.ValueOf(target).Elem()
                for i := 0; i < targElems.NumField(); i++ {
                        varName := targElems.Type().Field(i).Name
                        varType := targElems.Type().Field(i).Type
                        varValue := targElems.Field(i).Interface()
                        key := convertKey(varName)
                        keyVal := fmt.Sprint(varValue)
			if varName == "Servers" {
				keyVal = processStringList(target.Servers)
			}
                        if varType.Kind() == reflect.String {
                                targetString += tab8 + key + " = \"" + keyVal + "\"\n"
                        } else {
                                targetString += tab8 + key + " = " + keyVal + "\n"
                        }
                }
                if ii < len(targets)-1 {
                        targetString += tab8 + "},\n" + tab8 + "{\n"
                } else {
                        targetString += tab8 + "}\n"
                        targetString += tab4 + "]"
                }
        }
        return targetString

}

func processLivenessTests(tests []*gtm.LivenessTest) string {

        testsString := "[]\n"			// assume MT
        for ii, test := range tests {
                testsString = "[{\n"			// at least one
                testElems := reflect.ValueOf(test).Elem()
                for i := 0; i < testElems.NumField(); i++ {
                        varName := testElems.Type().Field(i).Name
                        varType := testElems.Type().Field(i).Type
                        varValue := testElems.Field(i).Interface()
                        key := convertKey(varName)
                        if key == "" {
                                continue
                        }
                        keyVal := fmt.Sprint(varValue)
                        if varName == "HttpHeaders" {
                                keyVal = processHttpHeaders(varValue.([]*gtm.HttpHeader))
                        }
                        if varType.Kind() == reflect.String {
                                testsString += tab8 + key + "  = \"" + keyVal + "\"\n"
                        } else {
                                testsString += tab8 + key + " = " + keyVal + "\n"
                        }
                }
                if ii < len(tests)-1 {
                        testsString += tab8 + "},\n" + tab8 + "{\n"
                } else {
                        testsString += tab8 + "}\n"
			testsString += tab4 + "]"
                }
        }
        return testsString

}

func processStaticRRSets(rrsets []*gtm.StaticRRSet) string {

        setString := "[]\n"                  // assume MT
        for ii, set := range rrsets {
                setString = "[{\n"           // at least one
                setElems := reflect.ValueOf(set).Elem()
                for i := 0; i < setElems.NumField(); i++ {
                        varName := setElems.Type().Field(i).Name
                        varType := setElems.Type().Field(i).Type
                        varValue := setElems.Field(i).Interface()
                        key := convertKey(varName)
                        keyVal := fmt.Sprint(varValue)
                        if varName == "Rdata" {
                                keyVal = processStringList(set.Rdata)
                        }       
                        if varType.Kind() == reflect.String {
                                setString += tab8 + key + " = \"" + keyVal + "\"\n"
                        } else {
                                setString += tab8 + key + " = " + keyVal + "\n"
                        }       
                }
                if ii < len(rrsets)-1 {
                        setString += tab8 + "},\n" + tab8 + "{\n"
                } else {
                        setString += tab8 + "}\n"
                        setString += tab4 + "]"
                }       
        }
        return setString

}

