// Package mtlskeystore contains code for exporting mTLS Keystore
package mtlskeystore

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/mtlskeystore"
	"github.com/akamai/cli-terraform/v2/pkg/edgegrid"
	"github.com/akamai/cli-terraform/v2/pkg/templates"
	"github.com/akamai/cli-terraform/v2/pkg/tools"
	"github.com/akamai/cli/v2/pkg/color"
	"github.com/akamai/cli/v2/pkg/terminal"
	"github.com/urfave/cli/v2"
)

type (

	// TFData represents the data used in mTLS Keystore.
	TFData struct {
		Certificate TFCertificate
		Section     string
	}

	// TFCertificate represents the data used for exporting mTLS Keystore client certificate.
	TFCertificate struct {
		ResourceName       string
		Name               string
		ID                 string
		ContractID         string
		GroupID            int64
		Geography          string
		KeyAlgorithm       string
		NotificationEmails []string
		SecureNetwork      string
		Subject            string
		Signer             string
		Versions           []TFClientCertificateVersion
	}

	// TFClientCertificateVersion represents the versions data of the client certificate needed for the configuration.
	TFClientCertificateVersion struct {
		Version     int64
		CreatedDate string
	}

	createCertificateParams struct {
		id                int64
		groupID           string
		contractID        string
		configSection     string
		client            mtlskeystore.MTLSKeystore
		templateProcessor templates.TemplateProcessor
	}
)

var (
	//go:embed templates/*
	templateFiles embed.FS

	additionalFunctions = tools.DecorateWithMultilineHandlingFunctions(map[string]any{
		"getLastIndex": tools.GetLastIndex,
	})

	// ErrFetchingClientCertificate is returned when client certificate could not be fetched.
	ErrFetchingClientCertificate = errors.New("error fetching client certificate")
	// ErrFetchingClientCertificateVersions is returned when client certificate versions could not be fetched.
	ErrFetchingClientCertificateVersions = errors.New("error fetching client certificate versions")
	// ErrSavingFiles is returned when an issue with processing templates occurs.
	ErrSavingFiles = errors.New("error saving terraform project files")
)

// CmdCreateCertificate is an entrypoint to export-mtls-keystore command.
func CmdCreateCertificate(c *cli.Context) error {

	// tfWorkPath is a target directory for generated terraform resources.
	var tfWorkPath = "./"
	if c.IsSet("tfworkpath") {
		tfWorkPath = c.String("tfworkpath")
	}
	mTLSKeyStorePath := filepath.Join(tfWorkPath, "mtlskeystore.tf")
	variablesPath := filepath.Join(tfWorkPath, "variables.tf")
	importPath := filepath.Join(tfWorkPath, "import.sh")
	if err := tools.CheckFiles(mTLSKeyStorePath, variablesPath, importPath); err != nil {
		return cli.Exit(color.RedString(err.Error()), 1)
	}

	idString := c.Args().Get(0)
	id, err := strconv.ParseInt(idString, 10, 64)
	if err != nil {
		return cli.Exit(color.RedString("Invalid client certificate ID: %s", idString), 1)
	}

	groupID := c.Args().Get(1)
	contractID := c.Args().Get(2)

	params := createCertificateParams{
		id:            id,
		groupID:       groupID,
		contractID:    contractID,
		configSection: edgegrid.GetEdgercSection(c),
		client:        mtlskeystore.Client(edgegrid.GetSession(c.Context)),
		templateProcessor: templates.FSTemplateProcessor{
			TemplatesFS: templateFiles,
			TemplateTargets: map[string]string{
				"mtlskeystore.tmpl": mTLSKeyStorePath,
				"variables.tmpl":    variablesPath,
				"imports.tmpl":      importPath,
			},
			AdditionalFuncs: additionalFunctions,
		},
	}

	if err := createCertificate(c.Context, params); err != nil {
		return cli.Exit(color.RedString("Error exporting mtls keystore certificate: %s", err), 1)
	}
	return nil
}

func createCertificate(ctx context.Context, params createCertificateParams) (e error) {
	term := terminal.Get(ctx)
	defer func() {
		if e != nil {
			term.Spinner().Fail()
		}
	}()
	term.Spinner().Start("Fetching client certificate")
	cert, err := params.client.GetClientCertificate(ctx, mtlskeystore.GetClientCertificateRequest{
		CertificateID: params.id,
	})
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFetchingClientCertificate, err)
	}
	term.Spinner().OK()

	var versions *mtlskeystore.ListClientCertificateVersionsResponse
	term.Spinner().Start("Fetching client certificate versions")
	versions, err = params.client.ListClientCertificateVersions(ctx, mtlskeystore.ListClientCertificateVersionsRequest{
		CertificateID: params.id,
	})
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFetchingClientCertificateVersions, err)
	}
	term.Spinner().OK()

	// For Akamai certificate, export only if at least one version is not pending deletion.
	if cert.Signer == string(mtlskeystore.SignerAkamai) && versions != nil {
		if err := checkAkamaiCertVersionsStatus(versions); err != nil {
			return err
		}
	}

	term.Spinner().Start("Extracting data")
	tfData, err := populateTFData(params, cert, versions)
	if err != nil {
		return fmt.Errorf("error populating terraform data: %w", err)
	}
	term.Spinner().OK()

	term.Spinner().Start("Saving TF configurations")
	if err = params.templateProcessor.ProcessTemplates(tfData); err != nil {
		return fmt.Errorf("%w: %s", ErrSavingFiles, err)
	}
	term.Spinner().OK()

	term.Printf("Terraform configuration for client certificate id '%d' was saved successfully\n", params.id)

	return nil
}

func populateTFData(param createCertificateParams, cert *mtlskeystore.GetClientCertificateResponse, versions *mtlskeystore.ListClientCertificateVersionsResponse) (TFData, error) {
	var tfVersions []TFClientCertificateVersion
	// For 'AKAMAI' versions are always nil, as we don't generate configuration for them.
	if versions != nil {
		for _, v := range versions.Versions {
			// We only add the versions that have nil VersionAlias and status different then pending deletion.
			if v.VersionAlias == nil && v.Status != string(mtlskeystore.CertificateVersionStatusDeletePending) {
				tfVersions = append(tfVersions, TFClientCertificateVersion{
					Version:     v.Version,
					CreatedDate: strings.TrimSuffix(v.CreatedDate.Format(time.RFC3339), "Z"),
				})
			}
		}
	}
	if len(tfVersions) == 0 && cert.Signer == string(mtlskeystore.SignerThirdParty) {
		return TFData{}, fmt.Errorf("certificate with ID '%d' has no versions or the versions are pending deletion", cert.CertificateID)
	}

	tfData := TFData{
		Section: param.configSection,
		Certificate: TFCertificate{
			Name:               cert.CertificateName,
			ResourceName:       strings.ReplaceAll(cert.CertificateName, "-", "_"),
			ID:                 strconv.FormatInt(cert.CertificateID, 10),
			Geography:          cert.Geography,
			KeyAlgorithm:       cert.KeyAlgorithm,
			NotificationEmails: cert.NotificationEmails,
			SecureNetwork:      cert.SecureNetwork,
			Signer:             cert.Signer,
			Subject:            cert.Subject,
			Versions:           tfVersions,
		},
	}

	var ctr, grp string
	var err error
	if param.contractID != "" && param.groupID != "" {
		ctr = param.contractID
		grp = param.groupID
	} else {
		ctr, grp, err = extractContractAndGroup(cert.Subject)
		if err != nil {
			return TFData{}, fmt.Errorf("unable to extract group and contract from certificate subject: %w.\nRe-run with following arguments: <certificate_id>  <group_id> <contract_id>", err)
		}
	}

	grpInt64, err := strconv.ParseInt(grp, 10, 64)
	if err != nil {
		return TFData{}, fmt.Errorf("failed to convert group ID to int64: %w", err)
	}
	tfData.Certificate.ContractID = ctr
	tfData.Certificate.GroupID = grpInt64

	return tfData, nil
}

// checkAkamaiCertVersionsStatus returns an error if the Akamai certificate has no versions or all versions are pending deletion.
func checkAkamaiCertVersionsStatus(versions *mtlskeystore.ListClientCertificateVersionsResponse) error {
	for _, v := range versions.Versions {
		if v.Status != string(mtlskeystore.CertificateVersionStatusDeletePending) {
			return nil
		}
	}
	return fmt.Errorf("certificate has no versions or the versions are pending deletion")
}

func extractContractAndGroup(subject string) (string, string, error) {
	// Capture the part before required '/CN=' label.
	re := regexp.MustCompile(`\/([^\/]+)\/CN=`)
	matches := re.FindStringSubmatch(subject)
	if len(matches) < 2 {
		return "", "", fmt.Errorf("unexpected format: '%s'", subject)
	}
	parts := strings.Fields(matches[1])
	if len(parts) < 2 {
		return "", "", fmt.Errorf("no group or contract: '%s'", subject)
	}
	return parts[len(parts)-2], parts[len(parts)-1], nil
}
