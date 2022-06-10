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
	"unicode"
)

// Keys that can be ignored, e.g. lists, read-only, don't want
var ignoredKeys = map[string]int{"AsMaps": 0, "Resources": 0, "Properties": 0, "Datacenters": 0, "CidrMaps": 0, "GeographicMaps": 0,
	"Links": 0, "Status": 0, "DefaultUnreachableThreshold": 0, "MinPingableRegionFraction": 0,
	"ServermonitorLivenessCount": 0, "RoundRobinPrefix": 0, "ServermonitorLoadCount": 0,
	"PingInterval": 0, "MaxTTL": 0, "DefaultHealthMax": 0, "MapUpdateInterval": 0, "MaxProperties": 0,
	"MaxResources": 0, "MaxTestTimeout": 0, "DefaultHealthMultiplier": 0, "ServermonitorPool": 0,
	"MinTTL": 0, "DefaultMaxUnreachablePenalty": 0, "DefaultHealthThreshold": 0, "MinTestInterval": 0,
	"PingPacketSize": 0, "ScorePenalty": 0, "LastModified": 0, "LastModifiedBy": 0, "ModificationComments": 0,
	"WeightedHashBitsForIPv4": 0, "WeightedHashBitsForIPv6": 0, "Virtual": 0}

// initialized with key names that don't follow mapping pattern. populated in convert key for secondary encounters
var mappedKeys = map[string]string{"DynamicTTL": "dynamic_ttl", "StaticTTL": "static_ttl", "StaticRRSets": "static_rr_sets",
	"TTL": "ttl", "DatacenterId": "datacenter_id", "HandoutCName": "handout_cname", "StickinessBonusPercentage": "stickiness_bonus_percentage",
	"CName": "cname", "BackupCName": "backup_cname"}

var tab4 = "    "
var tab8 = "        "
var tab12 = "            "

// misc
var gtmRConfigP2 = fmt.Sprintf(`    domain = akamai_gtm_domain.`)

var dependsClauseP1 = fmt.Sprintf(`    depends_on = [
        akamai_gtm_domain.`)

// utility method to process string lists
func processStringList(sl []string) string {

	switch len(sl) {
	case 0:
		return "[]"
	case 1:
		return "[\"" + sl[0] + "\"]"
	default:
		keyVal := "["
		for i, sval := range sl {
			keyVal += "\"" + sval + "\""
			if i != len(sl)-1 {
				keyVal += ", "
			}
		}
		keyVal += "]"
		return keyVal
	}

}

// utility method to process int64 lists
func processNumList(sl []int64) string {

	switch len(sl) {
	case 0:
		return "[]"
	case 1:
		return "[" + strconv.FormatInt(sl[0], 10) + "]"
	default:
		keyVal := "["
		for i, ival := range sl {
			keyVal += strconv.FormatInt(ival, 10)
			if i != len(sl)-1 {
				keyVal += ", "
			}
		}
		keyVal += "]"
		return keyVal
	}

}

// utility method to convert camelCased struct field names to provider field naming convention
func convertKey(inKey string, _ string, _ reflect.Kind) string {

	if _, ok := ignoredKeys[inKey]; ok {
		return ""
	}
	if val, ok := mappedKeys[inKey]; ok {
		return val
	}

	outKey := ""
	for i, char := range inKey {
		if unicode.IsUpper(char) {
			if i != 0 {
				outKey += "_"
			}
			outKey += string(unicode.ToLower(char))
		} else {
			outKey += string(char)
		}
	}
	mappedKeys[inKey] = outKey
	return outKey

}
