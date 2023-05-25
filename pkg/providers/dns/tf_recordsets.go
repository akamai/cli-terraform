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

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v6/pkg/dns"
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
func processRecordsets(ctx context.Context, client dns.DNS, zone string, resourceZoneName string, zoneTypeMap map[string]map[string]bool, fileUtils fileUtils, config configStruct) (map[string]Types, error) {

	// returned variable. That map later will be used to create import script
	var importScriptConfig = make(map[string]Types)

	queryArgs := getQueryArguments()
	nameRecordSetsResp, err := client.GetRecordsets(ctx, zone, queryArgs)
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
		for _, recordset := range nameRecordSetsResp.Recordsets {
			if !shouldProcessRecordset(zoneTypeMap, recordset, config) {
				continue
			}
			updateImportScriptConfig(importScriptConfig, recordset)

			recordMap := getRecordMap(ctx, client, recordset)
			modName := createUniqueRecordsetName(resourceZoneName, recordset.Name, recordset.Type)
			data := RecordsetData{BlockName: modName, ResourceFields: recordMap, TfWorkPath: config.tfWorkPath}
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
		nameRecordSetsResp, err = client.GetRecordsets(ctx, zone, queryArgs)
		if err != nil {
			return importScriptConfig, fmt.Errorf("failed to read record set %s", err.Error())
		}
	}

	return importScriptConfig, nil

}

func updateImportScriptConfig(importScriptConfig map[string]Types, recordset dns.Recordset) {
	if _, ok := importScriptConfig[recordset.Name]; !ok {
		importScriptConfig[recordset.Name] = Types{}
	}
	importScriptConfig[recordset.Name] = append(importScriptConfig[recordset.Name], recordset.Type)
}

func shouldProcessRecordset(zoneTypeMap map[string]map[string]bool, recordset dns.Recordset, config configStruct) bool {
	if config.fetchConfig.ConfigOnly {
		// combination of recordnames and config only valid
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

func getQueryArguments() dns.RecordsetQueryArgs {
	v, _ := mem.VirtualMemory()
	maxPageSize := (v.Free / 2) / 512 // use max half of free memory. Assume avg recordset size is 512 bytes
	if maxPageSize > uint64(maxInt/512) {
		maxPageSize = uint64(maxInt / 512)
	}
	pagesize := int(maxPageSize)

	// get recordsets
	queryArgs := dns.RecordsetQueryArgs{PageSize: pagesize, SortBy: "name, type", Page: 1}
	return queryArgs
}

// getRecordMap returns all fields that will be exported into generated resource. The fields name bases on recordset type
func getRecordMap(ctx context.Context, client dns.DNS, recordset dns.Recordset) map[string]string {
	// keys of that map depends on recordset.Type
	recordFields := client.ParseRData(ctx, recordset.Type, recordset.Rdata) //returns map[string]interface{}
	// required fields
	recordFields["name"] = recordset.Name
	//recordFields["active"] = true                   // how set?
	recordFields["recordtype"] = recordset.Type
	recordFields["ttl"] = recordset.TTL
	recordMap := make(map[string]string)
	for fname, fval := range recordFields {
		if (fname == "priority" || fname == "priority_increment") && recordset.Type == "MX" {
			fval = 0
		}
		if recordset.Type == "SOA" && fname == "serial" {
			continue // computed
		}
		if recordset.Type == "AKAMAITLC" && (fname == "dns_name" || fname == "answer_type") {
			continue // computed
		}
		if fname == "svc_params" && (recordset.Type == "SVCB" || recordset.Type == "HTTPS") {
			if createParamsMap(splitSvcParams(fmt.Sprint(fval))) == nil {
				continue
			}
		}
		switch fval.(type) {
		case string:
			recordMap[fname] = "\"" + handleSpecialCharacters(recordset, fval) + "\""

		case []string:
			// target
			recordMap[fname] = fmt.Sprint(recordValueForSlice(fval, recordset))
		default:
			recordMap[fname] = fmt.Sprint(fval)
		}
	}
	return recordMap
}

func handleSpecialCharacters(rs dns.Recordset, fval interface{}) string {
	strval := fmt.Sprint(fval)
	if rs.Type == "HTTPS" || rs.Type == "SVCB" {
		strval = strings.ReplaceAll(strval, "\"", "\\\"")
	}
	if strings.HasPrefix(strval, "\"") {
		strval = strings.Trim(strval, "\"")
		strval = "\\\"" + strval + "\\\""
	}
	return strval
}

func recordValueForSlice(fval interface{}, rs dns.Recordset) string {
	listString := ""
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
	return listString
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
func createUniqueRecordsetName(resourceZoneName, rName, rType string) string {

	return strings.TrimRight(fmt.Sprintf("%s_%s_%s",
		normalizeResourceName(resourceZoneName),
		normalizeResourceName(rName),
		rType), "_")

}
