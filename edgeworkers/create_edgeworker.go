package edgeworkers

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"text/template"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/edgeworkers"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	common "github.com/akamai/cli-common-golang"
	"github.com/akamai/cli-terraform/templates"
	"github.com/akamai/cli-terraform/tools"
	"github.com/fatih/color"
	"github.com/urfave/cli"
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
	}
)

var (
	// ErrFetchingEdgeWorker is returned when fetching edgeworker fails
	ErrFetchingEdgeWorker = errors.New("unable to fetch edgeworker with given edgeworker_id")
)

// CmdCreateEdgeWorker is an entrypoint to create-edgeworker command
func CmdCreateEdgeWorker(c *cli.Context) error {
	// TODO context should be retrieved from cli.Context once we migrate to urfave/cli v2
	ctx := context.Background()
	if c.NArg() == 0 {
		if err := cli.ShowCommandHelp(c, c.Command.Name); err != nil {
			return cli.NewExitError(color.RedString("Error displaying help command"), 1)
		}
		return cli.NewExitError(color.RedString("EdgeWorker id is required"), 1)
	}
	config, err := tools.GetEdgegridConfig(c)
	if err != nil {
		return err
	}

	sess, err := session.New(
		session.WithSigner(config),
	)
	if err != nil {
		return cli.NewExitError(color.RedString(err.Error()), 1)
	}
	client := edgeworkers.Client(sess)

	if c.IsSet("tfworkpath") {
		tools.TFWorkPath = c.String("tfworkpath")
	}
	tools.TFWorkPath = filepath.FromSlash(tools.TFWorkPath)
	if stat, err := os.Stat(tools.TFWorkPath); err != nil || !stat.IsDir() {
		return cli.NewExitError(color.RedString("Destination work path is not accessible"), 1)
	}

	edgeWorkerPath := filepath.Join(tools.TFWorkPath, "edgeworker.tf")
	variablesPath := filepath.Join(tools.TFWorkPath, "variables.tf")
	importPath := filepath.Join(tools.TFWorkPath, "import.sh")

	bundleDir := tools.TFWorkPath
	if c.IsSet("bundlepath") {
		bundleDir = c.String("bundlepath")
	}
	bundleDir = filepath.FromSlash(bundleDir)
	if stat, err := os.Stat(bundleDir); err != nil || !stat.IsDir() {
		return cli.NewExitError(color.RedString("Bundle path is not accessible"), 1)
	}

	err = tools.CheckFiles(edgeWorkerPath, variablesPath, importPath)
	if err != nil {
		return cli.NewExitError(color.RedString(err.Error()), 1)
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
		},
	}

	edgeWorkerID, err := strconv.Atoi(c.Args().First())
	if err != nil {
		return cli.NewExitError(color.RedString("edgeworker_id is not a valid integer"), 1)
	}
	section := tools.GetEdgercSection(c)

	if err = createEdgeWorker(ctx, edgeWorkerID, bundleDir, section, client, processor); err != nil {
		return cli.NewExitError(color.RedString(fmt.Sprintf("Error exporting edgeworker HCL: %s", err)), 1)
	}
	return nil
}

func createEdgeWorker(ctx context.Context, edgeWorkerID int, bundleDir, section string, client edgeworkers.Edgeworkers, templateProcessor templates.TemplateProcessor) error {
	fmt.Println("Configuring EdgeWorker")
	common.StartSpinner(fmt.Sprintf("Fetching EdgeWorker %d", edgeWorkerID), "")

	edgeWorker, err := client.GetEdgeWorkerID(ctx, edgeworkers.GetEdgeWorkerIDRequest{
		EdgeWorkerID: edgeWorkerID,
	})
	if err != nil {
		common.StopSpinnerFail()
		return fmt.Errorf("%w: %s", ErrFetchingEdgeWorker, err)
	}

	localBundle, err := getEdgeWorkerBundle(ctx, edgeWorkerID, bundleDir, client)
	if err != nil {
		common.StopSpinnerFail()
		return fmt.Errorf("%w: %s", ErrFetchingEdgeWorker, err)
	}

	tfEdgeWorkerData := TFEdgeWorkerData{
		EdgeWorkerID:   edgeWorkerID,
		Name:           edgeWorker.Name,
		GroupID:        edgeWorker.GroupID,
		ResourceTierID: edgeWorker.ResourceTierID,
		LocalBundle:    localBundle,
		Section:        section,
	}

	common.StopSpinnerOk()
	common.StartSpinner("Saving TF configurations ", "")
	if err := templateProcessor.ProcessTemplates(tfEdgeWorkerData); err != nil {
		common.StopSpinnerFail()
		return fmt.Errorf("%w: %s", ErrSavingFiles, err)
	}
	common.StopSpinnerOk()
	fmt.Printf("Terraform configuration for edgeworker '%s' with edgeworker_id '%d' was saved successfully\n", edgeWorker.Name, edgeWorkerID)

	return nil
}

// getEdgeWorkerBundle fetches the bundle content of the latest version given edgeWorkerID and returns the path to it
func getEdgeWorkerBundle(ctx context.Context, edgeWorkerID int, bundlePath string, client edgeworkers.Edgeworkers) (string, error) {
	versions, err := client.ListEdgeWorkerVersions(ctx, edgeworkers.ListEdgeWorkerVersionsRequest{
		EdgeWorkerID: edgeWorkerID,
	})
	if err != nil {
		return "", err
	}

	if len(versions.EdgeWorkerVersions) == 0 {
		return "", nil
	}

	var version string
	var createdTime time.Time
	for _, v := range versions.EdgeWorkerVersions {
		parsedCreatedTime, err := time.Parse(time.RFC3339, v.CreatedTime)
		if err != nil {
			return "", fmt.Errorf("cannot indicate the latest version for edgeworker bundle")
		}
		if parsedCreatedTime.After(createdTime) {
			createdTime = parsedCreatedTime
			version = v.Version
		}
	}
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
