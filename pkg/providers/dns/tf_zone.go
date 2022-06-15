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
	"reflect"

	dns "github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/configdns"
)

// Keys that can be ignored, e.g. lists, read-only, don't want
var ignoredZoneKeys = map[string]int{"LastActivationDate": 0, "LastModifiedBy": 0, "LastModifiedDate": 0,
	"AliasCount": 0, "ActivationState": 0, "VersionId": 0}

// process zone
func processZone(ctx context.Context, zone *dns.ZoneResponse, resourceZoneName string, modSegment bool, fileUtils fileUtils) (string, error) {
	data := Data{
		Zone:           zone.Zone,
		BlockName:      resourceZoneName,
		ResourceFields: gatherResourceFields(zone),
	}
	var zoneTF string
	if modSegment {
		err := fileUtils.createModuleTF(ctx, resourceZoneName, useTemplate(&data, "config.tmpl", true))
		if err != nil {
			return "", err
		}
		zoneTF = useTemplate(&data, "zone.tmpl", true)
	} else {
		zoneTF = useTemplate(&data, "full_zone.tmpl", true)
	}

	return zoneTF, nil

}

func gatherResourceFields(zone *dns.ZoneResponse) map[string]string {
	resourceFields := make(map[string]string)
	zoneElems := reflect.ValueOf(zone).Elem()
	for i := 0; i < zoneElems.NumField(); i++ {
		varName := zoneElems.Type().Field(i).Name
		varType := zoneElems.Type().Field(i).Type
		varValue := zoneElems.Field(i).Interface()
		keyVal := fmt.Sprint(varValue)
		if _, ok := ignoredZoneKeys[varName]; ok {
			continue
		}
		key := convertKey(varName, keyVal, varType.Kind())
		if key == "" {
			continue
		}
		if varName == "ContractID" {
			continue
		} else if varName == "Zone" {
			keyVal = "local.zone"
		} else if varName == "Masters" {
			keyVal = processStringList(zone.Masters)
		} else if varName == "TsigKey" {
			if varValue.(*dns.TSIGKey) == nil {
				continue
			}
			keyVal = processTsigKey(varValue.(*dns.TSIGKey))
		}

		if varName != "Zone" && varType.Kind() == reflect.String {
			resourceFields[key] = fmt.Sprintf("%q", keyVal)
		} else {
			resourceFields[key] = keyVal
		}
	}
	return resourceFields
}

func processTsigKey(key *dns.TSIGKey) string {

	keyBody := ""
	if len(key.Name) == 0 {
		// Name has to be not empty if legit
		return keyBody
	}
	keyBody += fmt.Sprintf(" {\n")
	keyBody += fmt.Sprintf("%sname = %q\n", tab8, key.Name)
	keyBody += fmt.Sprintf("%salgorithm = %q\n", tab8, key.Algorithm)
	keyBody += fmt.Sprintf("%ssecret = %q\n", tab8, key.Secret)
	keyBody += fmt.Sprintf("%s}", tab8)
	return keyBody
}
