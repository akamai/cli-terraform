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

// cidr config 
var gtmCidrmapConfigP1 = fmt.Sprintf(`
resource "akamai_gtm_cidrmap" `)

// geo config    
var gtmGeomapConfigP1 = fmt.Sprintf(`
resource "akamai_gtm_geomap" `)

// as config    
var gtmAsmapConfigP1 = fmt.Sprintf(`
resource "akamai_gtm_asmap" `)

// Process resource cidrmap
func processCidrmaps(cidrmaps []*gtm.CidrMap, cidrImportList map[string][]int, dcIL map[int]string, resourceDomainName string) (string) {

	mapsString := ""
	for _, cmap := range cidrmaps {
		if _, ok := cidrImportList[cmap.Name]; !ok {
			continue
		}
        	mapBody := ""
        	name := ""
        	mString := gtmCidrmapConfigP1
        	mElems := reflect.ValueOf(cmap).Elem()
        	for i := 0; i < mElems.NumField(); i++ {
                	varName := mElems.Type().Field(i).Name
                	varType := mElems.Type().Field(i).Type
                	varValue := mElems.Field(i).Interface()
                	key := convertKey(varName)
                	if key == "" {
				continue
			}
                	keyVal := fmt.Sprint(varValue)
                	if key == "name" { name = keyVal }
			switch varName {
			case "DefaultDatacenter":
				keyVal = processDefaultDatacenter(varValue.(*gtm.DatacenterBase))
                	case "Assignments":
				keyVal = processCidrAssignments(varValue.([]*gtm.CidrAssignment))
               	 	}
			if keyVal == "" && varType.Kind() == reflect.String {
				continue
			}
                	mapBody += tab4 + key + " = "
                	if varType.Kind() == reflect.String {
                        	mapBody += "\"" + keyVal + "\"\n"
                	} else {
                        	mapBody += keyVal + "\n"
                	}
        	}
        	mString += "\"" + name + "\" {\n"
        	mString += gtmRConfigP2 + resourceDomainName + ".name}\""
        	mString += mapBody
		mString += dependsClauseP1 + resourceDomainName + "\""
                // process dc dependencies (only one type in 1.4 schema)
                for _, dcDep := range cidrImportList[name] {
                        mString += ",\n"
                        mString += datacenterResource + "." + dcIL[dcDep]
                }
		mString += "\n"
		mString += tab4 + "]\n"
        	mString += "}\n"
		mapsString += mString
	}

        return mapsString

}

// Process resource geomap
func processGeomaps(geomaps []*gtm.GeoMap, geoImportList map[string][]int, dcIL map[int]string, resourceDomainName string) (string) {

	return ""

}

// Process resource asmap
func processAsmaps(asmaps []*gtm.AsMap, asImportList map[string][]int, dcIL map[int]string, resourceDomainName string) (string) {

	return ""

}

func processDefaultDatacenter(ddc *gtm.DatacenterBase) string {

        ddcString := "[{\n"        
        ddcElems := reflect.ValueOf(ddc).Elem()
        for i := 0; i < ddcElems.NumField(); i++ {
                varName := ddcElems.Type().Field(i).Name
                varType := ddcElems.Type().Field(i).Type
                varValue := ddcElems.Field(i).Interface()
                key := convertKey(varName)
                keyVal := fmt.Sprint(varValue)
                if varType.Kind() == reflect.String {
                        ddcString += tab8 + "\"" + key + "\" = \"" + keyVal + "\"\n"
                } else {
                        ddcString += tab8 + "\"" + key + "\" = " + keyVal + "\n"
                }
        }
        ddcString += tab4 + "}]"
        return ddcString

}

func processCidrAssignments(assigns []*gtm.CidrAssignment) string {

        assignString := "[]\n"                  // assume MT
        for ii, assign := range assigns {
                assignString = "[{\n"           // at least one
                aElems := reflect.ValueOf(assign).Elem()
		assignString += processAssignmentsBase(aElems, "Blocks", processStringList, (ii < len(assigns)-1))
        }
        return assignString

} 

type listCore func([]string) string 


func processAssignmentsBase(elem reflect.Value, assignKey string, fn listCore, last bool) string {

	assignStr := ""
        for i := 0; i < elem.NumField(); i++ {
                varName := elem.Type().Field(i).Name
                varType := elem.Type().Field(i).Type
                varValue := elem.Field(i).Interface()
                key := convertKey(varName)
                keyVal := fmt.Sprint(varValue)
                if varName == assignKey {
                        keyVal = fn(varValue.([]string))
                }
                if varType.Kind() == reflect.String {
                        assignStr += tab8 + "\"" + key + "\" = \"" + keyVal + "\"\n"
                } else {
                        assignStr += tab8 + "\"" + key + "\" = " + keyVal + "\n"
                }
                if last {
                        assignStr += tab8 + "},\n" + tab8 + "{\n"
                } else {
                        assignStr += tab8 + "}\n"
                        assignStr += tab4 + "]"
                }      
        }
	return assignStr

}


