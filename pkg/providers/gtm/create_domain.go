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
	"context"
	"embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"

	gtm "github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/configgtm"
	"github.com/akamai/cli-terraform/pkg/edgegrid"
	"github.com/akamai/cli-terraform/pkg/templates"
	"github.com/akamai/cli-terraform/pkg/tools"
	"github.com/akamai/cli/pkg/terminal"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

type (
	// TFDomainData represents the data used in domain templates
	TFDomainData struct {
		NormalizedName              string
		Name                        string
		Type                        string
		Comment                     string
		EmailNotificationList       []string
		DefaultTimeoutPenalty       int
		LoadImbalancePercentage     float64
		DefaultSSLClientPrivateKey  string
		DefaultErrorPenalty         int
		CnameCoalescingEnabled      bool
		LoadFeedback                bool
		DefaultSSLClientCertificate string
		EndUserMappingEnabled       bool
		Section                     string
		Datacenters                 []TFDatacenterData
		DefaultDatacenters          []TFDatacenterData
		Resources                   []*gtm.Resource
		CidrMaps                    []*gtm.CidrMap
		GeoMaps                     []*gtm.GeoMap
		AsMaps                      []*gtm.AsMap
		Properties                  []*gtm.Property
	}

	// TFDatacenterData represents the data used for processing a datacenter
	TFDatacenterData struct {
		ID                            int
		Nickname                      string
		City                          string
		CloneOf                       int
		CloudServerHostHeaderOverride bool
		CloudServerTargeting          bool
		Continent                     string
		Country                       string
		Latitude                      float64
		Longitude                     float64
		StateOrProvince               string
		DefaultLoadObject             *gtm.LoadObject
	}
)

//go:embed templates/*
var templateFiles embed.FS

var defaultDCs = map[int]struct{}{5400: {}, 5401: {}, 5402: {}}

var (
	subWithUnderscoreRegexp               = regexp.MustCompile(`[^\w-_]`)
	mustStartWithLetterOrUnderscoreRegexp = regexp.MustCompile("^[^a-zA-Z_]")
	// ErrFetchingDomain is returned when fetching domain fails
	ErrFetchingDomain = errors.New("unable to fetch domain with given name")
)

// CmdCreateDomain is an entrypoint to create-domain command
func CmdCreateDomain(c *cli.Context) error {
	ctx := c.Context
	if c.NArg() != 1 {
		if err := cli.ShowCommandHelp(c, c.Command.Name); err != nil {
			return cli.Exit(color.RedString("Error displaying help command"), 1)
		}
		return cli.Exit(color.RedString("Domain is required"), 1)
	}

	sess := edgegrid.GetSession(ctx)
	client := gtm.Client(sess)
	if c.IsSet("tfworkpath") {
		tools.TFWorkPath = c.String("tfworkpath")
	}
	tools.TFWorkPath = filepath.FromSlash(tools.TFWorkPath)
	if stat, err := os.Stat(tools.TFWorkPath); err != nil || !stat.IsDir() {
		return cli.Exit(color.RedString("Destination work path is not accessible"), 1)
	}

	datacentersPath := filepath.Join(tools.TFWorkPath, "datacenters.tf")
	domainPath := filepath.Join(tools.TFWorkPath, "domain.tf")
	importPath := filepath.Join(tools.TFWorkPath, "import.sh")
	mapsPath := filepath.Join(tools.TFWorkPath, "maps.tf")
	propertiesPath := filepath.Join(tools.TFWorkPath, "properties.tf")
	resourcesPath := filepath.Join(tools.TFWorkPath, "resources.tf")
	variablesPath := filepath.Join(tools.TFWorkPath, "variables.tf")

	templateToFile := map[string]string{
		"datacenters.tmpl": datacentersPath,
		"domain.tmpl":      domainPath,
		"imports.tmpl":     importPath,
		"maps.tmpl":        mapsPath,
		"properties.tmpl":  propertiesPath,
		"resources.tmpl":   resourcesPath,
		"variables.tmpl":   variablesPath,
	}

	err := tools.CheckFiles(datacentersPath, domainPath, importPath, mapsPath, propertiesPath, resourcesPath, variablesPath)
	if err != nil {
		return cli.Exit(color.RedString(err.Error()), 1)
	}

	processor := templates.FSTemplateProcessor{
		TemplatesFS:     templateFiles,
		TemplateTargets: templateToFile,
		AdditionalFuncs: template.FuncMap{
			"normalize": normalizeResourceName,
			"toUpper":   strings.ToUpper,
		},
	}

	domainName := c.Args().First()
	section := edgegrid.GetEdgercSection(c)
	if err := createDomain(ctx, client, domainName, section, processor); err != nil {
		return cli.Exit(color.RedString(fmt.Sprintf("Error exporting domain HCL: %s", err)), 1)
	}
	return nil
}

func createDomain(ctx context.Context, client gtm.GTM, domainName, section string, templateProcessor templates.TemplateProcessor) error {
	term := terminal.Get(ctx)

	term.Writeln("Configuring Domain")
	term.Spinner().Start(fmt.Sprintf("Fetching domain %s", domainName))
	domain, err := client.GetDomain(ctx, domainName)
	if err != nil {
		term.Spinner().Fail()
		return fmt.Errorf("%w: %s", ErrFetchingDomain, err)
	}

	tfDomainData := TFDomainData{
		Section:                     section,
		NormalizedName:              normalizeResourceName(strings.TrimSuffix(domain.Name, ".akadns.net")),
		Name:                        domain.Name,
		Type:                        domain.Type,
		Comment:                     domain.ModificationComments,
		EmailNotificationList:       domain.EmailNotificationList,
		DefaultTimeoutPenalty:       domain.DefaultTimeoutPenalty,
		LoadImbalancePercentage:     domain.LoadImbalancePercentage,
		DefaultSSLClientPrivateKey:  domain.DefaultSslClientPrivateKey,
		DefaultErrorPenalty:         domain.DefaultErrorPenalty,
		CnameCoalescingEnabled:      domain.CnameCoalescingEnabled,
		LoadFeedback:                domain.LoadFeedback,
		DefaultSSLClientCertificate: domain.DefaultSslClientCertificate,
		EndUserMappingEnabled:       domain.EndUserMappingEnabled,
		Resources:                   domain.Resources,
		CidrMaps:                    domain.CidrMaps,
		GeoMaps:                     domain.GeographicMaps,
		AsMaps:                      domain.AsMaps,
		Properties:                  domain.Properties,
	}

	tfDomainData.getDatacenters(domain)
	term.Spinner().OK()

	term.Spinner().Start("Saving TF configurations")
	if err := templateProcessor.ProcessTemplates(tfDomainData); err != nil {
		term.Spinner().Fail()
		return err
	}
	term.Spinner().OK()

	term.Writeln(fmt.Sprintf("Terraform configuration for policy '%s' was saved successfully\n", domain.Name))

	return nil
}

func (d *TFDomainData) getDatacenters(domain *gtm.Domain) {
	d.Datacenters = make([]TFDatacenterData, 0)
	d.DefaultDatacenters = make([]TFDatacenterData, 0)
	for _, dc := range domain.Datacenters {
		if _, ok := defaultDCs[dc.DatacenterId]; ok {
			d.DefaultDatacenters = append(d.DefaultDatacenters, TFDatacenterData{Nickname: dc.Nickname, ID: dc.DatacenterId})
		} else {
			d.Datacenters = append(d.Datacenters, TFDatacenterData{
				ID:                            dc.DatacenterId,
				Nickname:                      dc.Nickname,
				City:                          dc.City,
				CloneOf:                       dc.CloneOf,
				CloudServerHostHeaderOverride: dc.CloudServerHostHeaderOverride,
				CloudServerTargeting:          dc.CloudServerTargeting,
				Continent:                     dc.Continent,
				Country:                       dc.Country,
				Latitude:                      dc.Latitude,
				Longitude:                     dc.Longitude,
				StateOrProvince:               dc.StateOrProvince,
				DefaultLoadObject:             dc.DefaultLoadObject,
			})
		}
	}
}

// normalizeResourceName is a utility function to normalize resource names.
// A name must start with a letter or underscore and may contain only letters, digits, underscores, and dashes.
func normalizeResourceName(key string) string {
	key = subWithUnderscoreRegexp.ReplaceAllString(key, "_")

	if mustStartWithLetterOrUnderscoreRegexp.MatchString(key) {
		key = fmt.Sprintf("_%s", key)
	}
	return key
}

// FindDatacenterResourceName finds and returns datacenter resource name with given id
func (d TFDomainData) FindDatacenterResourceName(id int) (string, error) {
	for _, dc := range d.Datacenters {
		if dc.ID == id {
			return normalizeResourceName(dc.Nickname), nil
		}
	}
	return "", fmt.Errorf("cannot find datacenter resource with ID: %d", id)
}