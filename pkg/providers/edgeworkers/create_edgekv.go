// Package edgeworkers contains code for exporting edge workers and edge kv configuration
package edgeworkers

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/edgeworkers"
	"github.com/akamai/cli-terraform/pkg/edgegrid"
	"github.com/akamai/cli-terraform/pkg/templates"
	"github.com/akamai/cli-terraform/pkg/tools"
	"github.com/akamai/cli/pkg/terminal"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
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
)

// CmdCreateEdgeKV is an entrypoint to create-edgekv command
func CmdCreateEdgeKV(c *cli.Context) error {
	ctx := c.Context
	sess := edgegrid.GetSession(c.Context)
	client := edgeworkers.Client(sess)

	// tfWorkPath is a target directory for generated terraform resources
	var tfWorkPath = "./"
	if c.IsSet("tfworkpath") {
		tfWorkPath = c.String("tfworkpath")
	}

	edgeKVPath := filepath.Join(tfWorkPath, "edgekv.tf")
	variablesPath := filepath.Join(tfWorkPath, "variables.tf")
	importPath := filepath.Join(tfWorkPath, "import.sh")

	err := tools.CheckFiles(edgeKVPath, variablesPath, importPath)
	if err != nil {
		return cli.Exit(color.RedString(err.Error()), 1)
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
	section := edgegrid.GetEdgercSection(c)

	if err = createEdgeKV(ctx, namespace, network, section, client, processor); err != nil {
		return cli.Exit(color.RedString(fmt.Sprintf("Error exporting edgekv HCL: %s", err)), 1)
	}
	return nil
}

func createEdgeKV(ctx context.Context, namespace string, network edgeworkers.NamespaceNetwork, section string, client edgeworkers.Edgeworkers, templateProcessor templates.TemplateProcessor) error {
	term := terminal.Get(ctx)
	fmt.Println("Configuring EdgeKV")
	term.Spinner().Start("Fetching EdgeKV "+namespace, "")

	edgeKV, err := getEdgeKV(ctx, namespace, network, client)
	if err != nil {
		term.Spinner().Fail()
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

	term.Spinner().OK()
	term.Spinner().Start("Saving TF configurations ")
	if err := templateProcessor.ProcessTemplates(tfEdgeKVData); err != nil {
		term.Spinner().Fail()
		return fmt.Errorf("%w: %s", templates.ErrSavingFiles, err)
	}
	term.Spinner().OK()
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
