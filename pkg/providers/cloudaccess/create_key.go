// Package cloudaccess contains code for exporting CloudAccess key
package cloudaccess

import (
	"cmp"
	"context"
	"embed"
	"errors"
	"fmt"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/cloudaccess"
	"github.com/akamai/cli-terraform/pkg/edgegrid"
	"github.com/akamai/cli-terraform/pkg/templates"
	"github.com/akamai/cli-terraform/pkg/tools"
	"github.com/akamai/cli/v2/pkg/color"
	"github.com/akamai/cli/v2/pkg/terminal"
	"github.com/urfave/cli/v2"
)

type (
	// TFCloudAccessData represents the data used in CloudAccess
	TFCloudAccessData struct {
		Key     TFCloudAccessKey
		Section string
		Flag    bool
	}

	// TFCloudAccessKey represents the data used for export CloudAccess key
	TFCloudAccessKey struct {
		KeyResourceName      string
		AccessKeyName        string
		AuthenticationMethod string
		GroupID              int64
		ContractID           string
		AccessKeyUID         int64
		CredentialA          *Credential
		CredentialB          *Credential
		NetworkConfiguration *NetworkConfiguration
	}

	// Credential represents CLoudAccess credential
	Credential struct {
		CloudAccessKeyID string
	}

	// NetworkConfiguration represents CLoudAccess network configuration
	NetworkConfiguration struct {
		AdditionalCDN   *string
		SecurityNetwork string
	}
)

var (
	//go:embed templates/*
	templateFiles embed.FS

	additionalFunctions = tools.DecorateWithMultilineHandlingFunctions(map[string]any{})

	// ErrFetchingKey is returned when key could not be fetched
	ErrFetchingKey = errors.New("problem with fetching key")
	// ErrListingKeyVersions is returned when key versions could not be listed
	ErrListingKeyVersions = errors.New("problem with listing key versions")
	// ErrSavingFiles is returned when an issue with processing templates occurs
	ErrSavingFiles = errors.New("saving terraform project files")
	// ErrNoGroup is returned when key does not have group and contract assigned
	ErrNoGroup = errors.New("access key has no defined group or contract")
	// ErrNonUniqueCloudAccessKeyID is returned when key have the same `cloud_access_key_id` for both pairs of credentials
	ErrNonUniqueCloudAccessKeyID = errors.New("'cloud_access_key_id' should be unique for each pair of credentials")
)

// CmdCreateCloudAccess is an entrypoint to export-cloudaccess command
func CmdCreateCloudAccess(c *cli.Context) error {
	ctx := c.Context
	sess := edgegrid.GetSession(ctx)
	client := cloudaccess.Client(sess)

	// tfWorkPath is a target directory for generated terraform resources
	var tfWorkPath = "./"
	if c.IsSet("tfworkpath") {
		tfWorkPath = c.String("tfworkpath")
	}
	cloudAccessPath := filepath.Join(tfWorkPath, "cloudaccess.tf")
	variablesPath := filepath.Join(tfWorkPath, "variables.tf")
	importPath := filepath.Join(tfWorkPath, "import.sh")
	if err := tools.CheckFiles(cloudAccessPath, variablesPath, importPath); err != nil {
		return cli.Exit(color.RedString(err.Error()), 1)
	}

	templateToFile := map[string]string{
		"cloudaccess.tmpl": cloudAccessPath,
		"variables.tmpl":   variablesPath,
		"imports.tmpl":     importPath,
	}

	processor := templates.FSTemplateProcessor{
		TemplatesFS:     templateFiles,
		TemplateTargets: templateToFile,
		AdditionalFuncs: additionalFunctions,
	}

	keyUID, err := strconv.ParseInt(c.Args().Get(0), 10, 64)
	if err != nil {
		return cli.Exit(color.RedString(err.Error()), 1)
	}
	section := edgegrid.GetEdgercSection(c)

	// Get groupID and contractId from flags
	var groupID int64
	var contractID string

	if c.IsSet("group_id") {
		groupID, err = strconv.ParseInt(c.String("group_id"), 10, 64)
		if err != nil {
			return cli.Exit(color.RedString("Invalid group_id: %s", err.Error()), 1)
		}
		if groupID <= 0 {
			// Check if group ID is less than or equal to 0
			return errors.New("Invalid group ID: group ID must be greater than 0")
		}
		if !c.IsSet("contract_id") {
			return cli.Exit(color.RedString("contract_id is mandatory when group_id is provided"), 1)
		}
		contractID = c.String("contract_id")
	} else if c.IsSet("contract_id") {
		return cli.Exit(color.RedString("contract_id cannot be set without group_id"), 1)
	}

	if err = createCloudAccess(ctx, keyUID, groupID, contractID, section, client, processor); err != nil {
		return cli.Exit(color.RedString("Error exporting cloudaccess: %s", err), 1)
	}
	return nil
}

func createCloudAccess(ctx context.Context, accessKeyUID int64, groupID int64, contractID string, section string, client cloudaccess.CloudAccess, templateProcessor templates.TemplateProcessor) error {
	term := terminal.Get(ctx)
	term.Spinner().Start("Fetching cloudaccess key " + strconv.Itoa(int(accessKeyUID)))
	key, err := client.GetAccessKey(ctx, cloudaccess.AccessKeyRequest{
		AccessKeyUID: accessKeyUID,
	})
	if err != nil {
		term.Spinner().Fail()
		return fmt.Errorf("%w: %s", ErrFetchingKey, err)
	}

	if len(key.Groups) == 0 {
		return ErrNoGroup
	}

	versions, err := client.ListAccessKeyVersions(ctx, cloudaccess.ListAccessKeyVersionsRequest{
		AccessKeyUID: accessKeyUID,
	})
	if err != nil {
		term.Spinner().Fail()
		return fmt.Errorf("%w: %s", ErrListingKeyVersions, err)
	}
	if len(versions.AccessKeyVersions) > 1 {
		if *versions.AccessKeyVersions[0].CloudAccessKeyID == *versions.AccessKeyVersions[1].CloudAccessKeyID {
			term.Spinner().Fail()
			return fmt.Errorf("%w", ErrNonUniqueCloudAccessKeyID)
		}
	}
	tfCloudAccessData, err := populateCloudAccessData(section, key, groupID, contractID, versions.AccessKeyVersions)
	if err != nil {
		term.Spinner().Fail()
		return fmt.Errorf("error populating cloud access data: %w", err)
	}

	term.Spinner().Start("Saving TF configurations ")
	if err = templateProcessor.ProcessTemplates(tfCloudAccessData); err != nil {
		term.Spinner().Fail()
		return fmt.Errorf("%w: %s", ErrSavingFiles, err)
	}

	term.Spinner().OK()
	term.Printf("Terraform configuration for cloudaccess key '%s' was saved successfully\n", tfCloudAccessData.Key.AccessKeyName)

	return nil
}

func populateCloudAccessData(section string, key *cloudaccess.GetAccessKeyResponse, providedGroupID int64, providedContractID string, versions []cloudaccess.AccessKeyVersion) (TFCloudAccessData, error) {
	var netConf *NetworkConfiguration
	if key.NetworkConfiguration != nil {
		netConf = &NetworkConfiguration{
			SecurityNetwork: string(key.NetworkConfiguration.SecurityNetwork),
		}
		if key.NetworkConfiguration.AdditionalCDN != nil {
			netConf.AdditionalCDN = tools.StringPtr(string(*key.NetworkConfiguration.AdditionalCDN))
		}
	}

	var contractID string
	var groupID int64
	var flag bool
	// Check if both providedGroupId and providedContractId are supplied
	if providedGroupID != 0 && providedContractID != "" {
		// Validate the group and contract combination
		for _, group := range key.Groups {
			if group.GroupID == providedGroupID {
				for _, contract := range group.ContractIDs {
					if contract == providedContractID {
						groupID = providedGroupID
						contractID = providedContractID
						flag = true
					}
				}
			}
		}
		// If no match found, return error
		if !flag {
			return TFCloudAccessData{}, fmt.Errorf("invalid combination of groupId (%d) and contractId (%s) for this access key", providedGroupID, providedContractID)
		}
	} else if len(key.Groups) > 0 {
		groups := key.Groups
		slices.SortFunc(groups, func(a, b cloudaccess.Group) int {
			return cmp.Compare(a.GroupID, b.GroupID)
		})
		groupID = groups[len(groups)-1].GroupID
		if len(groups[len(groups)-1].ContractIDs) > 0 {
			contractID = groups[len(groups)-1].ContractIDs[0]
		}
	}

	tfCloudAccessData := TFCloudAccessData{
		Section: section,
		Key: TFCloudAccessKey{
			KeyResourceName:      strings.ReplaceAll(key.AccessKeyName, "-", "_"),
			AccessKeyName:        key.AccessKeyName,
			AuthenticationMethod: key.AuthenticationMethod,
			GroupID:              groupID,
			ContractID:           contractID,
			AccessKeyUID:         key.AccessKeyUID,
			NetworkConfiguration: netConf,
		},
		Flag: flag,
	}

	versionNum := len(versions)
	switch versionNum {
	case 0:
		return tfCloudAccessData, nil
	case 1:
		tfCloudAccessData.Key.CredentialA = &Credential{
			CloudAccessKeyID: *versions[0].CloudAccessKeyID,
		}
	default:
		slices.SortFunc(versions, func(a, b cloudaccess.AccessKeyVersion) int {
			return cmp.Compare(a.Version, b.Version)
		})
		// first version from the response from API is assigned to `credentials_b`, second version to `credentials_a`
		tfCloudAccessData.Key.CredentialA = &Credential{
			CloudAccessKeyID: *versions[0].CloudAccessKeyID,
		}
		tfCloudAccessData.Key.CredentialB = &Credential{
			CloudAccessKeyID: *versions[1].CloudAccessKeyID,
		}
	}

	return tfCloudAccessData, nil
}
