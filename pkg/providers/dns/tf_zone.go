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

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/dns"
)

// process zone
func processZone(ctx context.Context, zone *dns.GetZoneResponse, resourceZoneName string, modSegment bool, fileUtils fileUtils, tfWorkPath string) (string, error) {
	data := ZoneData{
		BlockName:             resourceZoneName,
		Zone:                  zone.Zone,
		Type:                  zone.Type,
		Masters:               zone.Masters,
		Comment:               zone.Comment,
		SignAndServe:          zone.SignAndServe,
		SignAndServeAlgorithm: zone.SignAndServeAlgorithm,
		TSIGKey:               zone.TSIGKey,
		Target:                zone.Target,
		EndCustomerID:         zone.EndCustomerID,
		TFWorkPath:            tfWorkPath,
	}
	var zoneTF string
	if modSegment {
		err := fileUtils.createModuleTF(ctx, resourceZoneName, useTemplate(&data, "config.tmpl", true), tfWorkPath)
		if err != nil {
			return "", err
		}
		zoneTF = useTemplate(&data, "zone.tmpl", true)
	} else {
		zoneTF = useTemplate(&data, "full_zone.tmpl", true)
	}

	return zoneTF, nil
}
