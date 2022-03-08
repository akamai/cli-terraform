package imaging

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"text/template"

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
		Policies  []TFPolicy
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

	// TFPolicy represents policy data used in templates
	TFPolicy struct {
		PolicyID             string
		ActivateOnProduction bool
		JSON                 string
	}
)

//go:embed templates/*
var templateFiles embed.FS

var (

	// RemoveSymbols is a regexp used to remove special characters from policy json file names.
	RemoveSymbols = regexp.MustCompile(`[^\w]`)
	// ErrFetchingPolicySet is returned when fetching policy set fails
	ErrFetchingPolicySet = errors.New("unable to fetch policy set with given name")
	// ErrFetchingPolicy is returned when fetching policy set fails
	ErrFetchingPolicy = errors.New("unable to fetch policy with given name")
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

	jsonDir := tools.TFWorkPath
	if c.IsSet("policy-json-dir") {
		jsonDir = c.String("policy-json-dir")
	}
	jsonDir = filepath.FromSlash(jsonDir)
	if stat, err := os.Stat(jsonDir); err != nil || !stat.IsDir() {
		return cli.NewExitError(color.RedString("Policy JSON path is not accessible"), 1)
	}

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
		AdditionalFuncs: template.FuncMap{
			"ToLower": func(val string) string {
				return strings.ToLower(val)
			},
			"RemoveSymbols": func(val string) string {
				return RemoveSymbols.ReplaceAllString(val, "_")
			},
		},
	}

	contractID, policySetID := c.Args().Get(0), c.Args().Get(1)
	section := edgegrid.GetEdgercSection(c)
	if err = createImaging(ctx, contractID, policySetID, jsonDir, section, client, processor); err != nil {
		return cli.Exit(color.RedString(fmt.Sprintf("Error exporting policy HCL: %s", err)), 1)
	}
	return nil
}

func createImaging(ctx context.Context, contractID, policySetID, jsonDir, section string, client imaging.Imaging, templateProcessor templates.TemplateProcessor) error {
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
	term.Spinner().OK()

	term.Spinner().Start("Fetching policies for the given policy set " + policySetID)
	policies, err := getPolicies(ctx, policySetID, contractID, client)
	if err != nil {
		term.Spinner().Fail()
		return fmt.Errorf("%w: %s", ErrFetchingPolicy, err)
	}

	var tfPoliciesData []TFPolicy
	switch policySet.Type {
	case string(imaging.TypeImage):
		tfPoliciesData, err = getPoliciesImageData(ctx, policies, policySetID, contractID, jsonDir, client)
	case string(imaging.TypeVideo):
		tfPoliciesData, err = getPoliciesVideoData(ctx, policies, policySetID, contractID, jsonDir, client)
	}
	if err != nil {
		term.Spinner().Fail()
		return fmt.Errorf("%w: %s", ErrFetchingPolicy, err)
	}

	term.Spinner().OK()

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

	// Only add Policies if at least one exists
	if tfPoliciesData != nil {
		tfData.Policies = tfPoliciesData
	}

	term.Spinner().Start("Saving TF configurations ")
	if err := templateProcessor.ProcessTemplates(tfData); err != nil {
		term.Spinner().Fail()
		return err
	}
	term.Spinner().OK()
	fmt.Printf("Terraform configuration for policy set '%s' was saved successfully\n", policySet.ID)

	return nil
}

func getPolicies(ctx context.Context, policySetID, contractID string, client imaging.Imaging) ([]imaging.PolicyOutput, error) {
	stagingPolicies, err := client.ListPolicies(ctx, imaging.ListPoliciesRequest{
		Network:     imaging.PolicyNetworkStaging,
		PolicySetID: policySetID,
		ContractID:  contractID,
	})
	if err != nil {
		return nil, err
	}

	return stagingPolicies.Items, nil
}

func getPoliciesImageData(ctx context.Context, policies []imaging.PolicyOutput, policySetID, contractID, jsonDir string, client imaging.Imaging) ([]TFPolicy, error) {
	var tfPoliciesData []TFPolicy

	for _, policyOutput := range policies {
		policy, ok := policyOutput.(*imaging.PolicyOutputImage)
		if !ok {
			return nil, fmt.Errorf("policy is not of type image")
		}

		policyJSON, err := getPolicyImageJSON(policy)
		if err != nil {
			return nil, err
		}

		var activateOnProduction bool
		policyProductionOutput, err := client.GetPolicy(ctx, imaging.GetPolicyRequest{
			PolicyID:    policy.ID,
			Network:     imaging.PolicyNetworkProduction,
			ContractID:  contractID,
			PolicySetID: policySetID,
		})
		if err != nil {
			var e *imaging.Error
			if ok := errors.As(err, &e); !ok || e.Status != http.StatusNotFound {
				return nil, err
			}
		} else {
			policyProduction, ok := policyProductionOutput.(*imaging.PolicyOutputImage)
			if !ok {
				return nil, fmt.Errorf("policy is not of type image")
			}
			policyProductionJSON, err := getPolicyImageJSON(policyProduction)
			if err != nil {
				return nil, err
			}

			activateOnProduction, err = equalPolicyImage(policyJSON, policyProductionJSON)
			if err != nil {
				return nil, err
			}
		}

		jsonPath := filepath.Join(jsonDir, RemoveSymbols.ReplaceAllString(policy.ID, "_")+".json")
		err = ioutil.WriteFile(jsonPath, []byte(policyJSON), 0755)
		if err != nil {
			return nil, err
		}

		tfPoliciesData = append(tfPoliciesData, TFPolicy{
			PolicyID:             policy.ID,
			ActivateOnProduction: activateOnProduction,
			JSON:                 jsonPath,
		})
	}

	return tfPoliciesData, nil
}

func getPoliciesVideoData(ctx context.Context, policies []imaging.PolicyOutput, policySetID, contractID, jsonDir string, client imaging.Imaging) ([]TFPolicy, error) {
	var tfPoliciesData []TFPolicy

	for _, policyOutput := range policies {
		policy, ok := policyOutput.(*imaging.PolicyOutputVideo)
		if !ok {
			return nil, fmt.Errorf("policy is not of type video")
		}

		policyJSON, err := getPolicyVideoJSON(policy)
		if err != nil {
			return nil, err
		}

		var activateOnProduction bool
		policyProductionOutput, err := client.GetPolicy(ctx, imaging.GetPolicyRequest{
			PolicyID:    policy.ID,
			Network:     imaging.PolicyNetworkProduction,
			ContractID:  contractID,
			PolicySetID: policySetID,
		})
		if err != nil {
			var e *imaging.Error
			if ok := errors.As(err, &e); !ok || e.Status != http.StatusNotFound {
				return nil, err
			}
		} else {
			policyProduction, ok := policyProductionOutput.(*imaging.PolicyOutputVideo)
			if !ok {
				return nil, fmt.Errorf("policy is not of type video")
			}
			policyProductionJSON, err := getPolicyVideoJSON(policyProduction)
			if err != nil {
				return nil, err
			}

			activateOnProduction, err = equalPolicyVideo(policyJSON, policyProductionJSON)
			if err != nil {
				return nil, err
			}
		}

		jsonPath := filepath.Join(jsonDir, RemoveSymbols.ReplaceAllString(policy.ID, "_")+".json")
		err = ioutil.WriteFile(jsonPath, []byte(policyJSON), 0755)
		if err != nil {
			return nil, err
		}

		tfPoliciesData = append(tfPoliciesData, TFPolicy{
			PolicyID:             policy.ID,
			ActivateOnProduction: activateOnProduction,
			JSON:                 jsonPath,
		})
	}

	return tfPoliciesData, nil
}

func getPolicyImageJSON(policy *imaging.PolicyOutputImage) (string, error) {
	policyJSON, err := json.MarshalIndent(policy, "", "  ")
	if err != nil {
		return "", err
	}

	// we store JSON as PolicyInput, so we need to convert it from PolicyOutput via JSON representation
	var policyInput imaging.PolicyInputImage
	if err := json.Unmarshal(policyJSON, &policyInput); err != nil {
		return "", err
	}

	policyJSON, err = json.MarshalIndent(policyInput, "", "  ")
	if err != nil {
		return "", err
	}

	return string(policyJSON), nil
}

func getPolicyVideoJSON(policy *imaging.PolicyOutputVideo) (string, error) {
	policyJSON, err := json.MarshalIndent(policy, "", "  ")
	if err != nil {
		return "", err
	}

	// we store JSON as PolicyInput, so we need to convert it from PolicyOutput via JSON representation
	var policyInput imaging.PolicyInputVideo
	if err := json.Unmarshal(policyJSON, &policyInput); err != nil {
		return "", err
	}

	policyJSON, err = json.MarshalIndent(policyInput, "", "  ")
	if err != nil {
		return "", err
	}

	return string(policyJSON), nil
}

func equalPolicyImage(old, new string) (bool, error) {
	if old == new {
		return true, nil
	}
	var oldPolicy, newPolicy imaging.PolicyInputImage
	if old == "" || new == "" {
		return old == new, nil
	}
	if err := json.Unmarshal([]byte(old), &oldPolicy); err != nil {
		return false, err
	}
	if err := json.Unmarshal([]byte(new), &newPolicy); err != nil {
		return false, err
	}

	return reflect.DeepEqual(oldPolicy, newPolicy), nil
}

func equalPolicyVideo(old, new string) (bool, error) {
	if old == new {
		return true, nil
	}
	var oldPolicy, newPolicy imaging.PolicyInputVideo
	if old == "" || new == "" {
		return old == new, nil
	}
	if err := json.Unmarshal([]byte(old), &oldPolicy); err != nil {
		return false, err
	}
	if err := json.Unmarshal([]byte(new), &newPolicy); err != nil {
		return false, err
	}

	return reflect.DeepEqual(oldPolicy, newPolicy), nil
}
