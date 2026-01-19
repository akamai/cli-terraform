// Package domainownership contains code for exporting Domain Ownership.
package domainownership

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"path/filepath"
	"slices"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/domainownership"
	"github.com/akamai/cli-terraform/v2/pkg/edgegrid"
	"github.com/akamai/cli-terraform/v2/pkg/templates"
	"github.com/akamai/cli-terraform/v2/pkg/tools"
	"github.com/akamai/cli/v2/pkg/color"
	"github.com/akamai/cli/v2/pkg/terminal"
	"github.com/urfave/cli/v2"
)

type (

	// TFData represents the data used in Domain Ownership.
	TFData struct {
		Section                string
		ResourceName           string
		Domains                []TFDomain
		ImportKeyForDomains    string
		ImportKeyForValidation string
		EdgeRCPath             string
	}

	// TFDomain represents a domain in Terraform.
	TFDomain struct {
		// DomainName is the name of the domain.
		DomainName string
		// ValidationScope indicates the scope of the validation, either HOST, DOMAIN, or WILDCARD.
		ValidationScope string
		// ValidationMethod indicates the method of validation, either DNS_CNAME, DNS_TXT or HTTP.
		ValidationMethod string
		// Validated indicates whether the domain has been validated (DomainStatus is VALIDATED).
		Validated bool
	}

	parsedDomain struct {
		domain          string
		validationScope domainownership.ValidationScope
	}

	createDomainOwnershipParams struct {
		domains           string
		configSection     string
		edgercPath        string
		templateProcessor templates.TemplateProcessor
	}
)

var (
	//go:embed templates/*
	templateFiles embed.FS

	additionalFunctions = map[string]any{
		"getLastIndex": tools.GetLastIndex,
	}

	validationScopes = []domainownership.ValidationScope{
		domainownership.ValidationScopeDomain,
		domainownership.ValidationScopeHost,
		domainownership.ValidationScopeWildcard,
	}

	// ErrParsingDomains is returned when domains could not be parsed.
	ErrParsingDomains = errors.New("error parsing domains")
	// ErrSearchingDomains is returned when searching domains failed.
	ErrSearchingDomains = errors.New("error searching domains")
	// ErrSavingFiles is returned when saving files failed.
	ErrSavingFiles = errors.New("error saving files")
)

// CmdCreateDomainOwnership is an entrypoint to the export-domainownership.
func CmdCreateDomainOwnership(c *cli.Context) error {

	// tfWorkPath is a target directory for generated terraform resources.
	var tfWorkPath = "./"
	if c.IsSet("tfworkpath") {
		tfWorkPath = c.String("tfworkpath")
	}
	domainownershipPath := filepath.Join(tfWorkPath, "domainownership.tf")
	variablesPath := filepath.Join(tfWorkPath, "variables.tf")
	importPath := filepath.Join(tfWorkPath, "import.sh")
	if err := tools.CheckFiles(domainownershipPath, variablesPath, importPath); err != nil {
		return cli.Exit(color.RedString("%s", err.Error()), 1)
	}

	params := createDomainOwnershipParams{
		domains:       c.Args().Get(0),
		configSection: edgegrid.GetEdgercSection(c),
		edgercPath:    edgegrid.GetEdgercPath(c),
		templateProcessor: templates.FSTemplateProcessor{
			TemplatesFS: templateFiles,
			TemplateTargets: map[string]string{
				"domainownership.tmpl": domainownershipPath,
				"variables.tmpl":       variablesPath,
				"import.tmpl":          importPath,
			},
			AdditionalFuncs: additionalFunctions,
		},
	}

	if err := createDomainOwnership(c.Context, params, domainownership.Client(edgegrid.GetSession(c.Context))); err != nil {
		return cli.Exit(color.RedString("Error exporting domain ownership: %s", err), 1)
	}
	return nil
}

func createDomainOwnership(ctx context.Context, params createDomainOwnershipParams, client domainownership.DomainOwnership) (e error) {
	term := terminal.Get(ctx)
	msg := "Fetching domains\n"
	term.Spinner().Start(msg)
	defer func() {
		if e != nil {
			term.Spinner().Fail()
		}
	}()

	inputDomains, err := parseInputDomains(params.domains)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrParsingDomains, err)
	}

	requestDomains := make([]domainownership.Domain, 0)
	for _, inputDomain := range inputDomains {
		domain := inputDomain.domain
		validationScope := inputDomain.validationScope
		if validationScope != "" {
			requestDomains = append(requestDomains, domainownership.Domain{
				DomainName:      domain,
				ValidationScope: validationScope})
		} else {
			for _, scope := range validationScopes {
				requestDomains = append(requestDomains, domainownership.Domain{
					DomainName:      domain,
					ValidationScope: scope,
				})
			}
		}
	}

	foundDomains := make(map[string]domainownership.SearchDomainItem)
	// SearchDomains has a limit of 1000 domains per request.
	for domains := range slices.Chunk(requestDomains, 1000) {
		request := domainownership.SearchDomainsRequest{
			IncludeAll: true,
			Body: domainownership.SearchDomainsBody{
				Domains: domains,
			},
		}
		response, err := client.SearchDomains(ctx, request)
		if err != nil {
			return fmt.Errorf("%w: %w", ErrSearchingDomains, err)
		}

		for _, domain := range response.Domains {
			if domain.ValidationLevel == "FQDN" {
				foundDomains[domain.DomainName+":"+domain.ValidationScope] = domain
			}
		}
	}

	// If domain was specified without validation scope, we have to ensure that there is only one matching domain.
	domains := make([]domainownership.SearchDomainItem, 0, len(foundDomains))
	for _, parsedDomain := range inputDomains {
		domain := parsedDomain.domain
		scope := parsedDomain.validationScope
		if scope == "" {
			// count how many matching domains we have in the response.
			matchingDomains := make([]string, 0)
			for _, scope := range validationScopes {
				if foundDomain, ok := foundDomains[domain+":"+string(scope)]; ok {
					matchingDomains = append(matchingDomains, domain+":"+string(scope))
					// it could be added outside the loop when we know that there is only one match, but this way we simplify the processing.
					domains = append(domains, foundDomain)
				}
			}
			if len(matchingDomains) == 0 {
				return fmt.Errorf("%w: domain '%s' was not found or the validation level was not 'FQDN'", ErrParsingDomains, domain)
			}
			if len(matchingDomains) > 1 {
				return fmt.Errorf("%w: domain '%s' specified without validation scope has multiple matches: %v", ErrParsingDomains, domain, matchingDomains)
			}
		}
		if scope != "" {
			// ensure that specific domain with scope exists.
			foundDomain, found := foundDomains[domain+":"+string(scope)]
			if !found {
				return fmt.Errorf("%w: domain '%s' with validation scope '%s' does not exist or the validation level was not 'FQDN'", ErrParsingDomains, domain, scope)
			}
			domains = append(domains, foundDomain)
		}
	}

	tfData, warnings := populateTFData(params, domains, inputDomains[0].domain)
	if len(warnings) > 0 {
		for _, warning := range warnings {
			term.Printf(warning)
		}
	}
	term.Spinner().OK()

	term.Spinner().Start("Saving TF configurations ")
	if err = params.templateProcessor.ProcessTemplates(tfData); err != nil {
		return fmt.Errorf("%w: %s", ErrSavingFiles, err)
	}
	term.Spinner().OK()

	term.Printf("Terraform configuration for %d domain(s) was saved successfully\n", len(tfData.Domains))

	return nil
}

func parseInputDomains(domainsIn string) ([]parsedDomain, error) {
	result := make([]parsedDomain, 0)
	seenDomains := make(map[string]struct{})
	ds := strings.Split(domainsIn, ",")
	if len(ds) > 1000 {
		return nil, fmt.Errorf("maximum of 1000 domains can be processed at once")
	}
	for _, d := range ds {
		parts := strings.Split(d, ":")
		switch len(parts) {
		case 1: // domain only, no validationScope.
			domain := parts[0]
			if _, ok := seenDomains[domain+":"]; ok {
				return nil, fmt.Errorf("domain '%s' specified multiple times without validation scope", domain)
			}
			for _, scope := range validationScopes {
				if _, ok := seenDomains[domain+":"+string(scope)]; ok {
					return nil, fmt.Errorf("domain '%s' specified multiple times with and without validation scope", domain)
				}
			}
			seenDomains[domain+":"] = struct{}{}

			result = append(result, parsedDomain{domain: domain, validationScope: ""})

		case 2: // domain:validationScope.
			domain := parts[0]
			validationScope := domainownership.ValidationScope(parts[1])
			if validationScope != domainownership.ValidationScopeHost && validationScope != domainownership.ValidationScopeWildcard && validationScope != domainownership.ValidationScopeDomain {
				return nil, fmt.Errorf("invalid validation scope '%s' for domain '%s', it should be 'HOST', 'WILDCARD' or 'DOMAIN'", validationScope, domain)
			}
			if _, ok := seenDomains[domain+":"]; ok {
				return nil, fmt.Errorf("domain '%s' specified multiple times with and without validation scope", domain)
			}
			if _, ok := seenDomains[domain+":"+string(validationScope)]; ok {
				return nil, fmt.Errorf("domain '%s' with validation scope '%s' specified multiple times", domain, validationScope)
			}

			seenDomains[domain+":"+string(validationScope)] = struct{}{}
			result = append(result, parsedDomain{domain: domain, validationScope: validationScope})

		default:
			return nil, fmt.Errorf("invalid domain format for '%s', it should be 'domain' or 'domain:validation_scope' separated with comma", d)
		}
	}
	return result, nil
}

func populateTFData(params createDomainOwnershipParams, foundDomains []domainownership.SearchDomainItem, firstDomainName string) (TFData, []string) {
	domains := make([]TFDomain, 0)
	importKeyForDomains := make([]string, 0)
	importKeyForValidation := make([]string, 0)
	warnings := make([]string, 0)
	for _, domain := range foundDomains {
		var validationMethod string
		if domain.ValidationMethod != nil {
			validationMethod = *domain.ValidationMethod
		} else {
			warnings = append(warnings, fmt.Sprintf("[WARN] ValidationMethod is nil for domain %s with scope %s. Please complete this field before import.\n", domain.DomainName, domain.ValidationScope))
		}
		domains = append(domains, TFDomain{
			DomainName:       domain.DomainName,
			ValidationScope:  domain.ValidationScope,
			ValidationMethod: validationMethod,
			Validated:        domain.DomainStatus == "VALIDATED",
		})
		importKeyForDomains = append(importKeyForDomains, fmt.Sprintf("%s:%s", domain.DomainName, domain.ValidationScope))
		if domain.DomainStatus == "VALIDATED" {
			importKeyForValidation = append(importKeyForValidation, fmt.Sprintf("%s:%s", domain.DomainName, domain.ValidationScope))
		}
	}

	return TFData{
		Section:                params.configSection,
		EdgeRCPath:             params.edgercPath,
		ResourceName:           tools.TerraformName(firstDomainName),
		ImportKeyForDomains:    strings.Join(importKeyForDomains, ","),
		ImportKeyForValidation: strings.Join(importKeyForValidation, ","),
		Domains:                domains,
	}, warnings
}
