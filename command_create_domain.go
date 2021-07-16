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
	"encoding/json"
	"fmt"
	configgtm "github.com/akamai/AkamaiOPEN-edgegrid-golang/configgtm-v1_4"
	akamai "github.com/akamai/cli-common-golang"
	"github.com/fatih/color"
	"github.com/urfave/cli"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

var defaultDCs = []int{5400, 5401, 5402}

// Terraform resource names
var domainResource = "akamai_gtm_domain"
var datacenterResource = "akamai_gtm_datacenter"
var defaultDatacenterDataSource = "akamai_gtm_default_datacenter"
var propertyResource = "akamai_gtm_property"
var resourceResource = "akamai_gtm_resource"
var asResource = "akamai_gtm_asmap"
var geoResource = "akamai_gtm_geomap"
var cidrResource = "akamai_gtm_cidrmap"

// Import List Struct
type importListStruct struct {
	Domain      string
	Datacenters map[int]string
	Properties  map[string][]int
	Resources   map[string][]int
	Cidrmaps    map[string][]int
	Geomaps     map[string][]int
	Asmaps      map[string][]int
}

var tfWorkPath = "./"
var createImportList = false
var createConfig = false

var domainName string
var fullImportList *importListStruct

// text for gtmvars.tf construction
var gtmvarsContent = fmt.Sprint(`variable "gtmsection" {
  default = "default"
}
// Notice: contractid unknown at time of import. Please update.
variable "contractid" {
  default = ""
}
// Notice: groupid unknown at time of import. Please update.
variable "groupid" {
  default = ""
}
`)

var nullFieldMap = &configgtm.NullFieldMapStruct{}

// retrieve Null Values for Domain
func getDomainNullValues() configgtm.NullPerObjectAttributeStruct {

	return nullFieldMap.Domain

}

// retrieve Null Values for Object Type
func getNullValuesList(objType string) map[string]configgtm.NullPerObjectAttributeStruct {

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
	return map[string]configgtm.NullPerObjectAttributeStruct{}
}

// command function create-domain
func cmdCreateDomain(c *cli.Context) error {

	config, err := akamai.GetEdgegridConfig(c)
	if err != nil {
		return err
	}

	configgtm.Init(config)

	if c.NArg() < 1 {
		cli.ShowCommandHelp(c, c.Command.Name)
		return cli.NewExitError(color.RedString("domain is required"), 1)
	}

	domainName = c.Args().Get(0)
	if c.IsSet("tfworkpath") {
		tfWorkPath = c.String("tfworkpath")
	}
	tfWorkPath = filepath.FromSlash(tfWorkPath)
	if c.IsSet("resources") {
		createImportList = true
	}
	if c.IsSet("createconfig") {
		createConfig = true
	}

	fmt.Println("Configuring Domain")
	akamai.StartSpinner("Fetching domain entity ", "")
	domain, err := configgtm.GetDomain(domainName)
	if err != nil {
		akamai.StopSpinnerFail()
		fmt.Println("Error: " + err.Error())
		return cli.NewExitError(color.RedString("Domain retrieval failed"), 1)
	}
	akamai.StopSpinnerOk()
	// use domain name sans suffix for domain resource name
	resourceDomainName := normalizeResourceName(strings.TrimSuffix(domainName, ".akadns.net"))
	if createImportList {
		akamai.StartSpinner("Inventorying domain objects ", "")
		// Inventory datacenters
		datacenters := make(map[int]string)
		for _, dc := range domain.Datacenters {
			// include Default DCs. Special handling elsewhere.
			datacenters[dc.DatacenterId] = dc.Nickname
		}
		// inventory properties and targets
		propTargets := make(map[string][]int)
		for _, p := range domain.Properties {
			targets := make([]int, 0)
			for _, tt := range p.TrafficTargets {
				targets = append(targets, tt.DatacenterId)
			}
			propTargets[p.Name] = targets
		}
		// inventory Resources
		resources := make(map[string][]int)
		for _, r := range domain.Resources {
			targets := make([]int, 0)
			for _, ri := range r.ResourceInstances {
				targets = append(targets, ri.DatacenterId)
			}
			resources[r.Name] = targets
		}
		// inventory CidrMaps
		cidrmaps := make(map[string][]int)
		for _, c := range domain.CidrMaps {
			targets := make([]int, 0)
			for _, a := range c.Assignments {
				targets = append(targets, a.DatacenterId)
			}
			cidrmaps[c.Name] = targets
		}
		// inventory GeoMaps
		geomaps := make(map[string][]int)
		for _, g := range domain.GeographicMaps {
			targets := make([]int, 0)
			for _, a := range g.Assignments {
				targets = append(targets, a.DatacenterId)
			}
			geomaps[g.Name] = targets
		}
		// inventory ASMaps
		asmaps := make(map[string][]int)
		for _, as := range domain.AsMaps {
			targets := make([]int, 0)
			for _, a := range as.Assignments {
				targets = append(targets, a.DatacenterId)
			}
			asmaps[as.Name] = targets
		}
		akamai.StopSpinnerOk()
		akamai.StartSpinner("Creating Resources list file ", "")
		// pathname and exists?
		if stat, err := os.Stat(tfWorkPath); err == nil && stat.IsDir() {
			importListFilename := createImportListFilename(resourceDomainName)
			if _, err := os.Stat(importListFilename); err == nil {
				akamai.StopSpinnerFail()
				return cli.NewExitError(color.RedString("Resource list file exists. Remove to continue."), 1)
			}
			fullImportList = &importListStruct{}
			fullImportList.Domain = domainName
			fullImportList.Properties = propTargets
			fullImportList.Datacenters = datacenters
			fullImportList.Resources = resources
			fullImportList.Cidrmaps = cidrmaps
			fullImportList.Geomaps = geomaps
			fullImportList.Asmaps = asmaps
			json, err := json.MarshalIndent(fullImportList, "", "  ")
			if err != nil {
				akamai.StopSpinnerFail()
				return cli.NewExitError(color.RedString("Unable to generate json formatted resource list"), 1)
			}
			f, err := os.Create(importListFilename)
			if err != nil {
				akamai.StopSpinnerFail()
				return cli.NewExitError(color.RedString("Unable to create resources file"), 1)
			}
			defer f.Close()
			_, err = f.WriteString(string(json))
			if err != nil {
				akamai.StopSpinnerFail()
				return cli.NewExitError(color.RedString("Unable to write resources file"), 1)
			}
			f.Sync()
		} else {
			// Path doesnt exist. Bail
			akamai.StopSpinnerFail()
			return cli.NewExitError(color.RedString("Destination work path is not accessible."), 1)
		}
		akamai.StopSpinnerOk()
	}

	if createConfig {
		// Read in resources list
		importList, err := retrieveImportList(resourceDomainName)
		if err != nil {
			akamai.StopSpinnerFail()
			return cli.NewExitError(color.RedString("Failed to read json resources file"), 1)
		}
		akamai.StartSpinner("Creating domain configuration file ", "")
		// see if configuration file already exists and exclude any resources already represented.
		domainTFfileHandle, tfConfig, configImportList, err := reconcileResourceTargets(importList, resourceDomainName)
		if err != nil {
			akamai.StopSpinnerFail()
			return cli.NewExitError(color.RedString("Failed to open/create config file."), 1)
		}
		defer domainTFfileHandle.Close()
		//initialize Null Fields Struct
		nullFieldMap, err = domain.NullFieldMap()
		if err != nil {
			akamai.StopSpinnerFail()
			return cli.NewExitError(color.RedString("Failed to initialize Domain null fields map"), 1)
		}
		// build tf file
		if len(tfConfig) == 0 {
			// if tf pre existed, domain has to exist by definition
			tfConfig = processDomain(domain, resourceDomainName)
		}
		tfConfig += processDatacenters(domain.Datacenters, configImportList.Datacenters, resourceDomainName)
		tfConfig += processProperties(domain.Properties, configImportList.Properties, importList.Datacenters, resourceDomainName)
		tfConfig += processResources(domain.Resources, configImportList.Resources, importList.Datacenters, resourceDomainName)
		tfConfig += processCidrmaps(domain.CidrMaps, configImportList.Cidrmaps, importList.Datacenters, resourceDomainName)
		tfConfig += processGeomaps(domain.GeographicMaps, configImportList.Geomaps, importList.Datacenters, resourceDomainName)
		tfConfig += processAsmaps(domain.AsMaps, configImportList.Asmaps, importList.Datacenters, resourceDomainName)
		tfConfig += "\n"
		_, err = domainTFfileHandle.Write([]byte(tfConfig))
		if err != nil {
			akamai.StopSpinnerFail()
			return cli.NewExitError(color.RedString("Failed to save domain configuration file."), 1)
		}
		domainTFfileHandle.Sync()

		// Need create gtmvars.tf dependency
		gtmvarsFilename := filepath.Join(tfWorkPath, "gtmvars.tf")
		gtmvarsHandle, err := os.Create(gtmvarsFilename)
		if err != nil {
			akamai.StopSpinnerFail()
			return cli.NewExitError(color.RedString("Unable to create gtmvars config file"), 1)
		}
		defer gtmvarsHandle.Close()
		_, err = gtmvarsHandle.WriteString(gtmvarsContent)
		if err != nil {
			akamai.StopSpinnerFail()
			return cli.NewExitError(color.RedString("Unable to write gtmvars config file"), 1)
		}
		gtmvarsHandle.Sync()
		akamai.StopSpinnerOk()
		akamai.StartSpinner("Creating domain import script file ", "")
		importScriptFilename := filepath.Join(tfWorkPath, resourceDomainName+"_resource_import.script")
		if _, err := os.Stat(importScriptFilename); err == nil {
			// File exists. Bail
			akamai.StopSpinnerOk()
		}
		scriptContent, err := buildImportScript(configImportList, resourceDomainName)
		if err != nil {
			akamai.StopSpinnerFail()
			return cli.NewExitError(color.RedString("Import script content generation failed"), 1)
		}
		f, err := os.Create(importScriptFilename)
		if err != nil {
			akamai.StopSpinnerFail()
			return cli.NewExitError(color.RedString("Unable to create import script file"), 1)
		}
		defer f.Close()
		_, err = f.WriteString(scriptContent)
		if err != nil {
			akamai.StopSpinnerFail()
			return cli.NewExitError(color.RedString("Unable to write import script file"), 1)
		}
		f.Sync()
		akamai.StopSpinnerOk()
	}

	fmt.Println("Domain configuration completed")

	return nil
}

func retrieveImportList(rscName string) (*importListStruct, error) {

	// check if createImportList set. If so, already have ....
	if createImportList {
		return fullImportList, nil
	}
	importListFilename := createImportListFilename(rscName)
	if _, err := os.Stat(importListFilename); err != nil {
		return nil, err
	}
	importData, err := ioutil.ReadFile(importListFilename)
	if err != nil {
		return nil, err
	}
	importList := &importListStruct{}
	err = json.Unmarshal(importData, importList)
	if err != nil {
		return nil, err
	}

	return importList, nil

}

// Utility method to create full import list file path
func createImportListFilename(resourceName string) string {

	return filepath.Join(tfWorkPath, resourceName+"_resources.json")

}

// Utility method to create full tf file path
func createTFFilename(resourceName string) string {

	return filepath.Join(tfWorkPath, resourceName+".tf")

}

func buildImportScript(importList *importListStruct, resourceDomainName string) (string, error) {

	// build import script
	var import_prefix = "terraform import "
	var import_file = ""
	// Init TF
	import_file += "terraform init\n"
	// domain
	if !checkForResource(domainResource, resourceDomainName) {
		// Assuming a domain name cannot contain spaces ....
		import_file += import_prefix + domainResource + "." + resourceDomainName + " " + importList.Domain + "\n"
	}
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
		if !checkForResource(datacenterResource, normalName) {
			import_file += import_prefix + datacenterResource + "." + normalName + " " + importList.Domain + ":" + strconv.Itoa(id) + "\n"
		}
	}
	// properties
	for name, _ := range importList.Properties {
		normalName := normalizeResourceName(name)
		if !checkForResource(propertyResource, normalName) {
			import_file += import_prefix + propertyResource + "." + normalName + " "
			if strings.Contains(name, " ") {
				import_file += `"` + importList.Domain + ":" + name + `"` + "\n"
			} else {
				import_file += importList.Domain + ":" + name + "\n"
			}
		}
	}
	// resources
	for name, _ := range importList.Resources {
		normalName := normalizeResourceName(name)
		if !checkForResource(resourceResource, normalName) {
			import_file += import_prefix + resourceResource + "." + normalName + " "
			if strings.Contains(name, " ") {
				import_file += `"` + importList.Domain + ":" + name + `"` + "\n"
			} else {
				import_file += importList.Domain + ":" + name + "\n"
			}
		}
	}
	// cidrmaps
	for name, _ := range importList.Cidrmaps {
		normalName := normalizeResourceName(name)
		if !checkForResource(cidrResource, normalName) {
			import_file += import_prefix + cidrResource + "." + normalName + " "
			if strings.Contains(name, " ") {
				import_file += `"` + importList.Domain + ":" + name + `"` + "\n"
			} else {
				import_file += importList.Domain + ":" + name + "\n"
			}
		}
	}
	// geomaps
	for name, _ := range importList.Geomaps {
		normalName := normalizeResourceName(name)
		if !checkForResource(geoResource, normalName) {
			import_file += import_prefix + geoResource + "." + normalName + " "
			if strings.Contains(name, " ") {
				import_file += `"` + importList.Domain + ":" + name + `"` + "\n"
			} else {
				import_file += importList.Domain + ":" + name + "\n"
			}
		}
	}
	// asmaps
	for name, _ := range importList.Asmaps {
		normalName := normalizeResourceName(name)
		if !checkForResource(asResource, normalName) {
			import_file += import_prefix + asResource + "." + normalName + " "
			if strings.Contains(name, " ") {
				import_file += `"` + importList.Domain + ":" + name + `"` + "\n"
			} else {
				import_file += importList.Domain + ":" + name + "\n"
			}
		}
	}

	return import_file, nil

}

// remove any resources already present in existing domain tf configuration
func reconcileResourceTargets(importList *importListStruct, domainName string) (*os.File, string, *importListStruct, error) {

	var tfScratchLen int64
	tfFilename := createTFFilename(domainName)
	if tfInfo, err := os.Stat(tfFilename); err != nil && os.IsExist(err) {
		tfScratchLen = tfInfo.Size()
	}
	tfScratch := make([]byte, tfScratchLen)
	var tfHandle *os.File
	tfHandle, err := os.OpenFile(tfFilename, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil && err != io.EOF {
		fmt.Println(err.Error())
		return nil, "", importList, err
	}
	if tfScratchLen == 0 {
		return tfHandle, "", importList, nil
	}
	charsRead, err := tfHandle.Read(tfScratch)
	if err != nil && err != io.EOF {
		fmt.Println(err.Error())
		return nil, "", importList, err
	}
	_, err = tfHandle.Seek(0, 0)
	if err != nil {
		fmt.Println(err.Error())
		return nil, "", importList, err
	}
	if charsRead == 0 {
		return tfHandle, "", importList, err
	}
	tfConfig := fmt.Sprintf("%s", tfScratch[0:charsRead-1])
	// need walk thru each resource type
	for id, nickname := range importList.Datacenters {
		normalName := normalizeResourceName(nickname)
		if strings.Contains(tfConfig, "\""+normalName+"\"") {
			fmt.Println("Datacenter " + nickname + " found in existing tf file")
			delete(importList.Datacenters, id)
		}
	}
	for name, _ := range importList.Properties {
		normalName := normalizeResourceName(name)
		if strings.Contains(tfConfig, "\""+normalName+"\"") {
			fmt.Println("Property " + name + " found in existing tf file")
			delete(importList.Properties, name)
		}
	}
	for name, _ := range importList.Resources {
		normalName := normalizeResourceName(name)
		if strings.Contains(tfConfig, "\""+normalName+"\"") {
			fmt.Println("Resource " + name + " found in existing tf file")
			delete(importList.Resources, name)
		}
	}
	for name, _ := range importList.Cidrmaps {
		normalName := normalizeResourceName(name)
		if strings.Contains(tfConfig, "\""+normalName+"\"") {
			fmt.Println("Cidrmap " + name + " found in existing tf file")
			delete(importList.Cidrmaps, name)
		}
	}
	for name, _ := range importList.Geomaps {
		normalName := normalizeResourceName(name)
		if strings.Contains(tfConfig, "\""+normalName+"\"") {
			fmt.Println("Geomap " + name + " found in existing tf file")
			delete(importList.Geomaps, name)
		}
	}
	for name, _ := range importList.Asmaps {
		normalName := normalizeResourceName(name)
		if strings.Contains(tfConfig, "\""+normalName+"\"") {
			fmt.Println("Asmap " + name + " found in existing tf file")
			delete(importList.Asmaps, name)
		}
	}
	return tfHandle, tfConfig, importList, err

}
