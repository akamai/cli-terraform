// Package cloudcertificates contains code for exporting Cloud Certificates lineages.
package cloudcertificates

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	edgegridcloudcertificates "github.com/akamai/AkamaiOPEN-edgegrid-golang/v15/pkg/cloudcertificates"
	"github.com/akamai/cli-terraform/v3/pkg/edgegrid"
	"github.com/akamai/cli-terraform/v3/pkg/templates"
	"github.com/akamai/cli-terraform/v3/pkg/tools"
	"github.com/akamai/cli/v2/pkg/color"
	"github.com/akamai/cli/v2/pkg/terminal"
	"github.com/urfave/cli/v3"
)

type (
	tfData struct {
		ResourceName        string
		Lineage             tfLineage
		Upload              tfUpload
		UploadGenerationID  int64
		Activation          tfActivation
		AcknowledgeWarnings bool
		EdgercPath          string
		Section             string
	}

	tfLineage struct {
		LineageID     int64
		ContractID    string
		GroupID       int64
		GeoClass      string
		LineageName   string
		LineageType   string
		SecureNetwork string
		KeySpecs      []keySpec
		SANs          []string
		Subject       *subject
	}

	keySpec struct {
		KeyType string
		KeySize string
	}

	subject struct {
		CommonName         string
		Organization       string
		OrganizationalUnit string
		Country            string
		State              string
		Locality           string
	}

	tfUpload struct {
		LineageID           int64
		GenerationID        int64
		AcknowledgeWarnings bool
		IsCommented         bool
		Algorithms          []uploadAlgorithm
	}

	tfActivation struct {
		StagingActivationID          int64
		ProductionActivationID       int64
		StagingActivationIsLatest    bool
		ProductionActivationIsLatest bool
		ImportID                     string
		HasAnyActivation             bool
	}

	uploadAlgorithm struct {
		KeyType              string
		SignedCertificatePEM *string
		TrustChainPEM        *string
		IsCommented          bool
	}

	createCloudCertificateParams struct {
		lineageName         string
		acknowledgeWarnings bool
		edgercPath          string
		configSection       string
		client              edgegridcloudcertificates.CloudCertificates
		templateProcessor   templates.TemplateProcessor
	}
)

var (
	//go:embed templates/*
	templateFiles embed.FS

	additionalFunctions = tools.DecorateWithMultilineHandlingFunctions(map[string]any{})

	// ErrFetchingLineage is returned when retrieving a certificate lineage fails.
	ErrFetchingLineage = errors.New("error fetching certificate lineage")
	// ErrSavingFiles is returned when processing templates fails.
	ErrSavingFiles = errors.New("error saving terraform project files")
	// ErrFetchingActivations is returned when retrieving lineage activations fails.
	ErrFetchingActivations = errors.New("error fetching certificate lineage activations")
)

// CmdCreateCloudCertificate is an entrypoint to the export-cloudcertificate command.
func CmdCreateCloudCertificate(ctx context.Context, cmd *cli.Command) error {
	tfWorkPath := "./"
	if cmd.IsSet("tfworkpath") {
		tfWorkPath = cmd.String("tfworkpath")
	}

	cloudCertificatePath := filepath.Join(tfWorkPath, "cloudcertificate.tf")
	variablesPath := filepath.Join(tfWorkPath, "variables.tf")
	importPath := filepath.Join(tfWorkPath, "import.sh")
	if err := tools.CheckFiles(cloudCertificatePath, variablesPath, importPath); err != nil {
		return cli.Exit(color.RedString("%s", err.Error()), 1)
	}

	params := createCloudCertificateParams{
		lineageName:   cmd.Args().Get(0),
		edgercPath:    edgegrid.GetEdgercPath(cmd),
		configSection: edgegrid.GetEdgercSection(cmd),
		client:        edgegridcloudcertificates.Client(edgegrid.GetSession(ctx)),
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
	params.acknowledgeWarnings = cmd.Bool("acknowledge-warnings")

	if err := createCloudCertificate(ctx, params); err != nil {
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

	term.Spinner().Start("Fetching certificate lineage '%s'", params.lineageName)
	lineageID, err := lineageIDByName(ctx, params.client, params.lineageName)
	if err != nil {
		return err
	}
	lineage, err := params.client.GetLineage(ctx, edgegridcloudcertificates.GetLineageRequest{
		LineageID: lineageID,
		ExpandGenerations: []edgegridcloudcertificates.ExpandGenerations{
			edgegridcloudcertificates.ExpandGenerationsHead,
			edgegridcloudcertificates.ExpandGenerationsCurrentProduction,
			edgegridcloudcertificates.ExpandGenerationsCurrentStaging,
		},
	})
	if err != nil {
		return fmt.Errorf("%w: %w", ErrFetchingLineage, err)
	}
	if lineage == nil {
		return fmt.Errorf("%w: empty response", ErrFetchingLineage)
	}
	term.Spinner().OK()

	lineageData := edgegridcloudcertificates.Lineage(*lineage)
	var stagingActivation, productionActivation *edgegridcloudcertificates.GetActivationStatusResponse
	if lineageData.CurrentStaging != nil || lineageData.CurrentProduction != nil {
		term.Spinner().Start("Fetching certificate activations for lineage '%s'", params.lineageName)
		var fetchErr error
		stagingActivation, productionActivation, fetchErr = fetchActivationData(ctx, params.client, lineageData)
		if fetchErr != nil {
			return fetchErr
		}
		term.Spinner().OK()
	}

	term.Spinner().Start("Extracting data for lineage '%s'", params.lineageName)
	data, warnings := populateTFData(params, lineageData, stagingActivation, productionActivation)
	term.Spinner().OK()
	for _, warning := range warnings {
		term.Printf("%s", warning)
	}

	term.Spinner().Start("Saving TF configurations for lineage '%s'", params.lineageName)
	if err := params.templateProcessor.ProcessTemplates(data); err != nil {
		return fmt.Errorf("%w: %w", ErrSavingFiles, err)
	}
	term.Spinner().OK()

	term.Printf("Terraform configuration for certificate lineage '%s' was saved successfully\n", params.lineageName)

	return nil
}

func lineageIDByName(ctx context.Context, client edgegridcloudcertificates.CloudCertificates, lineageName string) (int64, error) {
	cursor := ""
	for {
		lineages, err := client.ListLineages(ctx, edgegridcloudcertificates.ListLineagesRequest{
			LineageName: lineageName,
			After:       cursor,
		})
		if err != nil {
			return 0, fmt.Errorf("%w: %w", ErrFetchingLineage, err)
		}
		if lineages == nil {
			return 0, fmt.Errorf("%w: empty response", ErrFetchingLineage)
		}
		for _, lineage := range lineages.Lineages {
			if lineage.LineageName == lineageName {
				return lineage.LineageID, nil
			}
		}
		if lineages.NextCursor == nil {
			return 0, fmt.Errorf("%w: lineage with name %q not found", ErrFetchingLineage, lineageName)
		}
		cursor = *lineages.NextCursor
	}
}

func populateTFData(params createCloudCertificateParams, lineage edgegridcloudcertificates.Lineage, stagingActivation, productionActivation *edgegridcloudcertificates.GetActivationStatusResponse) (tfData, []string) {
	keySpecs := make([]keySpec, 0, len(lineage.KeySpecs))
	for _, spec := range lineage.KeySpecs {
		keySpecs = append(keySpecs, keySpec{
			KeyType: spec.KeyType,
			KeySize: spec.KeySize,
		})
	}
	sort.Slice(keySpecs, func(i, j int) bool {
		return keySpecs[i].KeyType < keySpecs[j].KeyType
	})

	sans := append([]string(nil), lineage.SANs...)
	sort.Strings(sans)

	upload, warnings := buildLineageUpload(params, lineage)
	activation, activationWarnings := buildLineageActivation(lineage, upload.GenerationID, stagingActivation, productionActivation)
	warnings = append(warnings, activationWarnings...)
	return tfData{
		ResourceName: tools.SanitizeResourceName(lineage.LineageName),
		EdgercPath:   params.edgercPath,
		Section:      params.configSection,
		Lineage: tfLineage{
			LineageID:     lineage.LineageID,
			ContractID:    lineage.ContractID,
			GroupID:       lineage.GroupID,
			GeoClass:      lineage.GeoClass,
			LineageName:   lineage.LineageName,
			LineageType:   lineage.LineageType,
			SecureNetwork: lineage.SecureNetwork,
			KeySpecs:      keySpecs,
			SANs:          sans,
			Subject:       lineageSubject(lineage.Subject),
		},
		Upload:              upload,
		UploadGenerationID:  upload.GenerationID,
		Activation:          activation,
		AcknowledgeWarnings: params.acknowledgeWarnings,
	}, warnings
}

func fetchActivationData(ctx context.Context, client edgegridcloudcertificates.CloudCertificates, lineage edgegridcloudcertificates.Lineage) (*edgegridcloudcertificates.GetActivationStatusResponse, *edgegridcloudcertificates.GetActivationStatusResponse, error) {
	hasStaging := lineage.CurrentStaging != nil
	hasProduction := lineage.CurrentProduction != nil

	var stagingActivation, productionActivation *edgegridcloudcertificates.GetActivationStatusResponse
	cursor := ""
	for {
		response, err := client.ListActivations(ctx, edgegridcloudcertificates.ListActivationsRequest{
			LineageID: lineage.LineageID,
			Cursor:    cursor,
		})
		if err != nil {
			return nil, nil, fmt.Errorf("%w: %w", ErrFetchingActivations, err)
		}
		if response == nil {
			return nil, nil, fmt.Errorf("%w: empty response", ErrFetchingActivations)
		}

		for i, current := range response.Items {
			switch {
			case hasStaging && stagingActivation == nil &&
				current.TargetEnvironment == string(edgegridcloudcertificates.TargetNetworkStaging) &&
				current.GenerationID == lineage.CurrentStaging.ID():

				stagingActivation = &response.Items[i]
			case hasProduction && productionActivation == nil &&
				current.TargetEnvironment == string(edgegridcloudcertificates.TargetNetworkProduction) &&
				current.GenerationID == lineage.CurrentProduction.ID():

				productionActivation = &response.Items[i]
			}
		}

		foundAllRequiredActivations := (!hasStaging || stagingActivation != nil) && (!hasProduction || productionActivation != nil)
		if foundAllRequiredActivations {
			return stagingActivation, productionActivation, nil
		}
		if response.NextCursor == nil {
			var missingNetworks []string
			if hasStaging && stagingActivation == nil {
				missingNetworks = append(missingNetworks, string(edgegridcloudcertificates.TargetNetworkStaging))
			}
			if hasProduction && productionActivation == nil {
				missingNetworks = append(missingNetworks, string(edgegridcloudcertificates.TargetNetworkProduction))
			}
			return nil, nil, fmt.Errorf(
				"%w: no activation found for a current network generation on %s",
				ErrFetchingActivations,
				strings.Join(missingNetworks, " and "),
			)
		}
		cursor = *response.NextCursor
	}
}

func buildLineageActivation(lineage edgegridcloudcertificates.Lineage, uploadGenerationID int64, stagingActivation, productionActivation *edgegridcloudcertificates.GetActivationStatusResponse) (tfActivation, []string) {
	activation := tfActivation{}
	var warnings []string

	if stagingActivation != nil {
		if isActivationSuccessful(*stagingActivation) {
			activation.StagingActivationID = stagingActivation.ActivationID
		} else {
			warnings = append(warnings, incompleteActivationWarning(
				lineage.LineageName, *stagingActivation, edgegridcloudcertificates.TargetNetworkStaging,
			))
		}
	}
	if productionActivation != nil {
		if isActivationSuccessful(*productionActivation) {
			activation.ProductionActivationID = productionActivation.ActivationID
		} else {
			warnings = append(warnings, incompleteActivationWarning(
				lineage.LineageName, *productionActivation, edgegridcloudcertificates.TargetNetworkProduction,
			))
		}
	}

	activation.HasAnyActivation = activation.StagingActivationID != 0 || activation.ProductionActivationID != 0
	activation.StagingActivationIsLatest = uploadGenerationID != 0 &&
		lineage.CurrentStaging != nil && lineage.CurrentStaging.ID() == uploadGenerationID
	activation.ProductionActivationIsLatest = uploadGenerationID != 0 &&
		lineage.CurrentProduction != nil && lineage.CurrentProduction.ID() == uploadGenerationID
	activation.ImportID = calculateActivationImportID(lineage, &activation)
	return activation, warnings
}

func calculateActivationImportID(lineage edgegridcloudcertificates.Lineage, activation *tfActivation) string {
	ids := []string{strconv.FormatInt(lineage.LineageID, 10)}
	if activation.StagingActivationID != 0 {
		ids = append(ids, strconv.FormatInt(activation.StagingActivationID, 10))
	}
	if activation.ProductionActivationID != 0 {
		ids = append(ids, strconv.FormatInt(activation.ProductionActivationID, 10))
	}
	return strings.Join(ids, ",")
}

func isActivationSuccessful(current edgegridcloudcertificates.GetActivationStatusResponse) bool {
	return current.ActivationStatus == "COMPLETE" || current.ActivationStatus == "PARTIAL_SUCCESS"
}

func incompleteActivationWarning(lineageName string, activation edgegridcloudcertificates.GetActivationStatusResponse, network edgegridcloudcertificates.TargetNetwork) string {
	return fmt.Sprintf(
		"[WARN] Activation %d for lineage %q on the %s network and generation %d has status %s; the network activation will be commented out.\n",
		activation.ActivationID,
		lineageName,
		network,
		activation.GenerationID,
		activation.ActivationStatus,
	)
}

func lineageSubject(source edgegridcloudcertificates.Subject) *subject {
	if source == (edgegridcloudcertificates.Subject{}) {
		return nil
	}

	return &subject{
		CommonName:         source.CommonName,
		Organization:       source.Organization,
		OrganizationalUnit: source.OrganizationalUnit,
		Country:            source.Country,
		State:              source.State,
		Locality:           source.Locality,
	}
}

func buildLineageUpload(params createCloudCertificateParams, lineage edgegridcloudcertificates.Lineage) (tfUpload, []string) {
	generation := calculateUploadGeneration(lineage)

	var algorithms []uploadAlgorithm
	var warnings []string
	if generation != nil {
		algorithms = make([]uploadAlgorithm, 0, len(generation.Common().Algorithms))
		for _, algorithm := range generation.Common().Algorithms {
			isReadyForUse := algorithm.CertificateStatus == edgegridcloudcertificates.CertificateStatusReadyForUse
			wasCorrectlyUploaded := isReadyForUse && algorithm.SignedCertificatePEM != nil
			if isReadyForUse && algorithm.SignedCertificatePEM == nil {
				warnings = append(warnings, fmt.Sprintf(
					"[WARN] Algorithm %q for lineage %q is READY_FOR_USE but has no signed certificate; it will be commented out.\n",
					algorithm.KeyType,
					lineage.LineageName,
				))
			}
			algorithms = append(algorithms, uploadAlgorithm{
				KeyType:              algorithm.KeyType,
				SignedCertificatePEM: algorithm.SignedCertificatePEM,
				TrustChainPEM:        algorithm.TrustChainPEM,
				IsCommented:          !wasCorrectlyUploaded,
			})
		}
	} else {
		algorithms = make([]uploadAlgorithm, 0, len(lineage.KeySpecs))
		for _, keySpec := range lineage.KeySpecs {
			algorithms = append(algorithms, uploadAlgorithm{
				KeyType:     keySpec.KeyType,
				IsCommented: true,
			})
		}
	}
	sort.Slice(algorithms, func(i, j int) bool {
		return algorithms[i].KeyType < algorithms[j].KeyType
	})
	isCommented := true
	for _, algorithm := range algorithms {
		if !algorithm.IsCommented {
			isCommented = false
			break
		}
	}

	upload := tfUpload{
		LineageID:           lineage.LineageID,
		AcknowledgeWarnings: params.acknowledgeWarnings,
		IsCommented:         isCommented,
		Algorithms:          algorithms,
	}
	if generation != nil {
		upload.GenerationID = generation.ID()
	}
	return upload, warnings
}

func calculateUploadGeneration(lineage edgegridcloudcertificates.Lineage) edgegridcloudcertificates.GenerationPointer {
	switch {
	case lineage.Head != nil:
		return lineage.Head
	case lineage.CurrentProduction != nil:
		return lineage.CurrentProduction
	default:
		return nil
	}
}
