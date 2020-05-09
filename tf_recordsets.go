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
	dns "github.com/akamai/AkamaiOPEN-edgegrid-golang/configdns-v2"
	"github.com/shirou/gopsutil/mem"
	"strings"
)

const (
	MinUint uint = 0
	MaxUint      = ^MinUint
	MaxInt       = int(MaxUint >> 1)
)

// module preamble
var dnsRecConfigP1 = fmt.Sprintf(`variable "zonename" {
    description = "zone name for this name record set config"
}

`)

var dnsModuleConfig3 = fmt.Sprintf(`"

    zonename = local.zone
`)

// recordset
var dnsRecordsetConfigP1 = fmt.Sprintf(`
resource "akamai_dns_record" `)

//
// misc
var dnsRConfigP2 = fmt.Sprintf(`    zone = local.zone
`)

// Process recordset resources
func processRecordsets(zone string, resourceZoneName string, zoneTypeMap map[string]map[string]bool, fetchconfig fetchConfigStruct) (map[string]Types, error) {

	var configuredMap = make(map[string]Types, 0) // returned variable

	v, _ := mem.VirtualMemory()
	maxPageSize := (v.Free / 2) / 512 // use max half of free memory. Assume avg recordset size is 512 bytes
	if maxPageSize > uint64(MaxInt/512) {
		maxPageSize = uint64(MaxInt / 512)
	}
	pagesize := int(maxPageSize)

	// get recordsets
	queryArgs := dns.RecordsetQueryArgs{PageSize: pagesize, SortBy: "name, type", Page: 1}
	nameRecordSetsResp, err := dns.GetRecordsets(zone, queryArgs)
	if err != nil {
		return configuredMap, fmt.Errorf("Failed to read record set. %s", err.Error())
	}
	var recordFields map[string]interface{}
	for {
		recordBody := ""
		rString := ""
		listString := ""
		modstring := ""
		tfModule := ""
		if fetchConfig.ConfigOnly {
			// can specify record names with config only
			for _, recname := range recordNames {
				zoneTypeMap[recname] = map[string]bool{}
			}
		}
		for _, rs := range nameRecordSetsResp.Recordsets {
			if fetchConfig.ConfigOnly {
				// combination of recordnames and config only valid
				if len(recordNames) > 0 {
					if _, ok := zoneTypeMap[rs.Name]; !ok {
						continue
					}
				}
			} else {
				if _, ok := zoneTypeMap[rs.Name]; !ok {
					continue
				}
				if !fetchConfig.NamesOnly && !zoneTypeMap[rs.Name][rs.Type] {
					continue
				}
			}
			// update configuredMap
			if _, ok := configuredMap[rs.Name]; !ok {
				configuredMap[rs.Name] = Types{}
			}
			configuredMap[rs.Name] = append(configuredMap[rs.Name], rs.Type)

			recordFields = dns.ParseRData(rs.Type, rs.Rdata) //returns map[string]interface{}
			// required fields
			recordFields["name"] = rs.Name
			//recordFields["active"] = true                   // how set?
			recordFields["recordtype"] = rs.Type
			recordFields["ttl"] = rs.TTL
			recordBody = ""
			strval := ""
			rString = dnsRecordsetConfigP1
			for fname, fval := range recordFields {
				if (fname == "priority" || fname == "priority_increment") && rs.Type == "MX" {
					fval = 0
				}
				if rs.Type == "SOA" && fname == "serial" {
					continue // computed
				}
				if rs.Type == "AKAMAITLC" && (fname == "dns_name" || fname == "answer_type") {
					continue // computed
				}
				recordBody += tab4 + fname + " = "
				switch fval.(type) {
				case string:
					strval = fmt.Sprint(fval)
					if strings.HasPrefix(strval, "\"") {
						strval = strings.Trim(strval, "\"")
						strval = "\\\"" + strval + "\\\""
					}
					recordBody += "\"" + strval + "\"\n"

				case []string:
					// target
					listString = ""
					if len(fval.([]string)) > 0 {
						listString += "["
						if rs.Type == "MX" {
							for _, rstr := range rs.Rdata {
								listString += "\"" + rstr + "\""
								listString += ", "
							}
						} else {
							for _, str := range fval.([]string) {
								if strings.HasPrefix(str, "\"") {
									str = strings.Trim(str, "\"")
									if strings.Contains(str, "\\\"") {
										strings.ReplaceAll(str, "\\\"", "\\\\\\\"")
									}
									str = "\\\"" + str + "\\\""
								}
								listString += "\"" + str + "\""
								listString += ", "
							}
						}
						listString = strings.TrimRight(listString, ", ")
						listString += "]"
					} else {
						listString += "[]"
					}
					recordBody += fmt.Sprint(listString) + "\n"

				default:
					recordBody += fmt.Sprint(fval) + "\n"
				}
			}
			rString += "\"" + createRecordsetNormalName(resourceZoneName, rs.Name, rs.Type) + "\" {\n"
			rString += dnsRConfigP2
			rString += recordBody
			rString += "}\n"
			if fetchConfig.ModSegment {
				// process as module
				modName := createRecordsetNormalName(resourceZoneName, rs.Name, rs.Type)
				tfModule = dnsModuleConfig1 + modName
				tfModule += dnsModuleConfig2 + createNamedModulePath(modName)
				tfModule += dnsModuleConfig3
				tfModule += "}\n"
				appendRootModuleTF(tfModule)
				modstring = dnsRecConfigP1
				modstring += dnsModZoneConfigP1 + "var.zonename\n" + "}\n"
				modstring += rString
				createModuleTF(modName, modstring)
			} else {
				// add to toplevel TF
				appendRootModuleTF(rString)
			}
		}

		if nameRecordSetsResp.Metadata.Page == nameRecordSetsResp.Metadata.LastPage || nameRecordSetsResp.Metadata.LastPage == 0 {
			break
		}
		queryArgs.Page += 1
		nameRecordSetsResp, err = dns.GetRecordsets(zone, queryArgs)
		if err != nil {
			return configuredMap, fmt.Errorf("Failed to read record set. %s", err.Error())
		}
	}

	return configuredMap, nil

}

// create unique resource record name
func createRecordsetNormalName(resourceZoneName, rName, rType string) string {

	return strings.TrimRight(fmt.Sprintf("%s_%s_%s",
		normalizeResourceName(resourceZoneName),
		normalizeResourceName(rName),
		rType), "_")

}
