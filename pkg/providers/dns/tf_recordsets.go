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

package dns

import (
	"fmt"
	"regexp"
	"strings"

	dns "github.com/akamai/AkamaiOPEN-edgegrid-golang/configdns-v2"
	"github.com/shirou/gopsutil/mem"
)

const (
	minUint uint = 0
	maxUint      = ^minUint
	maxInt       = int(maxUint >> 1)
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

// Util func to split params string
func splitSvcParams(params string) []string {

	r := regexp.MustCompile(`[^\s"]+|"([^"]*)"`)
	paramslice := r.FindAllString(params, -1)
	return paramslice
}

// Util func to walk params and create a map
func createParamsMap(params []string) *map[string]string {

	paramsMap := map[string]string{}

	for _, param := range params {
		keyval := strings.Split(param, "=")
		if len(keyval) == 0 || len(keyval) > 2 {
			continue // weird but skip
		}
		if len(keyval) == 1 {
			paramsMap[strings.TrimSpace(keyval[0])] = "" // no value
			continue
		}
		paramsMap[strings.TrimSpace(keyval[0])] = strings.TrimSpace(keyval[1])
	}
	if len(paramsMap) < 1 {
		return nil
	}
	return &paramsMap

}

// Process recordset resources
func processRecordsets(zone string, resourceZoneName string, zoneTypeMap map[string]map[string]bool, _ fetchConfigStruct) (map[string]Types, error) {

	var configuredMap = make(map[string]Types) // returned variable

	v, _ := mem.VirtualMemory()
	maxPageSize := (v.Free / 2) / 512 // use max half of free memory. Assume avg recordset size is 512 bytes
	if maxPageSize > uint64(maxInt/512) {
		maxPageSize = uint64(maxInt / 512)
	}
	pagesize := int(maxPageSize)

	// get recordsets
	queryArgs := dns.RecordsetQueryArgs{PageSize: pagesize, SortBy: "name, type", Page: 1}
	nameRecordSetsResp, err := dns.GetRecordsets(zone, queryArgs)
	if err != nil {
		return configuredMap, fmt.Errorf("failed to read record set %s", err.Error())
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
				var paramsMap *map[string]string
				if fname == "svc_params" && (rs.Type == "SVCB" || rs.Type == "HTTPS") {
					paramsMap = createParamsMap(splitSvcParams(fmt.Sprint(fval)))
					if paramsMap == nil {
						continue
					}
				}
				recordBody += tab4 + fname + " = "
				switch fval.(type) {
				case string:
					strval = fmt.Sprint(fval)
					if rs.Type == "HTTPS" || rs.Type == "SVCB" {
						strval = strings.ReplaceAll(strval, "\"", "\\\"")
					}
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
						} else if rs.Type == "CAA" {
							for _, rstr := range rs.Rdata {
								caaparts := strings.Split(rstr, " ")
								caaparts[2] = strings.ReplaceAll(caaparts[2], "\"", "\\\"")
								listString += "\"" + strings.Join(caaparts, " ") + "\""
								listString += ", "
							}
						} else {
							for _, str := range fval.([]string) {
								if strings.HasPrefix(str, "\"") {
									str = strings.Trim(str, "\"")
									str = processString(str)
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
				if err := appendRootModuleTF(tfModule); err != nil {
					return nil, err
				}
				modstring = dnsRecConfigP1
				modstring += dnsModZoneConfigP1 + "var.zonename\n" + "}\n"
				modstring += rString
				if err := createModuleTF(modName, modstring); err != nil {
					return nil, err
				}
			} else {
				// add to toplevel TF
				if err := appendRootModuleTF(rString); err != nil {
					return nil, err
				}
			}
		}

		if nameRecordSetsResp.Metadata.Page == nameRecordSetsResp.Metadata.LastPage || nameRecordSetsResp.Metadata.LastPage == 0 {
			break
		}
		queryArgs.Page++
		nameRecordSetsResp, err = dns.GetRecordsets(zone, queryArgs)
		if err != nil {
			return configuredMap, fmt.Errorf("failed to read record set %s", err.Error())
		}
	}

	return configuredMap, nil

}

// process string with embedded quotes
func processString(source string) string {

	if len(source) < 1 {
		return source
	}
	stringSlice := strings.Split(source, " ")
	if len(stringSlice) == 1 {
		return stringSlice[0]
	}
	cleanString := ""
	workingString := source
	var quoteIndex int
	for {
		quoteIndex = strings.Index(workingString, "\"")
		switch quoteIndex {
		case -1:
			cleanString += workingString
			return cleanString
		case 0:
			cleanString += "\\\""
		case len(workingString) - 1:
			if workingString[quoteIndex-1:quoteIndex] != "\\" {
				cleanString += workingString[:quoteIndex]
				cleanString += "\\\""
			} else {
				cleanString += workingString[:]
			}
			return cleanString
		default:
			if workingString[quoteIndex-1:quoteIndex] != "\\" {
				cleanString += workingString[:quoteIndex]
				cleanString += "\\\""
			} else {
				cleanString += workingString[:quoteIndex+1]
			}
		}
		workingString = workingString[quoteIndex+1:]
	}
}

// create unique resource record name
func createRecordsetNormalName(resourceZoneName, rName, rType string) string {

	return strings.TrimRight(fmt.Sprintf("%s_%s_%s",
		normalizeResourceName(resourceZoneName),
		normalizeResourceName(rName),
		rType), "_")

}
