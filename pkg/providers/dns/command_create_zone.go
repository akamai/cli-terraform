// Package dns contains code for exporting the DNS configuration.
package dns

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/dns"
	"github.com/akamai/cli-terraform/pkg/edgegrid"
	"github.com/akamai/cli-terraform/pkg/tools"
	"github.com/akamai/cli/pkg/color"
	"github.com/akamai/cli/pkg/terminal"
	"github.com/urfave/cli/v2"
)

// Types contains list of Name Types to organize types by name
type Types []string

// Import List Struct
type zoneImportListStruct struct {
	Zone       string
	RecordSets map[string]Types // zone record sets grouped by name
}

type configStruct struct {
	fetchConfig            fetchConfigStruct
	tfWorkPath             string
	shouldCreateImportList bool
	createConfig           bool
	recordNames            []string
	importScript           bool
}

type fetchConfigStruct struct {
	ConfigOnly bool
	ModSegment bool
	NamesOnly  bool
}

var zoneName string
var contractID string

var fullZoneImportList *zoneImportListStruct
var fullZoneConfigMap map[string]Types

// work defs
var moduleFolder = "modules"

// text for root module construction
var zoneTFFileHandle *os.File
var zoneTFConfig = ""

// CmdCreateZone is an entrypoint to create-zone command
func CmdCreateZone(c *cli.Context) error {
	ctx := c.Context
	log.SetOutput(io.Discard)

	sess := edgegrid.GetSession(ctx)
	configDNS := dns.Client(sess)

	// uppercase characters cause issues with TF and the generated config
	zoneName = strings.ToLower(c.Args().Get(0))

	configuration := setConfiguration(c)

	term := terminal.Get(ctx)
	fmt.Println("Configuring Zone")
	zoneObject, err := configDNS.GetZone(ctx, dns.GetZoneRequest{
		Zone: zoneName,
	})
	if err != nil {
		term.Spinner().Fail()
		fmt.Println("Error: " + err.Error())
		return cli.Exit(color.RedString("Zone retrieval failed"), 1)
	}
	contractID = zoneObject.ContractID // grab for use later
	// normalize zone name for zone resource name
	resourceZoneName := normalizeResourceName(zoneName)
	if configuration.shouldCreateImportList {
		err := createImportList(ctx, term, configDNS, resourceZoneName, configuration)
		if err != nil {
			return err
		}
		term.Spinner().OK()
	}

	if configuration.createConfig {
		// Read in resources list
		zoneImportList, err := retrieveZoneImportList(resourceZoneName, configuration)
		if err != nil {
			term.Spinner().Fail()
			return cli.Exit(color.RedString("Failed to read json zone resources file"), 1)
		}
		// if segmenting record sets by name, make sure module folder exists
		if configuration.fetchConfig.ModSegment {
			modulePath := filepath.Join(configuration.tfWorkPath, moduleFolder)
			if !createDirectory(modulePath) {
				term.Spinner().Fail()
				return cli.Exit(color.RedString("Failed to create modules folder."), 1)
			}
		}
		term.Spinner().Start("Creating zone configuration file ")
		err = createZoneConfigFile(ctx, zoneImportList, resourceZoneName, zoneObject, configDNS, configuration)
		if err != nil {
			term.Spinner().Fail()
			return err
		}

		err = createDNSVarsConfig(term, configuration.tfWorkPath)
		if err != nil {
			return err
		}
		term.Spinner().OK()
	}

	if configuration.importScript {
		term.Spinner().Start("Creating zone import script file")
		err := createImportScript(resourceZoneName, term, configuration)
		if err != nil {
			term.Spinner().Fail()
			return err
		}
		term.Spinner().OK()
	}

	fmt.Println("Zone configuration completed")

	return nil
}

func createImportList(ctx context.Context, term terminal.Terminal, configDNS dns.DNS, resourceZoneName string, configuration configStruct) error {
	term.Spinner().Start("Inventorying zone and recordsets ")
	recordSets, err := inventorZone(ctx, configDNS, configuration)
	if err != nil {
		term.Spinner().Fail()
		fmt.Println("Error: " + err.Error())
		return err
	}
	term.Spinner().OK()

	term.Spinner().Start("Creating Zone Resources list file ")
	err = createZoneResourceListFile(resourceZoneName, recordSets, configuration.tfWorkPath)
	if err != nil {
		term.Spinner().Fail()
		return err
	}
	return nil
}

func setConfiguration(c *cli.Context) configStruct {
	var executionConfig = configStruct{
		tfWorkPath: "./",
	}

	if c.IsSet("tfworkpath") {
		executionConfig.tfWorkPath = c.String("tfworkpath")
	}
	executionConfig.tfWorkPath = filepath.FromSlash(executionConfig.tfWorkPath)
	if c.IsSet("resources") {
		executionConfig.shouldCreateImportList = true
	}
	if c.IsSet("createconfig") {
		executionConfig.createConfig = true
	}
	if c.IsSet("configonly") {
		executionConfig.fetchConfig.ConfigOnly = true
	}
	if c.IsSet("namesonly") {
		executionConfig.fetchConfig.NamesOnly = true
	}
	if c.IsSet("recordname") {
		executionConfig.recordNames = c.StringSlice("recordname")
	}
	if c.IsSet("segmentconfig") {
		executionConfig.fetchConfig.ModSegment = true
	}
	if c.IsSet("importscript") {
		executionConfig.importScript = true
	}

	return executionConfig
}

func createZoneConfigFile(ctx context.Context, zoneImportList *zoneImportListStruct, resourceZoneName string, zoneObject *dns.GetZoneResponse, configDNS dns.DNS, configuration configStruct) (err error) {
	// see if configuration file already exists and exclude any resources already represented.
	var configImportList *zoneImportListStruct
	var zoneTypeMap map[string]map[string]bool

	zoneTFFileHandle, zoneTFConfig, err = openZoneConfigFile(resourceZoneName, configuration.tfWorkPath)
	if err != nil {
		return cli.Exit(color.RedString("Failed to open/create zone config file."), 1)
	}
	configImportList, zoneTypeMap = reconcileZoneResourceTargets(zoneImportList, resourceZoneName, zoneTFConfig)

	defer func(zoneTFFileHandle *os.File) {
		if e := zoneTFFileHandle.Close(); e != nil {
			err = e
		}
	}(zoneTFFileHandle)
	fileUtils := fileUtilsProcessor{}

	err = calculateTfConfig(ctx, zoneObject, resourceZoneName, fileUtils, configuration)
	if err != nil {
		return err
	}
	err = fileUtils.appendRootModuleTF(zoneTFConfig)
	if err != nil {
		fmt.Println(err.Error())
		return cli.Exit(color.RedString("Failed. Couldn't write to zone config"), 1)
	}

	// process RecordSets.
	fullZoneConfigMap, err = processRecordSets(ctx, configDNS, configImportList.Zone, resourceZoneName, zoneTypeMap, fileUtils, configuration)
	if err != nil {
		return cli.Exit(color.RedString("Failed to process recordsets."), 1)
	}
	// Save config map for import script generation
	resourceConfigFilename := createResourceConfigFilename(resourceZoneName, configuration.tfWorkPath)

	return saveResourceConfigFile(resourceConfigFilename)
}

func calculateTfConfig(ctx context.Context, zoneObject *dns.GetZoneResponse, resourceZoneName string, fileUtils fileUtilsProcessor, config configStruct) error {
	// build tf file if none
	var err error
	if len(zoneTFConfig) > 0 {
		if strings.Contains(zoneTFConfig, "module") && strings.Contains(zoneTFConfig, "zonename") {
			if !config.fetchConfig.ModSegment {
				// already have a top level zone config and it's modularized!
				return cli.Exit(color.RedString("Failed. Existing zone config is modularized"), 1)
			}
		} else if config.fetchConfig.ModSegment {
			// already have a top level zone config and it's not modularized!
			return cli.Exit(color.RedString("Failed. Existing zone config is not modularized"), 1)
		}
	} else {
		// if tf pre-existed, zone has to exist by definition
		zoneTFConfig, err = processZone(ctx, zoneObject, resourceZoneName, config.fetchConfig.ModSegment, fileUtils, config.tfWorkPath)
		if err != nil {
			fmt.Println(err.Error())
			return cli.Exit(color.RedString("Failed. Couldn't initialize zone config"), 1)
		}
	}
	return nil
}

func saveResourceConfigFile(resourceConfigFilename string) (err error) {
	resourceConfigJSON, err := json.MarshalIndent(&fullZoneConfigMap, "", "  ")
	if err != nil {
		return cli.Exit(color.RedString("Unable to generate json formatted zone config"), 1)
	}
	f, err := os.Create(resourceConfigFilename)
	if err != nil {
		return cli.Exit(color.RedString("Unable to create resource config file"), 1)
	}
	defer func(f *os.File) {
		if e := f.Close(); e != nil {
			err = e
		}
	}(f)

	_, err = f.WriteString(string(resourceConfigJSON))
	if err != nil {
		return cli.Exit(color.RedString("Unable to write zone resource config file"), 1)
	}
	err = f.Sync()
	if err != nil {
		return err
	}
	return nil
}

func createDNSVarsConfig(term terminal.Terminal, tfWorkPath string) (err error) {
	// Need to create dnsvars.tf dependency
	dnsVarsFileName := filepath.Join(tfWorkPath, "dnsvars.tf")
	dnsVarsHandle, err := os.Create(dnsVarsFileName)
	//}
	if err != nil {
		term.Spinner().Fail()
		return cli.Exit(color.RedString("Unable to create dnsvars config file"), 1)
	}
	defer func(dnsVarsHandle *os.File) {
		if e := dnsVarsHandle.Close(); e != nil {
			err = e
		}
	}(dnsVarsHandle)
	_, err = dnsVarsHandle.WriteString(fmt.Sprintf(useTemplate(nil, "dnsvars.tmpl", true), contractID))
	if err != nil {
		term.Spinner().Fail()
		return cli.Exit(color.RedString("Unable to write dnsvars config file"), 1)
	}

	return dnsVarsHandle.Sync()
}

func createImportScript(resourceZoneName string, term terminal.Terminal, configuration configStruct) (err error) {
	fullZoneConfigMap, _ = retrieveZoneResourceConfig(resourceZoneName, configuration)
	importScriptFilename := filepath.Join(configuration.tfWorkPath, resourceZoneName+"_resource_import.script")
	if _, err := os.Stat(importScriptFilename); err == nil {
		term.Spinner().OK()
	}
	scriptContent, err := buildZoneImportScript(zoneName, fullZoneConfigMap, resourceZoneName)

	if err != nil {
		return cli.Exit(color.RedString("Import script content generation failed"), 1)
	}
	f, err := os.Create(importScriptFilename)
	if err != nil {
		return cli.Exit(color.RedString("Unable to create import script file"), 1)
	}
	defer func(f *os.File) {
		if e := f.Close(); e != nil {
			err = e
		}
	}(f)
	_, err = f.WriteString(scriptContent)
	if err != nil {
		return cli.Exit(color.RedString("Unable to write import script file"), 1)
	}
	err = f.Sync()

	return err
}

func createZoneResourceListFile(resourceZoneName string, recordSets map[string]Types, tfWorkPath string) error {
	importListFilename := createImportListFilename(resourceZoneName, tfWorkPath)
	if _, err := os.Stat(importListFilename); err == nil {
		return cli.Exit(color.RedString("Resource list file exists. Remove to continue."), 1)
	}
	fullZoneImportList = &zoneImportListStruct{}
	fullZoneImportList.Zone = zoneName
	fullZoneImportList.RecordSets = recordSets
	err := saveImportListToFile(importListFilename)
	if err != nil {
		return err
	}
	return nil
}

func saveImportListToFile(importListFilename string) (err error) {
	importListJSON, err := json.MarshalIndent(fullZoneImportList, "", "  ")
	if err != nil {
		return cli.Exit(color.RedString("Unable to generate json formatted zone resource list"), 1)
	}
	f, err := os.Create(importListFilename)
	if err != nil {
		return cli.Exit(color.RedString("Unable to create zone resources file"), 1)
	}
	defer func(f *os.File) {
		if e := f.Close(); e != nil {
			err = e
		}
	}(f)
	_, err = f.WriteString(string(importListJSON))
	if err != nil {
		return cli.Exit(color.RedString("Unable to write zone resources file"), 1)
	}
	err = f.Sync()

	return err
}

func inventorZone(ctx context.Context, configDNS dns.DNS, configuration configStruct) (map[string]Types, error) {
	recordSets := make(map[string]Types)
	// Retrieve all zone names
	if len(configuration.recordNames) == 0 {
		recordsetNames, err := configDNS.GetZoneNames(ctx, dns.GetZoneNamesRequest{
			Zone: zoneName,
		})
		if err != nil {
			return nil, cli.Exit(color.RedString("Zone Name retrieval failed"), 1)
		}
		configuration.recordNames = recordsetNames.Names
	}
	for _, zName := range configuration.recordNames {
		if configuration.fetchConfig.NamesOnly {
			recordSets[zName] = make([]string, 0)
		} else {
			nameTypesResp, err := configDNS.GetZoneNameTypes(ctx, dns.GetZoneNameTypesRequest{
				ZoneName: zName,
				Zone:     zoneName,
			})
			if err != nil {
				return nil, cli.Exit(color.RedString("Zone Name types retrieval failed"), 1)
			}
			recordSets[zName] = nameTypesResp.Types
		}
	}
	return recordSets, nil
}

// Utility function to create full resource config file path
func createResourceConfigFilename(resourceName, tfWorkPath string) string {
	return filepath.Join(tfWorkPath, resourceName+"_zoneconfig.json")
}

// Utility function to create named module path
func createNamedModulePath(modName, tfWorkPath string) string {
	return filepath.Join(tfWorkPath, moduleFolder, normalizeResourceName(modName))
}

// Utility func to create a directory
func createDirectory(dirName string) bool {

	stat, err := os.Stat(dirName)
	if err == nil && stat.IsDir() {
		return true
	}
	if os.IsNotExist(err) {
		errDir := os.MkdirAll(dirName, 0755)
		if errDir != nil {
			panic(err)
		}
		return true
	}
	if stat.Mode().IsRegular() {
		return false
	}

	return false
}

func buildZoneImportScript(zone string, zoneConfigMap map[string]Types, resourceName string) (string, error) {
	data := ImportData{
		Zone:          zone,
		ZoneConfigMap: zoneConfigMap,
		ResourceName:  resourceName,
	}
	return useTemplate(&data, "import-script.tmpl", true), nil
}

// remove any resources already present in existing zone tf configuration
func reconcileZoneResourceTargets(zoneImportList *zoneImportListStruct, zoneName, tfConfig string) (*zoneImportListStruct, map[string]map[string]bool) {

	zoneTypeMap := make(map[string]map[string]bool)
	// populate zoneTypeMap

	// need walk through each resource type
	for zName, typeList := range zoneImportList.RecordSets {
		typeMap := make(map[string]bool)
		revisedTypeList := make([]string, 0, len(typeList))
		for _, ntype := range typeList {
			normalName := createUniqueRecordsetName(zoneName, zName, ntype)
			if !strings.Contains(tfConfig, `"`+normalName+`"`) {
				typeMap[ntype] = true
				revisedTypeList = append(revisedTypeList, ntype)
			} else {
				fmt.Println("Recordset resource " + normalName + " found in existing tf file")
			}
		}
		zoneImportList.RecordSets[zName] = revisedTypeList
		zoneTypeMap[zName] = typeMap
	}

	return zoneImportList, zoneTypeMap
}

func openZoneConfigFile(zoneName string, tfWorkPath string) (*os.File, string, error) {
	tfFilename := tools.CreateTFFilename(zoneName, tfWorkPath)
	var tfHandle *os.File
	tfHandle, err := os.OpenFile(tfFilename, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil && err != io.EOF {
		fmt.Println(err.Error())
		return nil, "", err
	}
	tfInfo, err := os.Stat(tfFilename)
	if err != nil {
		fmt.Println(err.Error())
		return nil, "", err
	}
	tfScratch := make([]byte, tfInfo.Size())
	charsRead, err := tfHandle.Read(tfScratch)
	if err != nil && err != io.EOF {
		fmt.Println(err.Error())
		return nil, "", err
	}
	_, err = tfHandle.Seek(0, 0)
	if err != nil {
		fmt.Println(err.Error())
		return nil, "", err
	}
	tfConfig := ""
	if charsRead > 0 {
		tfConfig = string(tfScratch[0 : charsRead-1])
	}

	return tfHandle, tfConfig, nil
}

func retrieveZoneImportList(rscName string, configuration configStruct) (*zoneImportListStruct, error) {
	// check if shouldCreateImportList set. If so, already have ....
	if configuration.shouldCreateImportList {
		return fullZoneImportList, nil
	}
	if configuration.fetchConfig.ConfigOnly {
		fullZoneImportList := &zoneImportListStruct{Zone: zoneName}
		fullZoneImportList.RecordSets = make(map[string]Types)
		return fullZoneImportList, nil
	}
	importListFilename := createImportListFilename(rscName, configuration.tfWorkPath)
	if _, err := os.Stat(importListFilename); err != nil {
		return nil, err
	}
	importData, err := os.ReadFile(importListFilename)
	if err != nil {
		return nil, err
	}
	importList := &zoneImportListStruct{}
	err = json.Unmarshal(importData, importList)
	if err != nil {
		return nil, err
	}

	return importList, nil
}

func retrieveZoneResourceConfig(rscName string, config configStruct) (map[string]Types, error) {
	configList := make(map[string]Types)
	// check if createConfig set. If so, already have ....
	if config.createConfig {
		return fullZoneConfigMap, nil
	}
	resourceConfigFilename := createResourceConfigFilename(rscName, config.tfWorkPath)
	if _, err := os.Stat(resourceConfigFilename); err != nil {
		return configList, err
	}
	configData, err := os.ReadFile(resourceConfigFilename)
	if err != nil {
		return configList, err
	}
	err = json.Unmarshal(configData, &configList)
	if err != nil {
		return configList, err
	}

	return configList, nil
}
