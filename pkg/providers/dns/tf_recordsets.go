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
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/dns"
	"github.com/akamai/cli-terraform/pkg/tools"
	"github.com/shirou/gopsutil/mem"
)

const (
	minUint uint = 0
	maxUint      = ^minUint
	maxInt       = int(maxUint >> 1)
)

// Util func to split params string
func splitSvcParams(params string) []string {
	r := regexp.MustCompile(`[^\s"]+|"([^"]*)"`)
	paramSlice := r.FindAllString(params, -1)
	return paramSlice
}

// Util func to walk params and create a map
func createParamsMap(params []string) *map[string]string {
	paramsMap := map[string]string{}

	for _, param := range params {
		keyVal := strings.Split(param, "=")
		if len(keyVal) == 0 || len(keyVal) > 2 {
			continue // weird but skip
		}
		if len(keyVal) == 1 {
			paramsMap[strings.TrimSpace(keyVal[0])] = "" // no value
			continue
		}
		paramsMap[strings.TrimSpace(keyVal[0])] = strings.TrimSpace(keyVal[1])
	}
	if len(paramsMap) < 1 {
		return nil
	}

	return &paramsMap
}

// Process recordset resources
func processRecordSets(ctx context.Context, client dns.DNS, zone, resourceZoneName string, zoneTypeMap map[string]map[string]bool, fileUtils fileUtils, config configStruct) (map[string]Types, error) {
	// returned variable. That map later will be used to create import script
	var importScriptConfig = make(map[string]Types)

	queryArgs := getQueryArguments()
	nameRecordSetsResp, err := client.GetRecordSets(ctx, dns.GetRecordSetsRequest{
		Zone:      zone,
		QueryArgs: &queryArgs,
	})
	if err != nil {
		return importScriptConfig, fmt.Errorf("failed to read record set %s", err.Error())
	}
	for {
		if config.fetchConfig.ConfigOnly {
			// can specify record names with config only
			for _, recname := range config.recordNames {
				zoneTypeMap[recname] = map[string]bool{}
			}
		}
		for _, recordset := range nameRecordSetsResp.RecordSets {
			if !shouldProcessRecordset(zoneTypeMap, recordset, config) {
				continue
			}
			updateImportScriptConfig(importScriptConfig, recordset)

			recordMap := getRecordMap(ctx, client, recordset)
			modName := createUniqueRecordsetName(resourceZoneName, recordset.Name, recordset.Type)
			data := RecordsetData{BlockName: modName, ResourceFields: recordMap, TFWorkPath: config.tfWorkPath}
			if config.fetchConfig.ModSegment {
				// process as module
				if err := fileUtils.appendRootModuleTF(useTemplate(&data, "module-set.tmpl", false)); err != nil {
					return nil, err
				}
				if err := fileUtils.createModuleTF(ctx, modName, useTemplate(&data, "recordset-modsegment.tmpl", true), config.tfWorkPath); err != nil {
					return nil, err
				}
			} else {
				// add to toplevel TF
				if err := fileUtils.appendRootModuleTF(useTemplate(&data, "resource-set.tmpl", false)); err != nil {
					return nil, err
				}
			}
		}

		if nameRecordSetsResp.Metadata.Page == nameRecordSetsResp.Metadata.LastPage || nameRecordSetsResp.Metadata.LastPage == 0 {
			break
		}
		queryArgs.Page++
		nameRecordSetsResp, err = client.GetRecordSets(ctx, dns.GetRecordSetsRequest{
			Zone:      zone,
			QueryArgs: &queryArgs,
		})
		if err != nil {
			return importScriptConfig, fmt.Errorf("failed to read record set %s", err.Error())
		}
	}

	return importScriptConfig, nil
}

func updateImportScriptConfig(importScriptConfig map[string]Types, recordset dns.RecordSet) {
	if _, ok := importScriptConfig[recordset.Name]; !ok {
		importScriptConfig[recordset.Name] = Types{}
	}
	importScriptConfig[recordset.Name] = append(importScriptConfig[recordset.Name], recordset.Type)
}

func shouldProcessRecordset(zoneTypeMap map[string]map[string]bool, recordset dns.RecordSet, config configStruct) bool {
	if config.fetchConfig.ConfigOnly {
		// combination of record names and config only valid
		if len(config.recordNames) > 0 {
			if _, ok := zoneTypeMap[recordset.Name]; !ok {
				return false
			}
		}
	} else {
		if _, ok := zoneTypeMap[recordset.Name]; !ok {
			return false
		}
		if !config.fetchConfig.NamesOnly && !zoneTypeMap[recordset.Name][recordset.Type] {
			return false
		}
	}
	return true
}

func getQueryArguments() dns.RecordSetQueryArgs {
	v, _ := mem.VirtualMemory()
	maxPageSize := (v.Free / 2) / 512 // use max half of free memory. Assume avg recordset size is 512 bytes
	if maxPageSize > uint64(maxInt/512) {
		maxPageSize = uint64(maxInt / 512)
	}
	pageSize := int(maxPageSize)

	// get recordsets
	queryArgs := dns.RecordSetQueryArgs{PageSize: pageSize, SortBy: "name, type", Page: 1}
	return queryArgs
}

// getRecordMap returns all fields that will be exported into generated resource. The fields name bases on recordset type
func getRecordMap(ctx context.Context, client dns.DNS, recordset dns.RecordSet) map[string]string {
	// keys of that map depends on recordset.Type
	recordFields := client.ParseRData(ctx, recordset.Type, recordset.Rdata) //returns map[string]interface{}
	// required fields
	recordFields["name"] = recordset.Name
	recordFields["recordtype"] = recordset.Type
	recordFields["ttl"] = recordset.TTL
	recordMap := make(map[string]string)
	for fName, fVal := range recordFields {
		if (fName == "priority" || fName == "priority_increment") && recordset.Type == "MX" {
			fVal = 0
		}
		if recordset.Type == "SOA" && fName == "serial" {
			continue // computed
		}
		if recordset.Type == "AKAMAITLC" && (fName == "dns_name" || fName == "answer_type") {
			continue // computed
		}
		if fName == "svc_params" && (recordset.Type == "SVCB" || recordset.Type == "HTTPS") {
			if createParamsMap(splitSvcParams(fmt.Sprint(fVal))) == nil {
				continue
			}
		}
		switch fVal.(type) {
		case string:
			recordMap[fName] = "\"" + handleSpecialCharacters(recordset, fVal) + "\""

		case []string:
			// target
			recordMap[fName] = fmt.Sprint(recordValueForSlice(fVal, recordset))
		default:
			recordMap[fName] = fmt.Sprint(fVal)
		}
	}
	return recordMap
}

func handleSpecialCharacters(rs dns.RecordSet, fVal interface{}) string {
	strVal := fmt.Sprint(fVal)
	if rs.Type == "HTTPS" || rs.Type == "SVCB" {
		strVal = strings.ReplaceAll(strVal, "\"", "\\\"")
	}
	if strings.HasPrefix(strVal, "\"") {
		strVal = strings.Trim(strVal, "\"")
		strVal = "\\\"" + strVal + "\\\""
	}
	return strVal
}

func recordValueForSlice(fVal interface{}, rs dns.RecordSet) string {
	listString := ""
	if len(fVal.([]string)) > 0 {
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
			for _, str := range fVal.([]string) {
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
	return listString
}

// process string with embedded quotes
func processString(source string) string {
	return tools.EscapeQuotedStringLit(source)
}

// create unique resource record name
func createUniqueRecordsetName(resourceZoneName, rName, rType string) string {
	return strings.TrimRight(fmt.Sprintf("%s_%s_%s",
		normalizeResourceName(resourceZoneName),
		normalizeResourceName(rName),
		rType), "_")
}
