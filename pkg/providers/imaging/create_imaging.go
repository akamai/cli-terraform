// Package imaging contains code for exporting policies for images and videos
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
	"path"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"text/template"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/imaging"
	"github.com/akamai/cli-terraform/pkg/edgegrid"
	"github.com/akamai/cli-terraform/pkg/templates"
	"github.com/akamai/cli-terraform/pkg/tools"
	"github.com/akamai/cli/pkg/color"
	"github.com/akamai/cli/pkg/terminal"
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
		Policy               imaging.PolicyInput
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
	// ErrCreateDir is returned when error occurred creating directory
	ErrCreateDir = errors.New("cannot create directory")
)

// maxDepth value has to match the MaxPolicyDepth value in terraform imaging subprovider
const maxDepth = 7

// CmdCreateImaging is an entrypoint to create-imaging command
func CmdCreateImaging(c *cli.Context) error {
	ctx := c.Context
	sess := edgegrid.GetSession(ctx)
	client := imaging.Client(sess)

	// tfWorkPath is a target directory for generated terraform resources
	var tfWorkPath = "./"
	if c.IsSet("tfworkpath") {
		tfWorkPath = c.String("tfworkpath")
	}
	imagingPath := filepath.Join(tfWorkPath, "imaging.tf")
	variablesPath := filepath.Join(tfWorkPath, "variables.tf")
	importPath := filepath.Join(tfWorkPath, "import.sh")

	err := tools.CheckFiles(imagingPath, variablesPath, importPath)
	if err != nil {
		return cli.Exit(color.RedString(err.Error()), 1)
	}

	jsonDir := "."
	if c.IsSet("policy-json-dir") {
		jsonDir = c.String("policy-json-dir")
	}
	if !c.IsSet("policies-as-hcl") {
		jsonDirPath := path.Join(tfWorkPath, jsonDir)
		err = ensureDirExists(jsonDirPath)
		if err != nil {
			return cli.Exit(color.RedString(err.Error()), 1)
		}
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
	if err = createImaging(ctx, contractID, policySetID, tfWorkPath, jsonDir, section, client, processor, tools.PolicyAsHCL); err != nil {
		return cli.Exit(color.RedString(fmt.Sprintf("Error exporting policy HCL: %s", err)), 1)
	}
	return nil
}

func createImaging(ctx context.Context, contractID, policySetID, tfWorkPath, jsonDir, section string, client imaging.Imaging, templateProcessor templates.TemplateProcessor, policyAsHCL bool) error {
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
		tfPoliciesData, err = getPoliciesImageData(ctx, policies, policySetID, contractID, tfWorkPath, jsonDir, client, policyAsHCL)
	case string(imaging.TypeVideo):
		tfPoliciesData, err = getPoliciesVideoData(ctx, policies, policySetID, contractID, tfWorkPath, jsonDir, client, policyAsHCL)
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

func ensureDirExists(dirPath string) error {
	dirPath = filepath.FromSlash(dirPath)

	err := os.MkdirAll(dirPath, 0755)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrCreateDir, err)
	}
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

func getPoliciesImageData(ctx context.Context, policies []imaging.PolicyOutput, policySetID, contractID, tfWorkPath, jsonDir string, client imaging.Imaging, policyAsHCL bool) ([]TFPolicy, error) {
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

		if policyAsHCL {
			// we store JSON as PolicyInput, so we need to convert it from PolicyInput via JSON representation
			var policyInput imaging.PolicyInputImage
			if err := json.Unmarshal([]byte(policyJSON), &policyInput); err != nil {
				return nil, err
			}
			if depth := getDepth(policyInput, 0); depth > maxDepth {
				return nil, fmt.Errorf("policy has %d nested transformation levels, while only %d are allowed; please use JSON format instead", depth, maxDepth)
			}
			tfPoliciesData = append(tfPoliciesData, TFPolicy{
				PolicyID:             policy.ID,
				ActivateOnProduction: activateOnProduction,
				Policy:               &policyInput,
			})
		} else {
			jsonPath := filepath.Join(jsonDir, RemoveSymbols.ReplaceAllString(policy.ID, "_")+".json")
			err = ioutil.WriteFile(filepath.Join(tfWorkPath, jsonPath), []byte(policyJSON), 0644)
			if err != nil {
				return nil, err
			}

			tfPoliciesData = append(tfPoliciesData, TFPolicy{
				PolicyID:             policy.ID,
				ActivateOnProduction: activateOnProduction,
				JSON:                 jsonPath,
			})
		}
	}

	return tfPoliciesData, nil
}

func getPoliciesVideoData(ctx context.Context, policies []imaging.PolicyOutput, policySetID, contractID, tfWorkPath, jsonDir string, client imaging.Imaging, policyAsHCL bool) ([]TFPolicy, error) {
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

		if policyAsHCL {
			// we store JSON as PolicyInput, so we need to convert it from PolicyInput via JSON representation
			var policyInput imaging.PolicyInputVideo
			if err := json.Unmarshal([]byte(policyJSON), &policyInput); err != nil {
				return nil, err
			}
			tfPoliciesData = append(tfPoliciesData, TFPolicy{
				PolicyID:             policy.ID,
				ActivateOnProduction: activateOnProduction,
				Policy:               &policyInput,
			})
		} else {
			jsonPath := filepath.Join(jsonDir, RemoveSymbols.ReplaceAllString(policy.ID, "_")+".json")
			err = ioutil.WriteFile(filepath.Join(tfWorkPath, jsonPath), []byte(policyJSON), 0644)
			if err != nil {
				return nil, err
			}

			tfPoliciesData = append(tfPoliciesData, TFPolicy{
				PolicyID:             policy.ID,
				ActivateOnProduction: activateOnProduction,
				JSON:                 jsonPath,
			})
		}
	}

	return tfPoliciesData, nil
}

func getPolicyImageJSON(policy *imaging.PolicyOutputImage) (string, error) {
	policyJSON, err := json.MarshalIndent(policy, "", "  ")
	if err != nil {
		return "", err
	}

	// we store JSON as PolicyInput, so we need to convert it from PolicyInput via JSON representation
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

	// we store JSON as PolicyInput, so we need to convert it from PolicyInput via JSON representation
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

func getDepth(policy interface{}, depth int) int {
	if policy == nil {
		return 0
	}
	v := reflect.ValueOf(policy)
	if v.Kind() == reflect.Ptr {
		vp := v.Elem().Interface()
		v = reflect.ValueOf(vp)
	}
	newDepth := depth
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		var localDepth int
		switch field.Type().String() {
		case "imaging.Transformations":
			transformations := field.Interface().(imaging.Transformations)
			for _, transformation := range transformations {
				d := getDepth(transformation, depth+2)
				if d > localDepth {
					localDepth = d
				}
			}
		case "imaging.TransformationType":
			value := field.Interface()
			if value != nil {
				localDepth = getDepth(value.(imaging.TransformationType), depth+1)
			}
		case "imaging.ImageType":
			value := field.Interface()
			if value != nil {
				localDepth = getDepth(field.Interface().(imaging.ImageType), depth+1)
			}
		}
		if localDepth > newDepth {
			newDepth = localDepth
		}
	}
	return newDepth
}
