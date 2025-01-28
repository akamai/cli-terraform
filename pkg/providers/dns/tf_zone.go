package dns

import (
	"context"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/dns"
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
		OutboundZoneTransfer:  zone.OutboundZoneTransfer,
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
