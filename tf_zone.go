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
	"reflect"
)

// Keys that can be ignored, e.g. lists, read-only, don't want
var ignoredZoneKeys map[string]int = map[string]int{"LastActivationDate": 0, "LastModifiedBy": 0, "LastModifiedDate": 0,
	"AliasCount": 0, "ActivationState": 0, "VersionId": 0}

// header, zone
var tfHeaderContent = fmt.Sprint(`terraform {
  required_version = ">= 0.13"
  required_providers {
    akamai = {
      source = "akamai/akamai"
      version = "~> 1.6.1"
    }
  }
}

`)

var dnsZoneResourceConfig = fmt.Sprintf(`resource "akamai_dns_zone" "`)

var dnsZoneConfigP2 = fmt.Sprintf(`    contract = var.contractid
    group = var.groupid
`)

var dnsZoneConfigP3 = fmt.Sprintf(`    comment =  "Zone import"
`)

// module
var dnsModZoneConfigP1 = fmt.Sprintf(`locals {
    zone = `)

var dnsModZoneConfigP2 = fmt.Sprintf(`variable "contract" {
  description = "contract id for zone creation"
}

variable "group" {
  description = "group id for zone creation"
}

variable "name" {
  description = "zone name"
}

output "zonename" {
	value = akamai_dns_zone.`)

// process zone
func processZone(zone *dns.ZoneResponse, resourceZoneName string, modSegment bool) (string, error) {

	zoneBody := ""
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
		if varName == "ContractId" {
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
		zoneBody += tab4 + key + " = "
		if varName != "Zone" && varType.Kind() == reflect.String {
			zoneBody += "\"" + keyVal + "\"\n"
		} else {
			zoneBody += keyVal + "\n"
		}
	}
	zoneResourceString := dnsZoneResourceConfig + resourceZoneName + "\" {\n"
	zoneResourceString += dnsZoneConfigP2
	zoneResourceString += zoneBody
	zoneResourceString += "}\n\n"
	zoneTF := tfHeaderContent
	if modSegment {
		// create initial TF preamble and zone  module
		zoneTF += dnsModZoneConfigP1 + "\"" + zone.Zone + "\"\n" + "}\n\n"
		zoneTF += dnsModuleConfig1 + resourceZoneName
		zoneTF += dnsModuleConfig2 + createNamedModulePath(resourceZoneName) + "\"\n\n"
		zoneTF += dnsZoneConfigP2
		zoneTF += "    name = local.zone\n"
		zoneTF += "}\n\n"
		// create module config
		modConfig := dnsModZoneConfigP2 + resourceZoneName + ".name\n"
		modConfig += "}\n\n"
		modConfig += dnsModZoneConfigP1 + "var.name\n" + "}\n"
		modConfig += zoneResourceString
		err := createModuleTF(resourceZoneName, modConfig)
		if err != nil {
			return "", err
		}
	} else {
		// All in one config file
		zoneTF += dnsModZoneConfigP1 + "\"" + zone.Zone + "\"\n}\n\n"
		zoneTF += zoneResourceString
	}

	return zoneTF, nil

}

func processTsigKey(key *dns.TSIGKey) string {

	keyBody := ""
	if len(key.Name) == 0 {
		// Name has to be not empty if legit
		return keyBody
	}
	keyBody += " {\n"
	keyBody += tab8 + "name = \"" + key.Name + "\"\n"
	keyBody += tab8 + "algorithm = \"" + key.Algorithm + "\"\n"
	keyBody += tab8 + "secret = \"" + key.Secret + "\"\n"
	keyBody += tab8 + "}"
	return keyBody
}
