// Package clientlists contains code for exporting Client Lists
package clientlists

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/clientlists"
	"github.com/akamai/cli-terraform/v2/pkg/edgegrid"
	"github.com/akamai/cli-terraform/v2/pkg/templates"
	"github.com/akamai/cli-terraform/v2/pkg/tools"
	"github.com/akamai/cli/v2/pkg/color"
	"github.com/akamai/cli/v2/pkg/terminal"
	"github.com/urfave/cli/v2"
)

//go:embed templates/*
var templateFiles embed.FS

var (
	// ErrClientListNotFound list not found
	ErrClientListNotFound = errors.New("client list not found")
	// ErrSavingFiles saving terraform project files
	ErrSavingFiles = errors.New("saving terraform project files")
	// ErrSavingListItems saving client list items file
	ErrSavingListItems = errors.New("saving client list items file")
	// ErrActivationDetails retrieving activation details
	ErrActivationDetails = errors.New("retrieving activation details")
)

// TFData holds template data
type TFData struct {
	ClientList TFListData
	Section    string
	EdgercPath string
}

// TFActivationData holds template data for activation
type TFActivationData struct {
	HasActivation          bool
	Comments               string
	SiebelTicketID         string
	NotificationRecipients []string
}

// TFListData holds template data for list
type TFListData struct {
	ListID               string
	Name                 string
	Notes                string
	Type                 string
	ContractID           string
	GroupID              int64
	ItemsCount           int64
	Tags                 []string
	StagingActivation    TFActivationData
	ProductionActivation TFActivationData
}

// CmdCreateClientList is an entrypoint to create-clientlist command
func CmdCreateClientList(c *cli.Context) error {
	ctx := c.Context
	sess := edgegrid.GetSession(c.Context)
	client := clientlists.Client(sess)

	// tfWorkPath is a target directory for generated terraform resources
	var tfWorkPath = "./"
	if c.IsSet("tfworkpath") {
		tfWorkPath = c.String("tfworkpath")
	}

	listID := c.Args().First()

	clientListPath := filepath.Join(tfWorkPath, "client-list.tf")
	variablesPath := filepath.Join(tfWorkPath, "variables.tf")
	importPath := filepath.Join(tfWorkPath, "imports.sh")
	jsonPath := filepath.Join(tfWorkPath, fmt.Sprintf("%s.json", listID))

	if err := tools.CheckFiles(clientListPath, variablesPath, importPath, jsonPath); err != nil {
		return cli.Exit(color.RedString(err.Error()), 1)
	}

	templateToFile := map[string]string{
		"client-list.tmpl": clientListPath,
		"variables.tmpl":   variablesPath,
		"imports.tmpl":     importPath,
	}

	processor := templates.FSTemplateProcessor{
		TemplatesFS:     templateFiles,
		TemplateTargets: templateToFile,
	}

	section := edgegrid.GetEdgercSection(c)
	edgercPath := edgegrid.GetEdgercPath(c)

	if err := createClientList(ctx, listID, edgercPath, section, tfWorkPath, client, processor); err != nil {
		return cli.Exit(color.RedString("Error exporting client list: %s", err), 1)
	}

	return nil
}

func createClientList(ctx context.Context, listID, edgercPath, section, tfWorkPath string, client clientlists.ClientLists, processor templates.TemplateProcessor) error {
	term := terminal.Get(ctx)

	term.Spinner().Start("Fetching client list " + listID)
	clientList, err := client.GetClientList(ctx, clientlists.GetClientListRequest{
		ListID:       listID,
		IncludeItems: true,
	})
	if err != nil {
		term.Spinner().Fail()
		return fmt.Errorf("%w: %s", ErrClientListNotFound, err)
	}

	stagingActivation, errStaging := getActivationDataByNetwork(ctx, client, clientList, clientlists.Staging)
	productionActivation, errProd := getActivationDataByNetwork(ctx, client, clientList, clientlists.Production)
	if errStaging != nil || errProd != nil {
		term.Spinner().Fail()
		return fmt.Errorf("%w: %s", ErrActivationDetails, err)
	}
	term.Spinner().OK()

	tfData := TFData{
		ClientList: TFListData{
			ListID:               clientList.ListID,
			Name:                 clientList.Name,
			Type:                 string(clientList.Type),
			Tags:                 clientList.Tags,
			Notes:                clientList.Notes,
			ContractID:           clientList.ContractID,
			GroupID:              clientList.GroupID,
			ItemsCount:           clientList.ItemsCount,
			StagingActivation:    stagingActivation,
			ProductionActivation: productionActivation,
		},
		Section:    section,
		EdgercPath: edgercPath,
	}

	term.Spinner().Start("Saving TF configurations ")
	if err = processor.ProcessTemplates(tfData); err != nil {
		term.Spinner().Fail()
		return fmt.Errorf("%w: %s", ErrSavingFiles, err)
	}

	if err := saveListItemsJSON(clientList, tfWorkPath); err != nil {
		term.Spinner().Fail()
		return fmt.Errorf("%w: %s", ErrSavingListItems, err)
	}

	term.Spinner().OK()
	term.Printf("Terraform configuration for client list '%s' was saved successfully\n", clientList.Name)

	return nil
}

func saveListItemsJSON(clientList *clientlists.GetClientListResponse, tfWorkPath string) error {
	listItems := []clientlists.ListItemPayload{}

	for _, v := range clientList.Items {
		item := clientlists.ListItemPayload{
			Value: v.Value, // Value is always included
		}

		if v.Description != "" {
			item.Description = v.Description
		}

		if len(v.Tags) > 0 {
			tags := make([]string, len(v.Tags))
			copy(tags, v.Tags)
			item.Tags = tags
		}

		if v.ExpirationDate != "" {
			item.ExpirationDate = v.ExpirationDate
		}

		listItems = append(listItems, item)
	}

	jsonBody, err := json.MarshalIndent(listItems, "", "  ")
	if err != nil {
		return fmt.Errorf("can't marshall list items: %s", err)
	}

	path := filepath.Join(tfWorkPath, fmt.Sprintf("%s.json", clientList.ListID))

	if err = os.WriteFile(path, jsonBody, 0644); err != nil {
		return fmt.Errorf("can't write list items json: %s", err)
	}
	return nil
}

func getActivationDataByNetwork(ctx context.Context, client clientlists.ClientLists, cl *clientlists.GetClientListResponse, network clientlists.ActivationNetwork) (TFActivationData, error) {
	activationData := TFActivationData{HasActivation: false}
	var activationStatus clientlists.ActivationStatus

	if network == clientlists.Staging {
		activationStatus = clientlists.ActivationStatus(cl.StagingActivationStatus)
	} else {
		activationStatus = clientlists.ActivationStatus(cl.ProductionActivationStatus)
	}

	if activationStatus != clientlists.Inactive {
		res, err := client.GetActivationStatus(ctx, clientlists.GetActivationStatusRequest{
			ListID:  cl.ListID,
			Network: network,
		})
		if err != nil {
			return activationData, err
		}

		activationData = TFActivationData{
			HasActivation:          true,
			Comments:               res.Comments,
			SiebelTicketID:         res.SiebelTicketID,
			NotificationRecipients: res.NotificationRecipients,
		}
	}

	return activationData, nil
}
