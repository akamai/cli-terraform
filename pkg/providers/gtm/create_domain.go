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
	"io"
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
		Datacenters                 map[int]string
		DefaultDatacenters          map[int]string
		Properties                  map[string][]int
		Resources                   map[string][]int
		Cidrmaps                    map[string][]int
		Geomaps                     map[string][]int
		Asmaps                      map[string][]int
	}
)

//go:embed templates/*
var templateFiles embed.FS

var defaultDCs = map[int]struct{}{5400: {}, 5401: {}, 5402: {}}

// Terraform resource names
const (
	domainResource              = "akamai_gtm_domain"
	datacenterResource          = "akamai_gtm_datacenter"
	defaultDatacenterDataSource = "akamai_gtm_default_datacenter"
	propertyResource            = "akamai_gtm_property"
	resourceResource            = "akamai_gtm_resource"
	asResource                  = "akamai_gtm_asmap"
	geoResource                 = "akamai_gtm_geomap"
	cidrResource                = "akamai_gtm_cidrmap"
)

// TODO: remove and declare those variables in CmdCreateDomain once there is no appending to those files (DXE-698)
var (
	domainPath string
	importPath string
)

var nullFieldMap = &gtm.NullFieldMapStruct{}

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

	variablesPath := filepath.Join(tools.TFWorkPath, "variables.tf")
	domainPath = filepath.Join(tools.TFWorkPath, "domain.tf")
	importPath = filepath.Join(tools.TFWorkPath, "import.sh")

	templateToFile := map[string]string{
		"domain.tmpl":    domainPath,
		"variables.tmpl": variablesPath,
		"imports.tmpl":   importPath,
	}

	err := tools.CheckFiles(domainPath, variablesPath, importPath)
	if err != nil {
		return cli.Exit(color.RedString(err.Error()), 1)
	}

	processor := templates.FSTemplateProcessor{
		TemplatesFS:     templateFiles,
		TemplateTargets: templateToFile,
		AdditionalFuncs: template.FuncMap{
			"normalize": normalizeResourceName,
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

	fmt.Println("Configuring Domain")
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
	}

	createImportList(domain, &tfDomainData)
	term.Spinner().OK()

	term.Spinner().Start("Saving TF configurations ")
	if err := templateProcessor.ProcessTemplates(tfDomainData); err != nil {
		term.Spinner().Fail()
		return err
	}
	term.Spinner().OK()

	term.Spinner().Start("Creating domain configuration file ")
	if err := createConfig(ctx, client, domain, &tfDomainData); err != nil {
		term.Spinner().Fail()
		return err
	}
	term.Spinner().OK()

	fmt.Printf("Terraform configuration for policy '%s' was saved successfully\n", domain.Name)

	return nil
}

// retrieve Null Values for Object Type
func getNullValuesList(objType string) map[string]gtm.NullPerObjectAttributeStruct {

	switch objType {
	case "Properties":
		return nullFieldMap.Properties
	case "Datacenters":
		return nullFieldMap.Datacenters
	case "Resources":
		return nullFieldMap.Resources
	case "CidrMaps":
		return nullFieldMap.CidrMaps
	case "GeoMaps":
		return nullFieldMap.GeoMaps
	case "AsMaps":
		return nullFieldMap.AsMaps
	}
	// unknown
	return map[string]gtm.NullPerObjectAttributeStruct{}
}

func createImportList(domain *gtm.Domain, tfData *TFDomainData) {
	// Inventory datacenters
	tfData.Datacenters = make(map[int]string, len(domain.Datacenters))
	tfData.DefaultDatacenters = make(map[int]string, len(domain.Datacenters))
	for _, dc := range domain.Datacenters {
		if _, ok := defaultDCs[dc.DatacenterId]; ok {
			tfData.DefaultDatacenters[dc.DatacenterId] = dc.Nickname
		} else {
			tfData.Datacenters[dc.DatacenterId] = dc.Nickname
		}
	}
	// inventory properties and targets
	tfData.Properties = make(map[string][]int, len(domain.Properties))
	for _, p := range domain.Properties {
		targets := make([]int, 0, len(p.TrafficTargets))
		for _, tt := range p.TrafficTargets {
			targets = append(targets, tt.DatacenterId)
		}
		tfData.Properties[p.Name] = targets
	}
	// inventory Resources
	tfData.Resources = make(map[string][]int, len(domain.Resources))
	for _, r := range domain.Resources {
		targets := make([]int, 0, len(r.ResourceInstances))
		for _, ri := range r.ResourceInstances {
			targets = append(targets, ri.DatacenterId)
		}
		tfData.Resources[r.Name] = targets
	}
	// inventory CidrMaps
	tfData.Cidrmaps = make(map[string][]int, len(domain.CidrMaps))
	for _, c := range domain.CidrMaps {
		targets := make([]int, 0, len(c.Assignments))
		for _, a := range c.Assignments {
			targets = append(targets, a.DatacenterId)
		}
		tfData.Cidrmaps[c.Name] = targets
	}
	// inventory GeoMaps
	tfData.Geomaps = make(map[string][]int, len(domain.GeographicMaps))
	for _, g := range domain.GeographicMaps {
		targets := make([]int, 0, len(g.Assignments))
		for _, a := range g.Assignments {
			targets = append(targets, a.DatacenterId)
		}
		tfData.Geomaps[g.Name] = targets
	}
	// inventory ASMaps
	tfData.Asmaps = make(map[string][]int, len(domain.AsMaps))
	for _, as := range domain.AsMaps {
		targets := make([]int, 0, len(as.Assignments))
		for _, a := range as.Assignments {
			targets = append(targets, a.DatacenterId)
		}
		tfData.Asmaps[as.Name] = targets
	}
}

func createConfig(ctx context.Context, client gtm.GTM, domain *gtm.Domain, tfData *TFDomainData) error {
	domainTFfileHandle, err := os.OpenFile(domainPath, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil && err != io.EOF {
		return err
	}
	defer domainTFfileHandle.Close()

	//initialize Null Fields Struct
	nullFieldMap, err = client.NullFieldMap(ctx, domain)
	if err != nil {
		return fmt.Errorf("failed to initialize Domain null fields map")
	}
	// build tf file
	tfConfig := processDatacenters(domain.Datacenters, tfData.NormalizedName)
	tfConfig += processProperties(domain.Properties, tfData.Properties, tfData.Datacenters, tfData.NormalizedName)
	tfConfig += processResources(domain.Resources, tfData.Resources, tfData.Datacenters, tfData.NormalizedName)
	tfConfig += processCidrmaps(domain.CidrMaps, tfData.Cidrmaps, tfData.Datacenters, tfData.NormalizedName)
	tfConfig += processGeomaps(domain.GeographicMaps, tfData.Geomaps, tfData.Datacenters, tfData.NormalizedName)
	tfConfig += processAsmaps(domain.AsMaps, tfData.Asmaps, tfData.Datacenters, tfData.NormalizedName)
	tfConfig += "\n"

	_, err = domainTFfileHandle.Write([]byte(tfConfig))
	if err != nil {
		return fmt.Errorf("failed to save domain configuration file")
	}
	domainTFfileHandle.Sync()
	return nil
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
