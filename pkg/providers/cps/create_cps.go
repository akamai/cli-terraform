// Package cps contains code for exporting Certificate Provisioning System (CPS) configuration
package cps

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v4/pkg/cps"
	"github.com/akamai/cli-terraform/pkg/edgegrid"
	"github.com/akamai/cli-terraform/pkg/templates"
	"github.com/akamai/cli-terraform/pkg/tools"
	"github.com/akamai/cli/pkg/terminal"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

type (
	// TFCPSData represents the data used in CPS templates
	TFCPSData struct {
		Enrollment          cps.Enrollment
		EnrollmentID        int
		ContractID          string
		Section             string
		CertificateECDSA    string
		TrustChainECDSA     string
		CertificateRSA      string
		TrustChainRSA       string
		NoUploadCertificate bool
	}
)

//go:embed templates/*
var templateFiles embed.FS

var (
	// ErrFetchingEnrollment is returned when fetching enrollment fails
	ErrFetchingEnrollment = errors.New("unable to fetch enrollment with given id")
	// ErrFetchingCertificateHistory is returned when fetching certificate history fails
	ErrFetchingCertificateHistory = errors.New("unable to fetch certificate history with given id")
	// ErrUnsupportedEnrollmentType is returned when user try to export OV or EV enrollments
	ErrUnsupportedEnrollmentType = errors.New("supporting export of dv and third-party enrollments but got")
)

// CmdCreateCPS is an entrypoint to create-cps command
func CmdCreateCPS(c *cli.Context) error {
	ctx := c.Context
	sess := edgegrid.GetSession(ctx)
	client := cps.Client(sess)

	// tfWorkPath is a target directory for generated terraform resources
	var tfWorkPath = "./"
	if c.IsSet("tfworkpath") {
		tfWorkPath = c.String("tfworkpath")
	}
	enrollmentPath := filepath.Join(tfWorkPath, "enrollment.tf")
	variablesPath := filepath.Join(tfWorkPath, "variables.tf")
	importPath := filepath.Join(tfWorkPath, "import.sh")

	err := tools.CheckFiles(enrollmentPath, variablesPath, importPath)
	if err != nil {
		return cli.Exit(color.RedString(err.Error()), 1)
	}

	templateToFile := map[string]string{
		"enrollment.tmpl": enrollmentPath,
		"variables.tmpl":  variablesPath,
		"imports.tmpl":    importPath,
	}

	processor := templates.FSTemplateProcessor{
		TemplatesFS:     templateFiles,
		TemplateTargets: templateToFile,
	}

	enrollmentID, err := strconv.Atoi(c.Args().Get(0))
	if err != nil {
		return cli.Exit(color.RedString(err.Error()), 1)
	}
	contractID := c.Args().Get(1)
	section := edgegrid.GetEdgercSection(c)
	if err = createCPS(ctx, contractID, enrollmentID, section, client, processor); err != nil {
		return cli.Exit(color.RedString(fmt.Sprintf("Error exporting enrollment HCL: %s", err)), 1)
	}
	return nil
}

func createCPS(ctx context.Context, contractID string, enrollmentID int,
	section string, client cps.CPS, templateProcessor templates.TemplateProcessor) error {
	term := terminal.Get(ctx)

	fmt.Println("Exporting CPS configuration")

	term.Spinner().Start(fmt.Sprintf("Fetching enrollment for the given id %d", enrollmentID))
	enrollment, err := client.GetEnrollment(ctx, cps.GetEnrollmentRequest{
		EnrollmentID: enrollmentID,
	})
	if err != nil || enrollment == nil {
		term.Spinner().Fail()
		return fmt.Errorf("%w: %s", ErrFetchingEnrollment, err)
	}

	if enrollment.ValidationType != "third-party" && enrollment.ValidationType != "dv" {
		term.Spinner().Fail()
		return fmt.Errorf("%w: %s", ErrUnsupportedEnrollmentType, enrollment.ValidationType)
	}

	term.Spinner().OK()

	tfData := TFCPSData{
		Enrollment:   *enrollment,
		EnrollmentID: enrollmentID,
		ContractID:   contractID,
		Section:      section,
	}

	if enrollment.ValidationType == "third-party" {
		term.Spinner().Start("Retrieving certificate history ")
		certHistory, err := client.GetChangeHistory(ctx, cps.GetChangeHistoryRequest{EnrollmentID: enrollmentID})
		if err != nil {
			term.Spinner().Fail()
			return fmt.Errorf("%w: %s", ErrFetchingCertificateHistory, err)
		}
		certificateECDSA, trustChainECDSA, certificateRSA, trustChainRSA := getCertificatesFromChangeHistory(certHistory)
		tfData.CertificateECDSA = strings.ReplaceAll(certificateECDSA, "\n", "\\n")
		tfData.TrustChainECDSA = strings.ReplaceAll(trustChainECDSA, "\n", "\\n")
		tfData.CertificateRSA = strings.ReplaceAll(certificateRSA, "\n", "\\n")
		tfData.TrustChainRSA = strings.ReplaceAll(trustChainRSA, "\n", "\\n")
		if certificateECDSA == "" && certificateRSA == "" {
			tfData.NoUploadCertificate = true
		}
		term.Spinner().OK()
	}

	term.Spinner().Start("Saving TF configurations ")
	if err := templateProcessor.ProcessTemplates(tfData); err != nil {
		term.Spinner().Fail()
		return err
	}
	term.Spinner().OK()
	fmt.Printf("Terraform configuration for enrollment '%d' was saved successfully\n", enrollmentID)

	return nil
}

// getCertificatesFromChangeHistory creates attributes for a resource form GetChangeHistoryResponse
func getCertificatesFromChangeHistory(changeHistory *cps.GetChangeHistoryResponse) (string, string, string, string) {
	var certificateECDSA, certificateRSA, trustChainECDSA, trustChainRSA string
	for _, change := range changeHistory.Changes {
		if change.PrimaryCertificate.KeyAlgorithm == "RSA" {
			certificateRSA = change.PrimaryCertificate.Certificate
			trustChainRSA = change.PrimaryCertificate.TrustChain
		} else {
			certificateECDSA = change.PrimaryCertificate.Certificate
			trustChainECDSA = change.PrimaryCertificate.TrustChain
		}
		if len(change.MultiStackedCertificates) != 0 {
			if change.MultiStackedCertificates[0].KeyAlgorithm == "RSA" {
				certificateRSA = change.MultiStackedCertificates[0].Certificate
				trustChainRSA = change.MultiStackedCertificates[0].TrustChain
			} else {
				certificateECDSA = change.MultiStackedCertificates[0].Certificate
				trustChainECDSA = change.MultiStackedCertificates[0].TrustChain
			}
		}
		if certificateECDSA != "" || certificateRSA != "" {
			break
		}
	}
	return certificateECDSA, trustChainECDSA, certificateRSA, trustChainRSA
}
