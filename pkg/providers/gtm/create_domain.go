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
	"strconv"
	"strings"

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
		Section string
	}

	importListStruct struct {
		Domain      string
		Datacenters map[int]string
		Properties  map[string][]int
		Resources   map[string][]int
		Cidrmaps    map[string][]int
		Geomaps     map[string][]int
		Asmaps      map[string][]int
	}
)

//go:embed templates/*
var templateFiles embed.FS

var defaultDCs = []int{5400, 5401, 5402}

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
	// templatizing import script will be part of DXE-692
	importPath = filepath.Join(tools.TFWorkPath, "import.sh")

	templateToFile := map[string]string{
		"domain.tmpl":    domainPath,
		"variables.tmpl": variablesPath,
	}

	err := tools.CheckFiles(domainPath, variablesPath, importPath)
	if err != nil {
		return cli.Exit(color.RedString(err.Error()), 1)
	}

	processor := templates.FSTemplateProcessor{
		TemplatesFS:     templateFiles,
		TemplateTargets: templateToFile,
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
		Section: section,
	}
	importList := createImportList(domain)
	term.Spinner().OK()

	term.Spinner().Start("Saving TF configurations ")
	if err := templateProcessor.ProcessTemplates(tfDomainData); err != nil {
		term.Spinner().Fail()
		return err
	}
	term.Spinner().OK()

	// use domain name sans suffix for domain resource name
	resourceDomainName := normalizeResourceName(strings.TrimSuffix(domain.Name, ".akadns.net"))

	term.Spinner().Start("Creating domain configuration file ")
	if err := createConfig(ctx, client, domain, importList, resourceDomainName); err != nil {
		term.Spinner().Fail()
		return err
	}
	term.Spinner().OK()

	term.Spinner().Start("Creating domain import script file ")
	if err := buildImportScript(importList, resourceDomainName); err != nil {
		term.Spinner().Fail()
		return err
	}
	term.Spinner().OK()
	fmt.Printf("Terraform configuration for policy '%s' was saved successfully\n", domain.Name)

	return nil
}

func buildImportScript(importList *importListStruct, resourceDomainName string) error {
	// build import script
	var importPrefix = "terraform import "
	var importFile = ""
	// Init TF
	importFile += "terraform init\n"
	// domain
	// Assuming a domain name cannot contain spaces ....
	importFile += importPrefix + domainResource + "." + resourceDomainName + " " + importList.Domain + "\n"

	// datacenters
	for id, nickname := range importList.Datacenters {
		normalName := normalizeResourceName(nickname)
		// default datacenters special case.
		ddcfound := false
		for _, ddc := range defaultDCs {
			if id == ddc {
				ddcfound = true
			}
		}
		if ddcfound {
			continue
		}
		importFile += importPrefix + datacenterResource + "." + normalName + " " + importList.Domain + ":" + strconv.Itoa(id) + "\n"
	}

	// properties
	for name := range importList.Properties {
		normalName := normalizeResourceName(name)
		importFile += importPrefix + propertyResource + "." + normalName + " "
		if strings.Contains(name, " ") {
			importFile += `"` + importList.Domain + ":" + name + `"` + "\n"
		} else {
			importFile += importList.Domain + ":" + name + "\n"
		}
	}

	// resources
	for name := range importList.Resources {
		normalName := normalizeResourceName(name)

		importFile += importPrefix + resourceResource + "." + normalName + " "
		if strings.Contains(name, " ") {
			importFile += `"` + importList.Domain + ":" + name + `"` + "\n"
		} else {
			importFile += importList.Domain + ":" + name + "\n"
		}
	}

	// cidrmaps
	for name := range importList.Cidrmaps {
		normalName := normalizeResourceName(name)

		importFile += importPrefix + cidrResource + "." + normalName + " "
		if strings.Contains(name, " ") {
			importFile += `"` + importList.Domain + ":" + name + `"` + "\n"
		} else {
			importFile += importList.Domain + ":" + name + "\n"
		}
	}

	// geomaps
	for name := range importList.Geomaps {
		normalName := normalizeResourceName(name)

		importFile += importPrefix + geoResource + "." + normalName + " "
		if strings.Contains(name, " ") {
			importFile += `"` + importList.Domain + ":" + name + `"` + "\n"
		} else {
			importFile += importList.Domain + ":" + name + "\n"
		}
	}

	// asmaps
	for name := range importList.Asmaps {
		normalName := normalizeResourceName(name)

		importFile += importPrefix + asResource + "." + normalName + " "
		if strings.Contains(name, " ") {
			importFile += `"` + importList.Domain + ":" + name + `"` + "\n"
		} else {
			importFile += importList.Domain + ":" + name + "\n"
		}
	}

	if err := os.WriteFile(importPath, []byte(importFile), 0666); err != nil {
		return fmt.Errorf("unable to write import script file")
	}
	return nil
}

// retrieve Null Values for Domain
func getDomainNullValues() gtm.NullPerObjectAttributeStruct {
	return nullFieldMap.Domain
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

func createImportList(domain *gtm.Domain) *importListStruct {
	importList := importListStruct{Domain: domain.Name}
	// Inventory datacenters
	importList.Datacenters = make(map[int]string, len(domain.Datacenters))
	for _, dc := range domain.Datacenters {
		// include Default DCs. Special handling elsewhere.
		importList.Datacenters[dc.DatacenterId] = dc.Nickname
	}
	// inventory properties and targets
	importList.Properties = make(map[string][]int, len(domain.Properties))
	for _, p := range domain.Properties {
		targets := make([]int, 0, len(p.TrafficTargets))
		for _, tt := range p.TrafficTargets {
			targets = append(targets, tt.DatacenterId)
		}
		importList.Properties[p.Name] = targets
	}
	// inventory Resources
	importList.Resources = make(map[string][]int, len(domain.Resources))
	for _, r := range domain.Resources {
		targets := make([]int, 0, len(r.ResourceInstances))
		for _, ri := range r.ResourceInstances {
			targets = append(targets, ri.DatacenterId)
		}
		importList.Resources[r.Name] = targets
	}
	// inventory CidrMaps
	importList.Cidrmaps = make(map[string][]int, len(domain.CidrMaps))
	for _, c := range domain.CidrMaps {
		targets := make([]int, 0, len(c.Assignments))
		for _, a := range c.Assignments {
			targets = append(targets, a.DatacenterId)
		}
		importList.Cidrmaps[c.Name] = targets
	}
	// inventory GeoMaps
	importList.Geomaps = make(map[string][]int, len(domain.GeographicMaps))
	for _, g := range domain.GeographicMaps {
		targets := make([]int, 0, len(g.Assignments))
		for _, a := range g.Assignments {
			targets = append(targets, a.DatacenterId)
		}
		importList.Geomaps[g.Name] = targets
	}
	// inventory ASMaps
	importList.Asmaps = make(map[string][]int, len(domain.AsMaps))
	for _, as := range domain.AsMaps {
		targets := make([]int, 0, len(as.Assignments))
		for _, a := range as.Assignments {
			targets = append(targets, a.DatacenterId)
		}
		importList.Asmaps[as.Name] = targets
	}

	return &importList
}

func createConfig(ctx context.Context, client gtm.GTM, domain *gtm.Domain, importList *importListStruct, resourceDomainName string) error {
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
	tfConfig := "\n"
	tfConfig += processDomain(domain, resourceDomainName)
	tfConfig += processDatacenters(domain.Datacenters, importList.Datacenters, resourceDomainName)
	tfConfig += processProperties(domain.Properties, importList.Properties, importList.Datacenters, resourceDomainName)
	tfConfig += processResources(domain.Resources, importList.Resources, importList.Datacenters, resourceDomainName)
	tfConfig += processCidrmaps(domain.CidrMaps, importList.Cidrmaps, importList.Datacenters, resourceDomainName)
	tfConfig += processGeomaps(domain.GeographicMaps, importList.Geomaps, importList.Datacenters, resourceDomainName)
	tfConfig += processAsmaps(domain.AsMaps, importList.Asmaps, importList.Datacenters, resourceDomainName)
	tfConfig += "\n"

	_, err = domainTFfileHandle.Write([]byte(tfConfig))
	if err != nil {
		return fmt.Errorf("failed to save domain configuration file")
	}
	domainTFfileHandle.Sync()
	return nil
}
