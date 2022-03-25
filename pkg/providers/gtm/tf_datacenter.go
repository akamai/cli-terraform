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
	"strconv"

	gtm "github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/configgtm"
)

// default datacenter
var gtmDefaultDatacenterConfigP1 = fmt.Sprintf(`
data "akamai_gtm_default_datacenter" `)

// datacenter
var gtmDatacenterConfigP1 = fmt.Sprintf(`
resource "akamai_gtm_datacenter" `)

// Process datacenter resources
func processDatacenters(datacenters []*gtm.Datacenter, dcImportList map[int]string, resourceDomainName string) string {

	// Get Null values list
	var coreFieldsNullMap map[string]string
	nullFieldsMap := getNullValuesList("Datacenters")

	datacentersString := ""
	for _, datacenter := range datacenters {
		if _, ok := dcImportList[datacenter.DatacenterId]; !ok {
			continue
		}
		ddcFlag := false // ddc?
		for _, ddc := range defaultDCs {
			if datacenter.DatacenterId == ddc {
				ddcFlag = true
			}
		}
		// Retrieve Core null fields map
		if dcNullFieldObjectMap, ok := nullFieldsMap[strconv.Itoa(datacenter.DatacenterId)]; ok {
			coreFieldsNullMap = dcNullFieldObjectMap.CoreObjectFields
		} else {
			coreFieldsNullMap = map[string]string{}
		}
		datacenterBody := ""
		name := ""
		dcString := gtmDatacenterConfigP1
		if ddcFlag {
			dcString = gtmDefaultDatacenterConfigP1
			name = "default_datacenter_" + strconv.Itoa(datacenter.DatacenterId)
			datacenterBody += tab4 + "datacenter = " + strconv.Itoa(datacenter.DatacenterId) + "\n"
		} else {
			dcElems := reflect.ValueOf(datacenter).Elem()
			for i := 0; i < dcElems.NumField(); i++ {
				varName := dcElems.Type().Field(i).Name
				varType := dcElems.Type().Field(i).Type
				varValue := dcElems.Field(i).Interface()
				if _, ok := coreFieldsNullMap[varName]; ok {
					continue
				}
				keyVal := fmt.Sprint(varValue)
				key := convertKey(varName, keyVal, varType.Kind())
				if key == "" {
					continue
				}
				if key == "nickname" {
					name = keyVal
				}
				// datacenter_id is computed. Can't be in ignore map because it exists in property config
				if varName == "DatacenterId" {
					continue
				}
				if varName == "DefaultLoadObject" {
					if varValue.(*gtm.LoadObject) == nil {
						continue
					}
					loadObject := varValue.(*gtm.LoadObject)
					// hack. If all load object fields are defaults, assume its MT
					if len(loadObject.LoadServers) == 0 && loadObject.LoadObject == "" && loadObject.LoadObjectPort == 0 {
						continue
					}
					datacenterBody += tab4 + key + " {\n"
					datacenterBody += processLoadObject(varValue.(*gtm.LoadObject))
					datacenterBody += tab4 + "}\n"
					continue
				}
				datacenterBody += tab4 + key + " = "
				if varType.Kind() == reflect.String {
					datacenterBody += "\"" + keyVal + "\"\n"
				} else {
					datacenterBody += keyVal + "\n"
				}
			}
		}
		dcString += "\"" + normalizeResourceName(name) + "\" {\n"
		dcString += gtmRConfigP2 + resourceDomainName + ".name\n"
		dcString += datacenterBody
		if !ddcFlag {
			dcString += dependsClauseP1 + resourceDomainName + "\n"
			dcString += tab4 + "]\n"
		}
		dcString += "}\n"
		datacentersString += dcString
	}

	return datacentersString

}

func processLoadObject(lo *gtm.LoadObject) string {

	loBody := ""
	loBody += tab8 + "load_object = \"" + lo.LoadObject + "\"\n"
	loBody += tab8 + "load_object_port = " + strconv.Itoa(lo.LoadObjectPort) + "\n"
	lsList := processStringList(lo.LoadServers)
	if len(lo.LoadServers) < 1 {
		lsList = "[]"
	}
	loBody += tab8 + "load_servers = " + lsList + "\n"
	return loBody

}
