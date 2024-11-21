// Package gtm contains code for exporting the global traffic manager configuration.
package gtm

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/gtm"
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
		CNameCoalescingEnabled      bool
		LoadFeedback                bool
		DefaultSSLClientCertificate string
		EndUserMappingEnabled       bool
		Section                     string
		SignAndServe                bool
		SignAndServeAlgorithm       string
		Datacenters                 []TFDatacenterData
		DefaultDatacenters          []TFDatacenterData
		Resources                   []gtm.Resource
		CIDRMaps                    []gtm.CIDRMap
		GeoMaps                     []gtm.GeoMap
		ASMaps                      []gtm.ASMap
		Properties                  []gtm.Property
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

var additionalFunctions = tools.DecorateWithMultilineHandlingFunctions(map[string]any{
	"normalize":    normalizeResourceName,
	"toUpper":      strings.ToUpper,
	"isDefaultDC":  isDefaultDatacenter,
	"escapeString": tools.EscapeQuotedStringLit,
})

var (
	subWithUnderscoreRegexp               = regexp.MustCompile(`[^\w-_]`)
	mustStartWithLetterOrUnderscoreRegexp = regexp.MustCompile("^[^a-zA-Z_]")
	// ErrFetchingDomain is returned when fetching domain fails
	ErrFetchingDomain = errors.New("unable to fetch domain with given name")
)

// CmdCreateDomain is an entrypoint to create-domain command
func CmdCreateDomain(c *cli.Context) error {
	ctx := c.Context
	sess := edgegrid.GetSession(ctx)
	client := gtm.Client(sess)

	// tfWorkPath is a target directory for generated terraform resources
	var tfWorkPath = "./"
	if c.IsSet("tfworkpath") {
		tfWorkPath = c.String("tfworkpath")
	}

	datacentersPath := filepath.Join(tfWorkPath, "datacenters.tf")
	domainPath := filepath.Join(tfWorkPath, "domain.tf")
	importPath := filepath.Join(tfWorkPath, "import.sh")
	mapsPath := filepath.Join(tfWorkPath, "maps.tf")
	propertiesPath := filepath.Join(tfWorkPath, "properties.tf")
	resourcesPath := filepath.Join(tfWorkPath, "resources.tf")
	variablesPath := filepath.Join(tfWorkPath, "variables.tf")

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
		AdditionalFuncs: additionalFunctions,
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

	if _, err := term.Writeln("Configuring Domain"); err != nil {
		return err
	}

	term.Spinner().Start(fmt.Sprintf("Fetching domain %s", domainName))
	domain, err := client.GetDomain(ctx, gtm.GetDomainRequest{
		DomainName: domainName,
	})
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
		DefaultSSLClientPrivateKey:  domain.DefaultSSLClientPrivateKey,
		DefaultErrorPenalty:         domain.DefaultErrorPenalty,
		CNameCoalescingEnabled:      domain.CNameCoalescingEnabled,
		LoadFeedback:                domain.LoadFeedback,
		DefaultSSLClientCertificate: domain.DefaultSSLClientCertificate,
		EndUserMappingEnabled:       domain.EndUserMappingEnabled,
		Resources:                   domain.Resources,
		CIDRMaps:                    domain.CIDRMaps,
		GeoMaps:                     domain.GeographicMaps,
		ASMaps:                      domain.ASMaps,
		Properties:                  domain.Properties,
		SignAndServe:                domain.SignAndServe,
	}

	if domain.SignAndServeAlgorithm != nil {
		tfDomainData.SignAndServeAlgorithm = *domain.SignAndServeAlgorithm
	}

	tfDomainData.getDatacenters(domain)
	term.Spinner().OK()

	term.Spinner().Start("Saving TF configurations")
	if err := templateProcessor.ProcessTemplates(tfDomainData); err != nil {
		term.Spinner().Fail()
		return err
	}
	term.Spinner().OK()

	if _, err = term.Writeln(fmt.Sprintf("Terraform configuration for policy '%s' was saved successfully\n", domain.Name)); err != nil {
		return err
	}

	return nil
}

func (d *TFDomainData) getDatacenters(domain *gtm.GetDomainResponse) {
	d.Datacenters = make([]TFDatacenterData, 0)
	d.DefaultDatacenters = make([]TFDatacenterData, 0)
	for _, dc := range domain.Datacenters {
		if isDefaultDatacenter(dc.DatacenterID) {
			d.DefaultDatacenters = append(d.DefaultDatacenters, TFDatacenterData{Nickname: dc.Nickname, ID: dc.DatacenterID})
		} else {
			d.Datacenters = append(d.Datacenters, TFDatacenterData{
				ID:                            dc.DatacenterID,
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

func isDefaultDatacenter(id int) bool {
	_, ok := defaultDCs[id]
	return ok
}
