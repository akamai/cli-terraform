// Copyright 2019. Akamai Technologies, Inc
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
	"path/filepath"
	"io/ioutil"
	"os"
	configgtm "github.com/akamai/AkamaiOPEN-edgegrid-golang/configgtm-v1_4"
	akamai "github.com/akamai/cli-common-golang"
	"github.com/fatih/color"
	//"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
	"strconv"
	"strings"
	//"time"
)

var defaultDC = 5400

// Terraform resource names
var domainResource = "akamai_gtm_domain"
var datacenterResource = "akamai_gtm_datacenter"
var propertyResource = "akamai_gtm_property"
var resourceResource = "akamai_gtm_resource"
var asResource = "akamai_gtm_asmap"
var geoResource = "akamai_gtm_geomap"
var cidrResource = "akamai_gtm_cidrmap"

// Import List Struct
type importListStruct struct {
	Domain 		string
	Datacenters 	map[int]string
	Properties	map[string][]int
	Resources	map[string][]int
	Cidrmaps	map[string][]int
	Geomaps		map[string][]int
	Asmaps		map[string][]int
}

var tfWorkPath = "./"
//var propImports []string
var createImportList = false 
var createConfig = false
//var importProps = false
var domainName string

var fullImportList *importListStruct = nil
 
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
        //propImports = c.StringSlice("property")
        if c.IsSet("importlist") {
		createImportList = true
	}
        if c.IsSet("createconfig") {
		createConfig = true
	}
	/*
        if c.IsSet("importconfig") {
		importProps = true
	}
	*/
        /*	
	if c.IsSet("verbose") {
		verboseStatus = true
	}
	if c.IsSet("complete") {
		pComplete = true
	}
        */

	domain, err := configgtm.GetDomain(domainName)
	if err != nil {
		akamai.StopSpinnerFail()
		fmt.Println(err.Error())
		return cli.NewExitError(color.RedString("Domain retrieval failed"), 1)
	}
	// Inventory datacenters
	datacenters := make(map[int]string)
	for _, dc := range domain.Datacenters {
		// special case. ignore 5400
		if dc.DatacenterId == defaultDC {
			continue
		}
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

	if createImportList {
		// pathname and exists?
		if stat, err := os.Stat(tfWorkPath); err == nil && stat.IsDir() {
			importListFilename := createImportListFilename()
			if _, err := os.Stat(importListFilename); err == nil {
				// File exists. Bail
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
                                return cli.NewExitError(color.RedString("Unable to create import list"), 1)
                        }
			fmt.Println(string(json))
			f, err := os.Create(importListFilename)
			if err != nil {
                                akamai.StopSpinnerFail()
                                return cli.NewExitError(color.RedString("Unable to create import list file"), 1)
                        }
			defer f.Close()
			_, err = f.WriteString(string(json))
                        if err != nil {
                                akamai.StopSpinnerFail()
                                return cli.NewExitError(color.RedString("Unable to write import list file"), 1)
                        }
			f.Sync()
		} else {
			// Path doesnt exist. Bail 
			
		}
	}

	if createConfig {
		// Read in import list
		importList, err := retrieveImportList()
		if err != nil {
			akamai.StopSpinnerFail()
                        return cli.NewExitError(color.RedString("Failed to read json domain definition"), 1)
		}
		// build tf file
		resourceDomainName := strings.TrimSuffix(domainName, ".akadns.net")
		tfConfig := processDomain(domain, resourceDomainName)
		tfConfig += processDatacenters(domain.Datacenters, importList.Datacenters, resourceDomainName)
		tfConfig += processProperties(domain.Properties, importList.Properties, importList.Datacenters, resourceDomainName)
                tfConfig += processResources(domain.Resources, importList.Resources, importList.Datacenters, resourceDomainName)
                tfConfig += processCidrmaps(domain.CidrMaps, importList.Cidrmaps, importList.Datacenters, resourceDomainName)
                tfConfig += processGeomaps(domain.GeographicMaps, importList.Geomaps, importList.Datacenters, resourceDomainName)
                tfConfig += processAsmaps(domain.AsMaps, importList.Asmaps, importList.Datacenters, resourceDomainName)
		fmt.Println(tfConfig)

		// build tfstate

		// build import script
                importScriptFilename := filepath.Join(tfWorkPath, domainName + "_resource_import.script")
                if _, err := os.Stat(importScriptFilename); err == nil {
                        // File exists. Bail
                        akamai.StopSpinnerFail()
                        return cli.NewExitError(color.RedString("Import script file already exists"), 1)
                }
		scriptContent, err := buildImportScript(importList, resourceDomainName)
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
	}
	
	return nil
}

func retrieveImportList() (*importListStruct, error) {

	// check if createImportList set. If so, already have .... 
	if createImportList {
		fmt.Println("Using fullImport")
		return fullImportList, nil
	}
	importListFilename := createImportListFilename()
	//debug
	fmt.Println(importListFilename)
        if _, err := os.Stat(importListFilename); err != nil {
  		return nil, err
        }
	importData, err := ioutil.ReadFile(importListFilename)
	if err != nil {
		return nil, err
	}
	//debug 
	fmt.Println(string(importData))
	importList := &importListStruct{}
	err = json.Unmarshal(importData, importList)
	if err != nil {
		return nil, err
	}

	return importList, nil

}

func createImportListFilename() string {

	return filepath.Join(tfWorkPath, domainName + "_import_list.json")

}

func buildImportScript(importList *importListStruct, resourceDomainName string) (string, error) {

        // build import script
        var import_prefix = "terraform import "
        var import_file = ""
        // domain
        if !checkForResource(domainResource, domainName) {
                import_file += import_prefix + domainResource + "." + resourceDomainName + " " + importList.Domain + "\n"
        }
        // datacenters
        for id, nickname := range importList.Datacenters {
                if !checkForResource(datacenterResource, nickname) {
                        import_file += import_prefix + datacenterResource+ "." + nickname + " " + importList.Domain+":" + strconv.Itoa(id) + "\n"
                }
       	} 
        // properties
        for name, _ := range importList.Properties {
                if !checkForResource(propertyResource, name) {
                        import_file += import_prefix + propertyResource+ "." + name + " " + importList.Domain+":" + name + "\n"
                }
        }
        // resources
        for name, _ := range importList.Resources {
                if !checkForResource(resourceResource, name) {
                        import_file += import_prefix + resourceResource+ "." + name + " " + importList.Domain+":" + name + "\n"
                }
        }
        // cidrmaps
        for name, _ := range importList.Cidrmaps {
                if !checkForResource(cidrResource, name) {
                        import_file += import_prefix + cidrResource+ "." + name + " " + importList.Domain+":" + name + "\n"
                }
        }
        // geomaps
        for name, _ := range importList.Geomaps {
                if !checkForResource(geoResource, name) {
                        import_file += import_prefix + geoResource+ "." + name + " " + importList.Domain+":" + name + "\n"
                }
        }
        // asmaps
        for name, _ := range importList.Asmaps {
                if !checkForResource(asResource, name) {
                        import_file += import_prefix + asResource+ "." + name + " " + importList.Domain+":" + name + "\n"
                }
        }
	
	return import_file, nil

}









