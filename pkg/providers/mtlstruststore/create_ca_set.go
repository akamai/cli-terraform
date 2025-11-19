// Package mtlstruststore contains code for exporting mTLS Truststore
package mtlstruststore

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/mtlstruststore"
	"github.com/akamai/cli-terraform/v2/pkg/edgegrid"
	"github.com/akamai/cli-terraform/v2/pkg/templates"
	"github.com/akamai/cli-terraform/v2/pkg/tools"
	"github.com/akamai/cli/v2/pkg/color"
	"github.com/akamai/cli/v2/pkg/terminal"
	"github.com/urfave/cli/v2"
)

type (

	// TFData represents the data used in mTLS Truststore
	TFData struct {
		CASet      TruststoreCASet
		EdgercPath string
		Section    string
	}

	// TruststoreCASet represents the data used for exporting mTLS Truststore CA set
	TruststoreCASet struct {
		Name               string
		ResourceName       string
		ID                 string
		Description        *string
		AllowInsecureSHA1  bool
		VersionDescription *string
		Certificates       []TruststoreCertificate
		Networks           []NetworkInfo
	}

	// TruststoreCertificate represents the data used for exporting mTLS Truststore certificate
	TruststoreCertificate struct {
		CertificatePEM string
		Description    *string
	}

	// NetworkInfo represents the data used for exporting mTLS Truststore CA set activation
	// on a specific network
	NetworkInfo struct {
		NetworkName                  string
		HasActivation                bool
		HasActivationOnLatestVersion bool
	}
)

var (
	//go:embed templates/*
	templateFiles embed.FS

	additionalFunctions = tools.DecorateWithMultilineHandlingFunctions(map[string]any{
		"getLastIndex": tools.GetLastIndex,
		"toUpper":      strings.ToUpper,
	})

	// ErrFindingCASet is returned when CA set could not be found
	ErrFindingCASet = errors.New("error finding CA set")
	// ErrFetchingCASet is returned when CA set could not be fetched
	ErrFetchingCASet = errors.New("error fetching CA set")
	// ErrFetchingCASetVersion is returned when CA set version could not be fetched
	ErrFetchingCASetVersion = errors.New("error fetching CA set version")
	// ErrSavingFiles is returned when an issue with processing templates occurs
	ErrSavingFiles = errors.New("error saving terraform project files")
)

// CmdCreateCASet is an entrypoint to the export-mtls-truststore command
func CmdCreateCASet(c *cli.Context) error {

	// tfWorkPath is a target directory for generated terraform resources
	var tfWorkPath = "./"
	if c.IsSet("tfworkpath") {
		tfWorkPath = c.String("tfworkpath")
	}
	mTLSTrustStorePath := filepath.Join(tfWorkPath, "mtlstruststore.tf")
	variablesPath := filepath.Join(tfWorkPath, "variables.tf")
	importPath := filepath.Join(tfWorkPath, "import.sh")
	if err := tools.CheckFiles(mTLSTrustStorePath, variablesPath, importPath); err != nil {
		return cli.Exit(color.RedString("%s", err.Error()), 1)
	}

	params := createCASetParams{
		name:          c.Args().Get(0),
		userVersion:   c.Int64("version"),
		edgercPath:    edgegrid.GetEdgercPath(c),
		configSection: edgegrid.GetEdgercSection(c),
		client:        mtlstruststore.Client(edgegrid.GetSession(c.Context)),
		templateProcessor: templates.FSTemplateProcessor{
			TemplatesFS: templateFiles,
			TemplateTargets: map[string]string{
				"mtlstruststore.tmpl": mTLSTrustStorePath,
				"variables.tmpl":      variablesPath,
				"imports.tmpl":        importPath,
			},
			AdditionalFuncs: additionalFunctions,
		},
	}

	if err := createCASet(c.Context, params); err != nil {
		return cli.Exit(color.RedString("Error exporting mtls truststore: %s", err), 1)
	}
	return nil
}

type createCASetParams struct {
	name              string
	userVersion       int64
	edgercPath        string
	configSection     string
	client            mtlstruststore.MTLSTruststore
	templateProcessor templates.TemplateProcessor
}

func createCASet(ctx context.Context, params createCASetParams) (e error) {

	// Defensive check: the version should be already validated by the CLI framework
	if params.userVersion < 0 {
		return fmt.Errorf("version must be a positive integer or zero for no version, got %d",
			params.userVersion)
	}

	term := terminal.Get(ctx)
	msg := "Fetching CA set with name " + params.name
	if params.userVersion != 0 {
		msg += fmt.Sprintf(" and version %d", params.userVersion)
	}
	term.Spinner().Start(msg)
	defer func() {
		if e != nil {
			term.Spinner().Fail()
		}
	}()

	id, err := findCASetID(ctx, params.client, params.name)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFindingCASet, err)
	}

	caSet, err := params.client.GetCASet(ctx, mtlstruststore.GetCASetRequest{
		CASetID: id,
	})
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFetchingCASet, err)
	}

	var version int64
	if params.userVersion != 0 {
		version = params.userVersion
	} else {
		if caSet.LatestVersion == nil {
			return fmt.Errorf("%w: CA set '%s' has no versions", ErrFetchingCASet, params.name)
		}
		version = *caSet.LatestVersion
	}

	caSetVersion, err := params.client.GetCASetVersion(ctx, mtlstruststore.GetCASetVersionRequest{
		CASetID: id,
		Version: version,
	})
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFetchingCASetVersion, err)
	}

	tfData := populateTFData(params.edgercPath, params.configSection, caSet, caSetVersion)
	term.Spinner().OK()

	term.Spinner().Start("Saving TF configurations ")
	if err = params.templateProcessor.ProcessTemplates(tfData); err != nil {
		return fmt.Errorf("%w: %s", ErrSavingFiles, err)
	}
	term.Spinner().OK()

	term.Printf("Terraform configuration for CA set '%s' was saved successfully\n", params.name)

	return nil
}

func findCASetID(ctx context.Context, client mtlstruststore.MTLSTruststore, caSetName string) (string, error) {
	caSets, err := client.ListCASets(ctx, mtlstruststore.ListCASetsRequest{
		CASetNamePrefix: caSetName,
	})
	if err != nil {
		return "", fmt.Errorf("failed to list CA sets: %w", err)
	}

	var matchingSets []mtlstruststore.CASetResponse
	for _, caSet := range caSets.CASets {
		if caSet.CASetStatus == "NOT_DELETED" && caSet.CASetName == caSetName {
			matchingSets = append(matchingSets, caSet)
		}
	}
	if len(matchingSets) == 0 {
		return "", fmt.Errorf("no CA set found with name '%s'", caSetName)
	}
	if len(matchingSets) > 1 {
		return "", fmt.Errorf("multiple CA sets found with name '%s'", caSetName)
	}
	return matchingSets[0].CASetID, nil
}

func versionEquals(a, b *int64) bool {
	return a != nil && b != nil && *a == *b
}

func populateTFData(edgercPath, configSection string, caSet *mtlstruststore.GetCASetResponse,
	caSetVersion *mtlstruststore.GetCASetVersionResponse) TFData {
	var certs []TruststoreCertificate
	for _, cert := range caSetVersion.Certificates {
		certs = append(certs, TruststoreCertificate{
			CertificatePEM: cert.CertificatePEM,
			Description:    cert.Description,
		})
	}

	var networks []NetworkInfo

	networks = append(networks, NetworkInfo{
		NetworkName:                  "staging",
		HasActivation:                caSet.StagingVersion != nil,
		HasActivationOnLatestVersion: versionEquals(caSet.StagingVersion, caSet.LatestVersion),
	})

	networks = append(networks, NetworkInfo{
		NetworkName:                  "production",
		HasActivation:                caSet.ProductionVersion != nil,
		HasActivationOnLatestVersion: versionEquals(caSet.ProductionVersion, caSet.LatestVersion),
	})

	// CASetName: "Allowed characters are alphanumerics (a-z, A-Z, 0-9), underscore (_), hyphen (-)
	// and period (.) with no three consecutive periods. Length must be between 3 and 64 characters."
	// Hyphens are allowed in resource names, but we need to replace periods.
	resp := TFData{
		EdgercPath: edgercPath,
		Section:    configSection,
		CASet: TruststoreCASet{
			Name:               caSet.CASetName,
			ResourceName:       strings.ReplaceAll(caSet.CASetName, ".", "_"),
			ID:                 caSet.CASetID,
			Description:        caSet.Description,
			AllowInsecureSHA1:  caSetVersion.AllowInsecureSHA1,
			VersionDescription: caSetVersion.Description,
			Certificates:       certs,
			Networks:           networks,
		},
	}
	return resp
}
