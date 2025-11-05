// Package cloudcertificates contains code for exporting Cloud Certificate
package cloudcertificates

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/cloudcertificates"
	"github.com/akamai/cli-terraform/v2/pkg/edgegrid"
	"github.com/akamai/cli-terraform/v2/pkg/templates"
	"github.com/akamai/cli-terraform/v2/pkg/tools"
	"github.com/akamai/cli/v2/pkg/color"
	"github.com/akamai/cli/v2/pkg/terminal"
	"github.com/urfave/cli/v2"
)

// The suffix added by CCM to the certificate name when a certificate is renewed.
const renewedNameSuffix = ".renewed."

// The date format used in the renewed certificate name in format <base_name>.renewed.YYYY-MM-DD.
const renewedNameDateLayout = "2006-01-02"

// maxPageSize defines the maximum number of items to fetch per page when listing certificates.
const maxPageSize = 100

type (
	// TFData represents the data used in Cloud Certificate
	TFData struct {
		Certificate TFCertificate
		EdgercPath  string
		Section     string
	}

	// TFCertificate represents the data used for exporting Cloud Certificate.
	TFCertificate struct {
		ID                   string
		ContractID           string
		BaseName             string
		KeyType              string
		KeySize              string
		SecureNetwork        string
		SANs                 []string
		CertificateStatus    string
		Subject              *cloudcertificates.Subject
		SignedCertificatePEM *string
		TrustChainPEM        *string
		ResourceName         string
	}

	createCloudCertificateParams struct {
		name              string
		edgercPath        string
		configSection     string
		client            cloudcertificates.CloudCertificates
		templateProcessor templates.TemplateProcessor
	}
)

var (
	//go:embed templates/*
	templateFiles embed.FS

	additionalFunctions = tools.DecorateWithMultilineHandlingFunctions(map[string]any{})

	// ErrListingCloudCertificates is returned when listing cloud certificates fails.
	ErrListingCloudCertificates = errors.New("error listing cloud certificates")
	// ErrFetchingCloudCertificate is returned when fetching cloud certificate fails.
	ErrFetchingCloudCertificate = errors.New("error fetching cloud certificate")
	// ErrSavingFiles is returned when an issue with processing templates occurs.
	ErrSavingFiles = errors.New("error saving terraform project files")
)

// CmdCreateCloudCertificate is an entrypoint to the export-cloudcertificate command.
func CmdCreateCloudCertificate(c *cli.Context) error {
	// tfWorkPath is a target directory for generated terraform resources.
	var tfWorkPath = "./"
	if c.IsSet("tfworkpath") {
		tfWorkPath = c.String("tfworkpath")
	}
	cloudCertificatePath := filepath.Join(tfWorkPath, "cloudcertificate.tf")
	variablesPath := filepath.Join(tfWorkPath, "variables.tf")
	importPath := filepath.Join(tfWorkPath, "import.sh")
	if err := tools.CheckFiles(cloudCertificatePath, variablesPath, importPath); err != nil {
		return cli.Exit(color.RedString(err.Error()), 1)
	}

	name := c.Args().Get(0)

	params := createCloudCertificateParams{
		name:          name,
		edgercPath:    edgegrid.GetEdgercPath(c),
		configSection: edgegrid.GetEdgercSection(c),
		client:        cloudcertificates.Client(edgegrid.GetSession(c.Context)),
		templateProcessor: templates.FSTemplateProcessor{
			TemplatesFS: templateFiles,
			TemplateTargets: map[string]string{
				"cloudcertificate.tmpl": cloudCertificatePath,
				"variables.tmpl":        variablesPath,
				"import.tmpl":           importPath,
			},
			AdditionalFuncs: additionalFunctions,
		},
	}

	if err := createCloudCertificate(c.Context, params); err != nil {
		return cli.Exit(color.RedString("Error exporting cloud certificate: %s", err), 1)
	}
	return nil
}

func createCloudCertificate(ctx context.Context, params createCloudCertificateParams) (e error) {
	term := terminal.Get(ctx)
	defer func() {
		if e != nil {
			term.Spinner().Fail()
		}
	}()

	certs, err := listAllCertificates(ctx, params.client)
	if err != nil {
		return fmt.Errorf("%w: %w", ErrListingCloudCertificates, err)
	}

	var certID string
	for _, c := range certs.Certificates {
		if c.CertificateName == params.name {
			certID = c.CertificateID
			break
		}
	}
	if certID == "" {
		return fmt.Errorf("failed to fetch certificate: no certificate found with the name %q", params.name)
	}

	cert, err := params.client.GetCertificate(ctx, cloudcertificates.GetCertificateRequest{
		CertificateID: certID,
	})
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFetchingCloudCertificate, err)
	}
	term.Spinner().OK()

	term.Spinner().Start("Extracting data")
	tfData := populateTFData(params, cert.Certificate)
	term.Spinner().OK()

	term.Spinner().Start("Saving TF configurations")
	if err = params.templateProcessor.ProcessTemplates(tfData); err != nil {
		return fmt.Errorf("%w: %s", ErrSavingFiles, err)
	}
	term.Spinner().OK()

	term.Printf("Terraform configuration for cloud certificate name '%s' was saved successfully\n", params.name)

	return nil
}

func listAllCertificates(ctx context.Context, client cloudcertificates.CloudCertificates) (*cloudcertificates.ListCertificatesResponse, error) {
	var allCertificates cloudcertificates.ListCertificatesResponse
	request := cloudcertificates.ListCertificatesRequest{
		PageSize: maxPageSize,
		Page:     1,
	}

	for {
		certificatesResponse, err := client.ListCertificates(ctx, request)
		if err != nil {
			return nil, err
		}

		allCertificates.Certificates = append(allCertificates.Certificates, certificatesResponse.Certificates...)

		if certificatesResponse.Links.Next == nil {
			break
		}
		request.Page++
	}

	return &allCertificates, nil
}

func populateTFData(params createCloudCertificateParams, cert cloudcertificates.Certificate) TFData {
	var subject *cloudcertificates.Subject
	if cert.Subject != nil && !isEmptySubject(cert.Subject) {
		subject = cert.Subject
	}
	return TFData{
		EdgercPath: params.edgercPath,
		Section:    params.configSection,
		Certificate: TFCertificate{
			ID:                   cert.CertificateID,
			ContractID:           cert.ContractID,
			BaseName:             extractBaseName(cert.CertificateName),
			KeyType:              string(cert.KeyType),
			KeySize:              string(cert.KeySize),
			SecureNetwork:        cert.SecureNetwork,
			SANs:                 cert.SANs,
			CertificateStatus:    cert.CertificateStatus,
			Subject:              subject,
			SignedCertificatePEM: cert.SignedCertificatePEM,
			TrustChainPEM:        cert.TrustChainPEM,
			ResourceName:         sanitizeResourceName(cert.CertificateName),
		},
	}
}

// sanitizeResourceName replaces dots and spaces with underscores
// and ensures the resource name starts with a letter or underscore.
func sanitizeResourceName(name string) string {
	name = strings.ReplaceAll(name, ".", "_")
	name = strings.ReplaceAll(name, " ", "_")
	// If a first character is not a letter or underscore, prepend an underscore.
	if len(name) > 0 && !((name[0] >= 'A' && name[0] <= 'Z') || (name[0] >= 'a' && name[0] <= 'z') || name[0] == '_') {
		name = "_" + name
	}
	return name
}

func extractBaseName(name string) string {
	parts := strings.Split(name, renewedNameSuffix)
	if len(parts) != 2 || parts[0] == "" {
		// The name does not include the rotated suffix or starts with the suffix.
		return name
	}
	if _, err := time.Parse(renewedNameDateLayout, parts[1]); err == nil {
		// Valid date part, return the base name.
		return parts[0]
	}
	// Invalid date part, return the original name.
	return name
}

func isEmptySubject(subject *cloudcertificates.Subject) bool {
	return subject.CommonName == "" && subject.Organization == "" && subject.Country == "" &&
		subject.State == "" && subject.Locality == ""
}
