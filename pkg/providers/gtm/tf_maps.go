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

package gtm

import (
	"fmt"
	"reflect"
	"strings"

	gtm "github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/configgtm"
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
func processCidrmaps(cidrmaps []*gtm.CidrMap, cidrImportList map[string][]int, dcIL map[int]string, resourceDomainName string) string {

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
			keyVal := fmt.Sprint(varValue)
			key := convertKey(varName, keyVal, varType.Kind())
			if key == "" {
				continue
			}
			if key == "name" {
				name = keyVal
			}
			switch varName {
			case "DefaultDatacenter":
				mapBody += tab4 + key + " {\n"
				mapBody += processDefaultDatacenter(varValue.(*gtm.DatacenterBase), dcIL, true)
				mapBody += tab4 + "}\n"
				continue
			case "Assignments":
				mapBody += processCidrAssignments(varValue.([]*gtm.CidrAssignment), dcIL)
				continue
			}
			mapBody += tab4 + key + " = "
			if varType.Kind() == reflect.String {
				mapBody += "\"" + keyVal + "\"\n"
			} else {
				mapBody += keyVal + "\n"
			}
		}
		mString += "\"" + normalizeResourceName(name) + "\" {\n"
		mString += gtmRConfigP2 + resourceDomainName + ".name\n"
		mString += mapBody
		mString += dependsClauseP1 + resourceDomainName
		// process dc dependencies (only one type in 1.4 schema)
		for _, dcDep := range cidrImportList[name] {
			mString += ",\n"
			mString += tab8 + datacenterResource + "." + normalizeResourceName(dcIL[dcDep])
		}
		mString += "\n"
		mString += tab4 + "]\n"
		mString += "}\n"
		mapsString += mString
	}

	return mapsString

}

// Process resource geomap
func processGeomaps(geomaps []*gtm.GeoMap, geoImportList map[string][]int, dcIL map[int]string, resourceDomainName string) string {

	mapsString := ""
	for _, gmap := range geomaps {
		if _, ok := geoImportList[gmap.Name]; !ok {
			continue
		}
		mapBody := ""
		name := ""
		mString := gtmGeomapConfigP1
		mElems := reflect.ValueOf(gmap).Elem()
		for i := 0; i < mElems.NumField(); i++ {
			varName := mElems.Type().Field(i).Name
			varType := mElems.Type().Field(i).Type
			varValue := mElems.Field(i).Interface()
			keyVal := fmt.Sprint(varValue)
			key := convertKey(varName, keyVal, varType.Kind())
			if key == "" {
				continue
			}
			if key == "name" {
				name = keyVal
			}
			switch varName {
			case "DefaultDatacenter":
				mapBody += tab4 + key + " {\n"
				mapBody += processDefaultDatacenter(varValue.(*gtm.DatacenterBase), dcIL, true)
				mapBody += tab4 + "}\n"
				continue
			case "Assignments":
				mapBody += processGeoAssignments(varValue.([]*gtm.GeoAssignment), dcIL)
				continue
			}
			mapBody += tab4 + key + " = "
			if varType.Kind() == reflect.String {
				mapBody += "\"" + keyVal + "\"\n"
			} else {
				mapBody += keyVal + "\n"
			}
		}
		mString += "\"" + normalizeResourceName(name) + "\" {\n"
		mString += gtmRConfigP2 + resourceDomainName + ".name\n"
		mString += mapBody
		mString += dependsClauseP1 + resourceDomainName
		// process dc dependencies (only one type in 1.4 schema)
		for _, dcDep := range geoImportList[name] {
			mString += ",\n"
			mString += tab8 + datacenterResource + "." + normalizeResourceName(dcIL[dcDep])
		}
		mString += "\n"
		mString += tab4 + "]\n"
		mString += "}\n"
		mapsString += mString
	}

	return mapsString

}

// Process resource asmap
func processAsmaps(asmaps []*gtm.AsMap, asImportList map[string][]int, dcIL map[int]string, resourceDomainName string) string {

	mapsString := ""
	for _, amap := range asmaps {
		if _, ok := asImportList[amap.Name]; !ok {
			continue
		}
		mapBody := ""
		name := ""
		mString := gtmAsmapConfigP1
		mElems := reflect.ValueOf(amap).Elem()
		for i := 0; i < mElems.NumField(); i++ {
			varName := mElems.Type().Field(i).Name
			varType := mElems.Type().Field(i).Type
			varValue := mElems.Field(i).Interface()
			keyVal := fmt.Sprint(varValue)
			key := convertKey(varName, keyVal, varType.Kind())
			if key == "" {
				continue
			}
			if key == "name" {
				name = keyVal
			}
			switch varName {
			case "DefaultDatacenter":
				mapBody += tab4 + key + " {\n"
				mapBody += processDefaultDatacenter(varValue.(*gtm.DatacenterBase), dcIL, true)
				mapBody += tab4 + "}\n"
				continue
			case "Assignments":
				mapBody += processAsAssignments(varValue.([]*gtm.AsAssignment), dcIL)
				continue
			}
			mapBody += tab4 + key + " = "
			if varType.Kind() == reflect.String {
				mapBody += "\"" + keyVal + "\"\n"
			} else {
				mapBody += keyVal + "\n"
			}
		}
		mString += "\"" + normalizeResourceName(name) + "\" {\n"
		mString += gtmRConfigP2 + resourceDomainName + ".name\n"
		mString += mapBody
		mString += dependsClauseP1 + resourceDomainName
		// process dc dependencies (only one type in 1.4 schema)
		for _, dcDep := range asImportList[name] {
			mString += ",\n"
			mString += tab8 + datacenterResource + "." + normalizeResourceName(dcIL[dcDep])
		}
		mString += "\n"
		mString += tab4 + "]\n"
		mString += "}\n"
		mapsString += mString
	}

	return mapsString

}

func processDefaultDatacenter(ddc *gtm.DatacenterBase, dcIDs map[int]string, structreq bool) string {

	ddcString := ""
	ddcElems := reflect.ValueOf(ddc).Elem()
	for i := 0; i < ddcElems.NumField(); i++ {
		varName := ddcElems.Type().Field(i).Name
		varType := ddcElems.Type().Field(i).Type
		varValue := ddcElems.Field(i).Interface()
		keyVal := fmt.Sprint(varValue)
		key := convertKey(varName, keyVal, varType.Kind())
		if varType.Kind() == reflect.String {
			ddcString += tab8 + key + " = \"" + keyVal + "\"\n"
		} else {
			if varName != "DatacenterId" {
				// Special handling. If we are here for a default dc object, we don't have a resource configured.
				ddcString += tab8 + key + " = " + keyVal + "\n"
			} else {
				if structreq && keyVal == "5400" {
					ddcString += tab8 + key + " = " + defaultDatacenterDataSource + "." + "default_datacenter_5400.datacenter_id\n"
				} else {
					ddcString += tab8 + key + " = " + datacenterResource + "." + normalizeResourceName(dcIDs[varValue.(int)]) + ".datacenter_id\n"
				}
			}
		}
	}
	if !structreq {
		ddcString = strings.TrimSuffix(ddcString, "\n") // remove trailing new line
	}

	return ddcString

}

func processCidrAssignments(assigns []*gtm.CidrAssignment, dcIDs map[int]string) string {

	if len(assigns) == 0 {
		return ""
	}
	assignString := ""
	for _, assign := range assigns {
		aElems := reflect.ValueOf(assign).Elem()
		assignString += processAssignmentsBase(aElems, "Blocks", dcIDs)
	}

	return assignString

}

func processGeoAssignments(assigns []*gtm.GeoAssignment, dcIDs map[int]string) string {

	if len(assigns) == 0 {
		return ""
	}
	assignString := ""
	for _, assign := range assigns {
		aElems := reflect.ValueOf(assign).Elem()
		assignString += processAssignmentsBase(aElems, "Countries", dcIDs)
	}

	return assignString

}

func processAsAssignments(assigns []*gtm.AsAssignment, dcIDs map[int]string) string {

	if len(assigns) == 0 {
		return ""
	}
	assignString := ""
	for _, assign := range assigns {
		aElems := reflect.ValueOf(assign).Elem()
		assignString += processAssignmentsBase(aElems, "AsNumbers", dcIDs)
	}

	return assignString

}

func processAssignmentsBase(elem reflect.Value, assignKey string, dcIDs map[int]string) string {

	assignStr := tab4 + "assignment {\n"
	for i := 0; i < elem.NumField(); i++ {
		varName := elem.Type().Field(i).Name
		varType := elem.Type().Field(i).Type
		varValue := elem.Field(i).Interface()
		keyVal := fmt.Sprint(varValue)
		key := convertKey(varName, keyVal, varType.Kind())
		if varName == "DatacenterBase" {
			dcb := varValue.(gtm.DatacenterBase)
			keyVal = processDefaultDatacenter(&dcb, dcIDs, false)
			assignStr += keyVal + "\n"
		} else {
			if varName == assignKey {
				if assignKey == "AsNumbers" {
					keyVal = processNumList(varValue.([]int64))
				} else {
					keyVal = processStringList(varValue.([]string))
				}
			}
			if varType.Kind() == reflect.String {
				assignStr += tab8 + key + " = \"" + keyVal + "\"\n"
			} else {
				assignStr += tab8 + key + " = " + keyVal + "\n"
			}
		}
	}
	assignStr += tab4 + "}\n"

	return assignStr

}
