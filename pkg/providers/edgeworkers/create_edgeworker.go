package edgeworkers

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/edgeworkers"
	"github.com/akamai/cli-terraform/pkg/edgegrid"
	"github.com/akamai/cli-terraform/pkg/templates"
	"github.com/akamai/cli-terraform/pkg/tools"
	"github.com/akamai/cli/pkg/terminal"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

type (
	// TFEdgeWorkerData represents the data used in EdgeWorker templates
	TFEdgeWorkerData struct {
		EdgeWorkerID   int
		Name           string
		GroupID        int64
		ResourceTierID int
		LocalBundle    string
		Section        string
		Note           string
	}
)

var (
	// ErrFetchingEdgeWorker is returned when fetching edgeworker fails
	ErrFetchingEdgeWorker = errors.New("unable to fetch edgeworker with given edgeworker_id")
)

const (
	activationStatusComplete = "COMPLETE"
)

// CmdCreateEdgeWorker is an entrypoint to create-edgeworker command
func CmdCreateEdgeWorker(c *cli.Context) error {
	ctx := c.Context
	sess := edgegrid.GetSession(c.Context)
	client := edgeworkers.Client(sess)

	// tfWorkPath is a target directory for generated terraform resources
	var tfWorkPath = "./"
	if c.IsSet("tfworkpath") {
		tfWorkPath = c.String("tfworkpath")
	}

	edgeWorkerPath := filepath.Join(tfWorkPath, "edgeworker.tf")
	variablesPath := filepath.Join(tfWorkPath, "variables.tf")
	importPath := filepath.Join(tfWorkPath, "import.sh")

	bundleDir := tfWorkPath
	if c.IsSet("bundlepath") {
		bundleDir = c.String("bundlepath")
	}
	bundleDir = filepath.FromSlash(bundleDir)
	if stat, err := os.Stat(bundleDir); err != nil || !stat.IsDir() {
		return cli.Exit(color.RedString("Bundle path is not accessible"), 1)
	}

	err := tools.CheckFiles(edgeWorkerPath, variablesPath, importPath)
	if err != nil {
		return cli.Exit(color.RedString(err.Error()), 1)
	}
	templateToFile := map[string]string{
		"edgeworker.tmpl":           edgeWorkerPath,
		"edgeworker-variables.tmpl": variablesPath,
		"edgeworker-imports.tmpl":   importPath,
	}

	processor := templates.FSTemplateProcessor{
		TemplatesFS:     templateFiles,
		TemplateTargets: templateToFile,
		AdditionalFuncs: template.FuncMap{
			"ToLower": func(network edgeworkers.ActivationNetwork) string {
				return strings.ToLower(string(network))
			},
			"Escape": tools.Escape,
		},
	}

	edgeWorkerID, err := strconv.Atoi(c.Args().First())
	if err != nil {
		return cli.Exit(color.RedString("edgeworker_id is not a valid integer"), 1)
	}
	section := edgegrid.GetEdgercSection(c)

	if err = createEdgeWorker(ctx, edgeWorkerID, bundleDir, section, client, processor); err != nil {
		return cli.Exit(color.RedString(fmt.Sprintf("Error exporting edgeworker HCL: %s", err)), 1)
	}
	return nil
}

func createEdgeWorker(ctx context.Context, edgeWorkerID int, bundleDir, section string, client edgeworkers.Edgeworkers, templateProcessor templates.TemplateProcessor) error {
	term := terminal.Get(ctx)
	fmt.Println("Configuring EdgeWorker")
	term.Spinner().Start(fmt.Sprintf("Fetching EdgeWorker %d", edgeWorkerID), "")

	edgeWorker, err := client.GetEdgeWorkerID(ctx, edgeworkers.GetEdgeWorkerIDRequest{
		EdgeWorkerID: edgeWorkerID,
	})
	if err != nil {
		term.Spinner().Fail()
		return fmt.Errorf("%w: %s", ErrFetchingEdgeWorker, err)
	}

	versions, err := client.ListEdgeWorkerVersions(ctx, edgeworkers.ListEdgeWorkerVersionsRequest{
		EdgeWorkerID: edgeWorkerID,
	})
	if err != nil {
		return fmt.Errorf("%w: %s", ErrFetchingEdgeWorker, err)
	}

	var activation *edgeworkers.Activation

	localBundle := ""
	if len(versions.EdgeWorkerVersions) != 0 {
		var version string
		var createdTime time.Time
		for _, v := range versions.EdgeWorkerVersions {
			parsedCreatedTime, err := time.Parse(time.RFC3339, v.CreatedTime)
			if err != nil {
				return err
			}
			if parsedCreatedTime.After(createdTime) {
				createdTime = parsedCreatedTime
				version = v.Version
			}
		}

		activation, err = getCurrentActivation(ctx, client, edgeWorkerID, "STAGING", version)
		if err != nil {
			return err
		}
		localBundle, err = getEdgeWorkerBundle(ctx, edgeWorkerID, version, bundleDir, client)
		if err != nil {
			term.Spinner().Fail()
			return fmt.Errorf("%w: %s", ErrFetchingEdgeWorker, err)
		}
	}

	tfEdgeWorkerData := TFEdgeWorkerData{
		EdgeWorkerID:   edgeWorkerID,
		Name:           edgeWorker.Name,
		GroupID:        edgeWorker.GroupID,
		ResourceTierID: edgeWorker.ResourceTierID,
		LocalBundle:    localBundle,
		Section:        section,
	}

	if activation != nil {
		tfEdgeWorkerData.Note = activation.Note
	}

	term.Spinner().OK()
	term.Spinner().Start("Saving TF configurations ", "")
	if err := templateProcessor.ProcessTemplates(tfEdgeWorkerData); err != nil {
		term.Spinner().Fail()
		return fmt.Errorf("%w: %s", templates.ErrSavingFiles, err)
	}
	term.Spinner().OK()
	fmt.Printf("Terraform configuration for edgeworker '%s' with edgeworker_id '%d' was saved successfully\n", edgeWorker.Name, edgeWorkerID)

	return nil
}

// getEdgeWorkerBundle fetches the bundle content of the latest version given edgeWorkerID and returns the path to it
func getEdgeWorkerBundle(ctx context.Context, edgeWorkerID int, version, bundlePath string, client edgeworkers.Edgeworkers) (string, error) {
	bundleContent, err := client.GetEdgeWorkerVersionContent(ctx, edgeworkers.GetEdgeWorkerVersionContentRequest{
		EdgeWorkerID: edgeWorkerID,
		Version:      version,
	})
	if err != nil {
		return "", err
	}

	localBundle := filepath.Join(bundlePath, version+".tgz")
	f, err := os.OpenFile(localBundle, os.O_CREATE|os.O_WRONLY, 0755)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(f, bundleContent); err != nil {
		return "", err
	}
	if err := f.Close(); err != nil {
		return "", err
	}

	return localBundle, nil
}

// getCurrentActivation gets the latest activation.
func getCurrentActivation(ctx context.Context, client edgeworkers.Edgeworkers, edgeworkerID int, network, version string) (*edgeworkers.Activation, error) {
	activationsResp, err := client.ListActivations(ctx, edgeworkers.ListActivationsRequest{
		EdgeWorkerID: edgeworkerID,
		Version:      version,
	})
	if err != nil {
		return nil, err
	}

	activations := filterActivationsByNetwork(activationsResp.Activations, network)
	completedActivations := filterActivationsByCompletedStatus(activations)
	finalActivations := sortActivationsByDate(completedActivations)

	if len(finalActivations) != 0 {
		return &finalActivations[0], nil
	}

	return nil, nil
}

func filterActivationsByNetwork(acts []edgeworkers.Activation, net string) (activations []edgeworkers.Activation) {
	for _, act := range acts {
		if act.Network == net {
			activations = append(activations, act)
		}
	}
	return activations
}

func filterActivationsByCompletedStatus(activations []edgeworkers.Activation) (completedActivations []edgeworkers.Activation) {
	for _, activation := range activations {
		if activation.Status == activationStatusComplete {
			completedActivations = append(completedActivations, activation)
		}
	}
	return completedActivations
}

func sortActivationsByDate(activations []edgeworkers.Activation) []edgeworkers.Activation {
	sort.Slice(activations, func(i, j int) bool {
		timeLayout := time.RFC3339
		t1, err := time.Parse(timeLayout, activations[i].CreatedTime)
		if err != nil {
			panic(err)
		}
		t2, err := time.Parse(timeLayout, activations[j].CreatedTime)
		if err != nil {
			panic(err)
		}
		return t1.After(t2)
	})
	return activations
}
