package imaging

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/imaging"
	"github.com/akamai/cli-terraform/pkg/edgegrid"
	"github.com/akamai/cli-terraform/pkg/templates"
	"github.com/akamai/cli-terraform/pkg/tools"
	"github.com/akamai/cli/pkg/terminal"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

type (
	// TFImagingData represents the data used in imaging templates
	TFImagingData struct {
		PolicySet TFPolicySet
		Section   string
	}

	// TFPolicySet represents policy set data used in templates
	TFPolicySet struct {
		ID         string
		ContractID string
		Name       string
		Region     string
		Type       string
	}
)

//go:embed templates/*
var templateFiles embed.FS

var (
	// ErrFetchingPolicySet is returned when fetching policy set fails
	ErrFetchingPolicySet = errors.New("unable to fetch policy set with given name")
)

// CmdCreateImaging is an entrypoint to create-imaging command
func CmdCreateImaging(c *cli.Context) error {
	ctx := c.Context
	if c.NArg() < 2 {
		if c.NArg() == 0 {
			if err := cli.ShowCommandHelp(c, c.Command.Name); err != nil {
				return cli.Exit(color.RedString("Error displaying help command"), 1)
			}
		}
		return cli.Exit(color.RedString("Contract id and policy set id are required"), 1)
	}

	sess := edgegrid.GetSession(ctx)
	client := imaging.Client(sess)
	if c.IsSet("tfworkpath") {
		tools.TFWorkPath = c.String("tfworkpath")
	}
	tools.TFWorkPath = filepath.FromSlash(tools.TFWorkPath)
	if stat, err := os.Stat(tools.TFWorkPath); err != nil || !stat.IsDir() {
		return cli.Exit(color.RedString("Destination work path is not accessible"), 1)
	}

	imagingPath := filepath.Join(tools.TFWorkPath, "imaging.tf")
	variablesPath := filepath.Join(tools.TFWorkPath, "variables.tf")
	importPath := filepath.Join(tools.TFWorkPath, "import.sh")

	err := tools.CheckFiles(imagingPath, variablesPath, importPath)
	if err != nil {
		return cli.Exit(color.RedString(err.Error()), 1)
	}
	templateToFile := map[string]string{
		"imaging.tmpl":   imagingPath,
		"variables.tmpl": variablesPath,
		"imports.tmpl":   importPath,
	}

	processor := templates.FSTemplateProcessor{
		TemplatesFS:     templateFiles,
		TemplateTargets: templateToFile,
	}

	contractID, policySetID := c.Args().Get(0), c.Args().Get(1)
	section := edgegrid.GetEdgercSection(c)
	if err = createImaging(ctx, contractID, policySetID, section, client, processor); err != nil {
		return cli.Exit(color.RedString(fmt.Sprintf("Error exporting policy HCL: %s", err)), 1)
	}
	return nil
}

func createImaging(ctx context.Context, contractID, policySetID, section string, client imaging.Imaging, templateProcessor templates.TemplateProcessor) error {
	term := terminal.Get(ctx)

	fmt.Println("Exporting Image and Video Manager configuration")
	term.Spinner().Start("Fetching policy set " + policySetID)

	policySet, err := client.GetPolicySet(ctx, imaging.GetPolicySetRequest{
		PolicySetID: policySetID,
		ContractID:  contractID,
	})
	if err != nil {
		term.Spinner().Fail()
		return fmt.Errorf("%w: %s", ErrFetchingPolicySet, err)
	}

	tfData := TFImagingData{
		PolicySet: TFPolicySet{
			ID:         policySet.ID,
			ContractID: contractID,
			Name:       policySet.Name,
			Region:     string(policySet.Region),
			Type:       policySet.Type,
		},
		Section: section,
	}
	term.Spinner().OK()
	term.Spinner().Start("Saving TF configurations ")
	if err := templateProcessor.ProcessTemplates(tfData); err != nil {
		term.Spinner().Fail()
		return err
	}
	term.Spinner().OK()
	fmt.Printf("Terraform configuration for policy set '%s' was saved successfully\n", policySet.ID)

	return nil
}
