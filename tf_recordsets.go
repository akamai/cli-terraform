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
	"sort"
	"strconv"
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
func processRecordsets(importStruct *zoneImportListStruct, resourceZoneName string, zoneTypeMap map[string]map[string]bool, modsegment bool) error {

	var totalRecordsets int = 0
	// total desired recordsets
	for _, typeMap := range importStruct.Recordsets {
		totalRecordsets += len(typeMap)
	}

	v, _ := mem.VirtualMemory()
	maxPageSize := (v.Free / 2) / 512 // use max half of free memory. Assume avg recordset size is 512 bytes
	if maxPageSize > uint64(MaxInt/512) {
		maxPageSize = uint64(MaxInt / 512)
	}
	pagesize := int(maxPageSize)
	if totalRecordsets < int(maxPageSize) {
		pagesize = totalRecordsets
	}
	pagesize = 1

	// get recordsets
	queryArgs := dns.RecordsetQueryArgs{PageSize: pagesize, SortBy: "name, type", Page: 1}
	nameRecordSetsResp, err := dns.GetRecordsets(importStruct.Zone, queryArgs)
	if err != nil {
		return fmt.Errorf("Failed to read record set. %s", err.Error())
	}
	for {
		for _, rs := range nameRecordSetsResp.Recordsets {
			//for each record set ...
			if _, ok := zoneTypeMap[rs.Name]; !ok {
				continue
			}
			if !zoneTypeMap[rs.Name][rs.Type] {
				continue
			}
			recordFields := parseRData(rs.Type, rs.Rdata) //returns map[string]interface{}
			// required fields
			recordFields["name"] = rs.Name
			//recordFields["active"] = true                   // how set?
			recordFields["recordtype"] = rs.Type
			recordFields["ttl"] = rs.TTL
			recordBody := ""
			rString := dnsRecordsetConfigP1
			for fname, fval := range recordFields {
				recordBody += tab4 + fname + " = "
				switch fval.(type) {
				case string:
					recordBody += "\"" + fmt.Sprint(fval) + "\"\n"

				case []string:
					// target
					listString := ""
					if len(fval.([]string)) > 0 {
						listString += "["
						for _, str := range fval.([]string) {
							listString += str + ", "
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
			if modsegment {
				// process as module
				modName := createRecordsetNormalName(resourceZoneName, rs.Name, rs.Type)
				tfModule := dnsModuleConfig1 + modName
				tfModule += dnsModuleConfig2 + createNamedModulePath(modName)
				tfModule += dnsModuleConfig3
				tfModule += "}\n"
				appendRootModuleTF(tfModule)
				modstring := dnsRecConfigP1
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
		nameRecordSetsResp, err = dns.GetRecordsets(importStruct.Zone, queryArgs)
		if err != nil {
			return fmt.Errorf("Failed to read record set. %s", err.Error())
		}
	}

	return nil

}

// create unique resource record name
func createRecordsetNormalName(resourceZoneName, rName, rType string) string {

	return strings.TrimRight(fmt.Sprintf("%s_%s_%s",
		normalizeResourceName(resourceZoneName),
		normalizeResourceName(rName),
		rType), "_")

}

// parse RData in context of type. Return map of fields and values
func parseRData(rtype string, rdata []string) map[string]interface{} {

	fieldMap := make(map[string]interface{}, 0)
	//rdataFields := "" //rings.Split(rdata, " ")
	if len(rdata) == 0 {
		return fieldMap
	}
	newrdata := make([]string, 0, len(rdata))

	switch rtype {
	case "AFSDB":
		parts := strings.Split(rdata[0], " ")
		fieldMap["subtype"], _ = strconv.Atoi(parts[0])
		for _, rcontent := range rdata {
			parts := strings.Split(rcontent, " ")
			newrdata = append(newrdata, "\""+parts[1]+"\"")
		}
		fieldMap["target"] = newrdata

	case "DNSKEY":
		for _, rcontent := range rdata {
			parts := strings.Split(rcontent, " ")
			fieldMap["flags"], _ = strconv.Atoi(parts[0])
			fieldMap["protocol"], _ = strconv.Atoi(parts[1])
			fieldMap["algorithm"], _ = strconv.Atoi(parts[2])
			fieldMap["key"] = parts[3]
			break
		}

	case "DS":
		for _, rcontent := range rdata {
			parts := strings.Split(rcontent, " ")
			fieldMap["keytag"], _ = strconv.Atoi(parts[0])
			fieldMap["digest_type"], _ = strconv.Atoi(parts[1])
			fieldMap["algorithm"], _ = strconv.Atoi(parts[2])
			fieldMap["digest"] = parts[3]
			break
		}

	case "HINFO":
		for _, rcontent := range rdata {
			parts := strings.Split(rcontent, " ")
			fieldMap["hardware"] = parts[0]
			fieldMap["software"] = parts[1]
			break
		}

	case "MX":
		sort.Strings(rdata)
		parts := strings.Split(rdata[0], " ")
		fieldMap["priority"], _ = strconv.Atoi(parts[0])
		if len(rdata) > 1 {
			parts = strings.Split(rdata[1], " ")
			tpri, _ := strconv.Atoi(parts[0])
			fieldMap["priority_increment"] = tpri - fieldMap["priority"].(int)
		}
		for _, rcontent := range rdata {
			parts := strings.Split(rcontent, " ")
			newrdata = append(newrdata, "\""+parts[1]+"\"")
		}
		sort.Strings(newrdata)
		fieldMap["target"] = newrdata

	case "NAPTR":
		for _, rcontent := range rdata {
			parts := strings.Split(rcontent, " ")
			fieldMap["order"], _ = strconv.Atoi(parts[0])
			fieldMap["preference"], _ = strconv.Atoi(parts[1])
			fieldMap["flagsnaptr"] = parts[2]
			fieldMap["regexp"] = parts[3]
			fieldMap["replacement"] = parts[4]
			fieldMap["service"] = parts[5]
			break
		}

	case "NSEC3":
		for _, rcontent := range rdata {
			parts := strings.Split(rcontent, " ")
			fieldMap["flags"], _ = strconv.Atoi(parts[0])
			fieldMap["algorithm"], _ = strconv.Atoi(parts[1])
			fieldMap["iterations"], _ = strconv.Atoi(parts[2])
			fieldMap["salt"] = parts[3]
			fieldMap["next_hashed_owner_name"] = parts[4]
			fieldMap["type_bitmaps"] = parts[5]
			break
		}

	case "NSEC3PARAM":
		for _, rcontent := range rdata {
			parts := strings.Split(rcontent, " ")
			fieldMap["flags"], _ = strconv.Atoi(parts[0])
			fieldMap["algorithm"], _ = strconv.Atoi(parts[1])
			fieldMap["iterations"], _ = strconv.Atoi(parts[2])
			fieldMap["salt"] = parts[3]
			break
		}

	case "RP":
		for _, rcontent := range rdata {
			parts := strings.Split(rcontent, " ")
			fieldMap["mailbox"] = parts[0]
			fieldMap["txt"] = parts[1]
			break
		}

	case "RRSIG":
		for _, rcontent := range rdata {
			parts := strings.Split(rcontent, " ")
			fieldMap["type_covered"] = parts[0]
			fieldMap["algorithm"], _ = strconv.Atoi(parts[1])
			fieldMap["labels"], _ = strconv.Atoi(parts[2])
			fieldMap["original_ttl"], _ = strconv.Atoi(parts[3])
			fieldMap["expiration"] = parts[4]
			fieldMap["inception"] = parts[5]
			fieldMap["signature"] = parts[6]
			fieldMap["signer"] = parts[7]
			fieldMap["keytag"], _ = strconv.Atoi(parts[8])
			break
		}

	case "SRV":
		// pull out some fields
		parts := strings.Split(rdata[0], " ")
		fieldMap["priority"], _ = strconv.Atoi(parts[0])
		fieldMap["weight"], _ = strconv.Atoi(parts[1])
		fieldMap["port"], _ = strconv.Atoi(parts[2])
		// populate target
		for _, rcontent := range rdata {
			parts := strings.Split(rcontent, " ")
			newrdata = append(newrdata, "\""+parts[3]+"\"")
		}
		fieldMap["target"] = newrdata

	case "SSHFP":
		for _, rcontent := range rdata {
			parts := strings.Split(rcontent, " ")
			fieldMap["algorithm"], _ = strconv.Atoi(parts[0])
			fieldMap["fingerprint_type"], _ = strconv.Atoi(parts[1])
			fieldMap["fingerprint"] = parts[2]
			break
		}

	case "SOA":
		for _, rcontent := range rdata {
			parts := strings.Split(rcontent, " ")
			fieldMap["name_server"] = parts[0]
			fieldMap["email_address"] = parts[1]
			fieldMap["serial"], _ = strconv.Atoi(parts[2])
			fieldMap["refresh"], _ = strconv.Atoi(parts[3])
			fieldMap["retry"], _ = strconv.Atoi(parts[4])
			fieldMap["expiry"], _ = strconv.Atoi(parts[5])
			fieldMap["nxdomain_ttl"], _ = strconv.Atoi(parts[6])
			break
		}

	case "AKAMAICDN":
		fieldMap["edge_hostname"] = rdata[0]

	case "AKAMAITLC":
		parts := strings.Split(rdata[0], " ")
		fieldMap["answer_type"] = parts[0]
		fieldMap["dns_name"] = parts[1]

	default:
		for _, rcontent := range rdata {
			newrdata = append(newrdata, "\""+rcontent+"\"")
		}
		fieldMap["target"] = newrdata
	}

	return fieldMap

}
