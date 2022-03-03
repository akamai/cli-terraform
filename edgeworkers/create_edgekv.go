package edgeworkers

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/edgeworkers"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/session"
	common "github.com/akamai/cli-common-golang"
	"github.com/akamai/cli-terraform/templates"
	"github.com/akamai/cli-terraform/tools"
	"github.com/fatih/color"
	"github.com/urfave/cli"
)

type (
	// TFEdgeKVData represents the data used in EdgeKV templates
	TFEdgeKVData struct {
		Name        string
		Network     edgeworkers.NamespaceNetwork
		GroupID     int
		Retention   int
		GeoLocation string
		Section     string
	}
)

var (
	//go:embed templates/*
	templateFiles embed.FS

	// ErrFetchingEdgeKV is returned when fetching edgekv fails
	ErrFetchingEdgeKV = errors.New("unable to fetch edgekv with given namespace_name and network")
	// ErrSavingFiles is returned when an issue with processing templates occurs
	ErrSavingFiles = errors.New("saving terraform project files")
)

// CmdCreateEdgeKV is an entrypoint to create-edgekv command
func CmdCreateEdgeKV(c *cli.Context) error {
	// TODO context should be retrieved from cli.Context once we migrate to urfave/cli v2
	ctx := context.Background()
	if c.NArg() < 2 {
		if err := cli.ShowCommandHelp(c, c.Command.Name); err != nil {
			return cli.NewExitError(color.RedString("Error displaying help command"), 1)
		}
		return cli.NewExitError(color.RedString("EdgeKV namespace_name and network are required"), 1)
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

	edgeKVPath := filepath.Join(tools.TFWorkPath, "edgekv.tf")
	variablesPath := filepath.Join(tools.TFWorkPath, "variables.tf")
	importPath := filepath.Join(tools.TFWorkPath, "import.sh")

	err = tools.CheckFiles(edgeKVPath, variablesPath, importPath)
	if err != nil {
		return cli.NewExitError(color.RedString(err.Error()), 1)
	}
	templateToFile := map[string]string{
		"edgekv.tmpl":           edgeKVPath,
		"edgekv-variables.tmpl": variablesPath,
		"edgekv-imports.tmpl":   importPath,
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

	namespace := c.Args().First()
	network := edgeworkers.NamespaceNetwork(c.Args().Get(1))
	section := tools.GetEdgercSection(c)

	if err = createEdgeKV(ctx, namespace, network, section, client, processor); err != nil {
		return cli.NewExitError(color.RedString(fmt.Sprintf("Error exporting edgekv HCL: %s", err)), 1)
	}
	return nil
}

func createEdgeKV(ctx context.Context, namespace string, network edgeworkers.NamespaceNetwork, section string, client edgeworkers.Edgeworkers, templateProcessor templates.TemplateProcessor) error {
	fmt.Println("Configuring EdgeKV")
	common.StartSpinner("Fetching EdgeKV "+namespace, "")

	edgeKV, err := getEdgeKV(ctx, namespace, network, client)
	if err != nil {
		common.StopSpinnerFail()
		return fmt.Errorf("%w: %s", ErrFetchingEdgeKV, err)
	}

	tfEdgeKVData := TFEdgeKVData{
		Name:        edgeKV.Name,
		Network:     network,
		Retention:   *edgeKV.Retention,
		GeoLocation: edgeKV.GeoLocation,
		Section:     section,
	}

	// Only add GroupID if the API returns it
	if edgeKV.GroupID != nil {
		tfEdgeKVData.GroupID = *edgeKV.GroupID
	}

	common.StopSpinnerOk()
	common.StartSpinner("Saving TF configurations ", "")
	if err := templateProcessor.ProcessTemplates(tfEdgeKVData); err != nil {
		common.StopSpinnerFail()
		return fmt.Errorf("%w: %s", ErrSavingFiles, err)
	}
	common.StopSpinnerOk()
	fmt.Printf("Terraform configuration for edgeKV '%s' on network '%s' was saved successfully\n", edgeKV.Name, network)

	return nil
}

func getEdgeKV(ctx context.Context, namespace string, network edgeworkers.NamespaceNetwork, client edgeworkers.Edgeworkers) (*edgeworkers.Namespace, error) {
	edgeKV, err := client.GetEdgeKVNamespace(ctx, edgeworkers.GetEdgeKVNamespaceRequest{
		Network: network,
		Name:    namespace,
	})
	if err != nil {
		return nil, err
	}
	if edgeKV == nil {
		return nil, fmt.Errorf("edgeKV '%s' on network '%s' does not exist", namespace, network)
	}

	return edgeKV, nil
}
