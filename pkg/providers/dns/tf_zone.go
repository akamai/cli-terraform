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

	dns "github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/configdns"
)

// process zone
func processZone(ctx context.Context, zone *dns.ZoneResponse, resourceZoneName string, modSegment bool, fileUtils fileUtils, tfworkPath string) (string, error) {
	data := ZoneData{
		BlockName:             resourceZoneName,
		Zone:                  zone.Zone,
		Type:                  zone.Type,
		Masters:               zone.Masters,
		Comment:               zone.Comment,
		SignAndServe:          zone.SignAndServe,
		SignAndServeAlgorithm: zone.SignAndServeAlgorithm,
		TsigKey:               zone.TsigKey,
		Target:                zone.Target,
		EndCustomerID:         zone.EndCustomerID,
		TfWorkPath:            tfworkPath,
	}
	var zoneTF string
	if modSegment {
		err := fileUtils.createModuleTF(ctx, resourceZoneName, useTemplate(&data, "config.tmpl", true), tfworkPath)
		if err != nil {
			return "", err
		}
		zoneTF = useTemplate(&data, "zone.tmpl", true)
	} else {
		zoneTF = useTemplate(&data, "full_zone.tmpl", true)
	}

	return zoneTF, nil

}
