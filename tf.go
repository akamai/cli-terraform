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
	"strings"
	"unicode"
	"reflect"
	"strconv"
	gtm "github.com/akamai/AkamaiOPEN-edgegrid-golang/configgtm-v1_4"
	"fmt"

)

var ignoredKeys map[string]int = map[string]int{"AsMaps":0,"Resources":0,"Properties":0,"Datacenters":0,"CidrMaps":0,"GeographicMaps":0,
				"Links":0,"Status":0,"DefaultUnreachableThreshold":0,"MinPingableRegionFraction":0, 
                                "ServermonitorLivenessCount":0,"RoundRobinPrefix":0,"ServermonitorLoadCount":0,
                                "PingInterval":0,"MaxTTL":0,"DefaultHealthMax":0,"MapUpdateInterval":0,"MaxProperties":0,
                                "MaxResources":0,"MaxTestTimeout":0,"DefaultHealthMultiplier":0,"ServermonitorPool":0,
                                "MinTTL":0,"DefaultMaxUnreachablePenalty":0,"DefaultHealthThreshold":0,"MinTestInterval":0,
                                "PingPacketSize":0,"ScorePenalty":0,"LastModified":0,"LastModifiedBy":0,"ModificationComments":0,
                                "WeightedHashBitsForIPv4":0,"WeightedHashBitsForIPv6":0}

var mappedKeys map[string]string = map[string]string{"DynamicTTL":"dynamic_ttl","StaticTTL":"static_ttl","StaticRRSets":"static_rr_sets",
				"TTL":"ttl","DatacenterId":"datacenter_id"}

var intNullableKeys map[string]int = map[string]int{"UpperBound":0,"DefaultTimeOutPenalty":0,"CloneOf":0,"Latitude":0,"Longitude":0,
							"DefaultErrorPenalty":0,"DefaultTimeoutPenalty":0,"HealthThreshold":0,
							"HealthMultiplier":0,"HealthMax":0}

var tab4 = "    "
var tab8 = "        "
var tab12 = "            "

// header, domain
var gtmHeaderConfig = fmt.Sprintf(`
provider "akamai" {
  gtm_section = "${var.gtmsection}"
}

resource "akamai_gtm_domain" `)
var gtmDomainConfigP2 = fmt.Sprintf(`    contract = "${var.contractid}"
    group = "${var.groupid}"
    comment =  "Domain import"
`)

// misc
var gtmRConfigP2 = fmt.Sprintf(`    domain = "${akamai_gtm_domain.`)

var dependsClauseP1 = fmt.Sprintf(`    depends_on = [
        "akamai_gtm_domain.`)

// process domain 
func processDomain(domain *gtm.Domain, resourceDomainName string) (string) {

	domainBody := ""
        domainString := gtmHeaderConfig
    	domElems := reflect.ValueOf(domain).Elem()
    	for i := 0; i < domElems.NumField(); i++ {
        	varName := domElems.Type().Field(i).Name
        	varType := domElems.Type().Field(i).Type
        	varValue := domElems.Field(i).Interface()
		key := convertKey(varName)
		if key == "" {
			continue
		}
		keyVal := fmt.Sprint(varValue)
		if varName == "EmailNotificationList" {
			keyVal = processStringList(domain.EmailNotificationList)
		}
		domainBody += tab4 + key + " = "
		if varType.Kind() == reflect.String {
			domainBody += "\"" + keyVal + "\"\n"
		} else {
			domainBody += keyVal + "\n"
		} 
	}
	domainString += "\"" + resourceDomainName + "\" {\n"
	domainString += gtmDomainConfigP2
	domainString += domainBody
	domainString += "}\n\n"
	
	return domainString

}

func processStringList(sl []string) string {

	keyVal := "[" + strings.Join(sl, " ") + "]"
	if len(sl) > 0 {
		keyVal = strings.ReplaceAll(keyVal, " ", "\", \"")
        	keyVal = strings.Replace(keyVal, "[", "[\"", 1)
        	keyVal = strings.Replace(keyVal, "]", "\"]", 1)
	}
	return keyVal

}

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

func convertKey(inKey string) string {

	if _, ok := ignoredKeys[inKey]; ok {
		return ""
	}
	if _, ok := intNullableKeys[inKey]; ok {
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

