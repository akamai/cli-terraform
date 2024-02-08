// Package cloudlets contains code for exporting cloudlets configuration
package cloudlets

import (
	"context"
	"embed"
	"errors"
	"fmt"
	"path/filepath"
	"reflect"
	"sort"
	"text/template"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/cloudlets"
	v3 "github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/cloudlets/v3"
	"github.com/akamai/cli-terraform/pkg/edgegrid"
	"github.com/akamai/cli-terraform/pkg/templates"
	"github.com/akamai/cli-terraform/pkg/tools"
	"github.com/akamai/cli/pkg/terminal"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

type (
	// TFPolicyData represents the data used in policy templates
	TFPolicyData struct {
		Name                    string
		CloudletCode            string
		OriginDescription       string
		Description             string
		GroupID                 int64
		MatchRuleFormat         cloudlets.MatchRuleFormat
		MatchRules              any
		PolicyActivations       TFPolicyActivationsData
		LoadBalancers           []LoadBalancerVersion
		LoadBalancerActivations []cloudlets.LoadBalancerActivation
		Section                 string
		IsV3                    bool
	}
	// LoadBalancerVersion ...
	LoadBalancerVersion struct {
		cloudlets.LoadBalancerVersion
		OriginDescription string
	}

	// TFPolicyActivationsData represents data used in policy all activation resource templates
	TFPolicyActivationsData struct {
		Staging    *TFPolicyActivationData
		Production *TFPolicyActivationData
		IsV3       bool
	}

	// TFPolicyActivationData represents data used in policy activation resource templates
	TFPolicyActivationData struct {
		PolicyID   int64
		Version    int64
		Properties []string
		IsV3       bool
	}

	activationStrategy interface {
		initializeWithPolicy(ctx context.Context, policyName string) error
		validatePolicy() error
		populateWithLatestPolicyVersion(ctx context.Context) error
		getTFPolicyData(ctx context.Context, section string) (*TFPolicyData, error)
	}

	v2ActivationStrategy struct {
		client        cloudlets.Cloudlets
		policy        *cloudlets.Policy
		policyVersion *cloudlets.PolicyVersion
	}

	v3ActivationStrategy struct {
		client        v3.Cloudlets
		policy        *v3.Policy
		policyVersion *v3.PolicyVersion
	}
)

//go:embed templates/*
var templateFiles embed.FS

var (
	// ErrFetchingPolicy is returned when fetching policy fails
	ErrFetchingPolicy = errors.New("unable to fetch policy with given name")
	// ErrFetchingVersion is returned when fetching policy version fails
	ErrFetchingVersion = errors.New("unable to fetch latest policy version")
	// ErrCloudletTypeNotSupported is returned when a provided cloudlet type is not yet supported
	ErrCloudletTypeNotSupported = errors.New("cloudlet type not supported")

	errPolicyNotFound   = errors.New("policy does not exist")
	errVersionsNotFound = errors.New("no policy versions found for given policy")
)

// CmdCreatePolicy is an entrypoint to create-policy command
func CmdCreatePolicy(c *cli.Context) error {
	ctx := c.Context
	sess := edgegrid.GetSession(c.Context)
	clientV2 := cloudlets.Client(sess)
	clientV3 := v3.Client(sess)

	// tfWorkPath is a target directory for generated terraform resources
	var tfWorkPath = "./"
	if c.IsSet("tfworkpath") {
		tfWorkPath = c.String("tfworkpath")
	}

	policyPath := filepath.Join(tfWorkPath, "policy.tf")
	matchRulesPath := filepath.Join(tfWorkPath, "match-rules.tf")
	loadBalancerPath := filepath.Join(tfWorkPath, "load-balancer.tf")
	variablesPath := filepath.Join(tfWorkPath, "variables.tf")
	importPath := filepath.Join(tfWorkPath, "import.sh")

	err := tools.CheckFiles(policyPath, matchRulesPath, loadBalancerPath, variablesPath, importPath)
	if err != nil {
		return cli.Exit(color.RedString(err.Error()), 1)
	}
	templateToFile := map[string]string{
		"policy.tmpl":        policyPath,
		"match-rules.tmpl":   matchRulesPath,
		"load-balancer.tmpl": loadBalancerPath,
		"variables.tmpl":     variablesPath,
		"imports.tmpl":       importPath,
	}

	processor := templates.FSTemplateProcessor{
		TemplatesFS:     templateFiles,
		TemplateTargets: templateToFile,
		AdditionalFuncs: template.FuncMap{
			"deepequal": reflect.DeepEqual,
		},
	}

	policyName := c.Args().First()
	section := edgegrid.GetEdgercSection(c)
	if err = createPolicy(ctx, policyName, section, clientV2, clientV3, processor); err != nil {
		return cli.Exit(color.RedString(fmt.Sprintf("Error exporting policy HCL: %s", err)), 1)
	}
	return nil
}

func createPolicy(ctx context.Context, policyName, section string, clientV2 cloudlets.Cloudlets, clientV3 v3.Cloudlets, templateProcessor templates.TemplateProcessor) error {
	term := terminal.Get(ctx)

	fmt.Println("Configuring Policy")
	term.Spinner().Start("Fetching policy " + policyName)

	strategy, err := initializeStrategyForPolicy(ctx, policyName, clientV2, clientV3)
	if err != nil {
		term.Spinner().Fail()
		return fmt.Errorf("%w: %s", ErrFetchingPolicy, err)
	}
	if err = strategy.validatePolicy(); err != nil {
		term.Spinner().Fail()
		return err
	}

	err = strategy.populateWithLatestPolicyVersion(ctx)
	if err != nil {
		term.Spinner().Fail()
		return fmt.Errorf("%w: %s", ErrFetchingVersion, err)
	}

	tfPolicyData, err := strategy.getTFPolicyData(ctx, section)
	if err != nil {
		term.Spinner().Fail()
		return err
	}

	term.Spinner().OK()
	term.Spinner().Start("Saving TF configurations ")
	if err := templateProcessor.ProcessTemplates(*tfPolicyData); err != nil {
		term.Spinner().Fail()
		return err
	}
	term.Spinner().OK()
	fmt.Printf("Terraform configuration for policy '%s' was saved successfully\n", policyName)

	return nil
}

func (strategy *v2ActivationStrategy) getTFPolicyData(ctx context.Context, section string) (*TFPolicyData, error) {
	tfPolicyData := TFPolicyData{
		Section:           section,
		Name:              strategy.policy.Name,
		CloudletCode:      strategy.policy.CloudletCode,
		GroupID:           strategy.policy.GroupID,
		PolicyActivations: TFPolicyActivationsData{IsV3: false},
	}

	if strategy.policyVersion == nil {
		return &tfPolicyData, nil
	}

	tfPolicyData.Description = strategy.policyVersion.Description
	tfPolicyData.MatchRuleFormat = strategy.policyVersion.MatchRuleFormat
	tfPolicyData.MatchRules = strategy.policyVersion.MatchRules

	if activationStaging := getActiveVersionAndProperties(strategy.policy, cloudlets.PolicyActivationNetworkStaging); activationStaging != nil {
		tfPolicyData.PolicyActivations.Staging = activationStaging
	}
	if activationProd := getActiveVersionAndProperties(strategy.policy, cloudlets.PolicyActivationNetworkProduction); activationProd != nil {
		tfPolicyData.PolicyActivations.Production = activationProd
	}

	if tfPolicyData.CloudletCode == "ALB" {
		originIDs, err := getOriginIDs(strategy.policyVersion.MatchRules)
		if err != nil {

			return nil, fmt.Errorf("%w: %s", ErrFetchingVersion, err)
		}
		tfPolicyData.LoadBalancers, err = getLoadBalancers(ctx, strategy.client, originIDs)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrFetchingVersion, err)
		}
		tfPolicyData.LoadBalancerActivations, err = getLoadBalancerActivations(ctx, strategy.client, originIDs)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrFetchingVersion, err)
		}

	}
	return &tfPolicyData, nil
}

func initializeStrategyForPolicy(ctx context.Context, policyName string, clientV2 cloudlets.Cloudlets, clientV3 v3.Cloudlets) (activationStrategy, error) {
	v2Strategy := v2ActivationStrategy{client: clientV2}
	errV2 := v2Strategy.initializeWithPolicy(ctx, policyName)
	if errV2 == nil {
		return &v2Strategy, nil
	}
	v3Strategy := v3ActivationStrategy{client: clientV3}
	errV3 := v3Strategy.initializeWithPolicy(ctx, policyName)
	if errV3 == nil {
		return &v3Strategy, nil
	}
	return nil, fmt.Errorf("could not find policy %s: neither as V2 (%s) nor as V3 (%s)", policyName, errV2, errV3)
}

func getLoadBalancerActivations(ctx context.Context, client cloudlets.Cloudlets, originIDs []string) ([]cloudlets.LoadBalancerActivation, error) {
	activations := make([]cloudlets.LoadBalancerActivation, 0)
	for _, originID := range originIDs {
		activation, err := getApplicationLoadBalancerActivation(ctx, client, originID, cloudlets.LoadBalancerActivationNetworkProduction)
		if err != nil {
			return nil, err
		}
		if activation != nil {
			activations = append(activations, *activation)
		}

		activation, err = getApplicationLoadBalancerActivation(ctx, client, originID, cloudlets.LoadBalancerActivationNetworkStaging)
		if err != nil {
			return nil, err
		}
		if activation != nil {
			activations = append(activations, *activation)
		}
	}
	return activations, nil
}

func getLoadBalancers(ctx context.Context, client cloudlets.Cloudlets, originIDs []string) ([]LoadBalancerVersion, error) {
	loadBalancers := make([]LoadBalancerVersion, 0)
	for _, originID := range originIDs {
		versions, err := client.ListLoadBalancerVersions(ctx, cloudlets.ListLoadBalancerVersionsRequest{
			OriginID: originID,
		})
		if err != nil {
			return nil, err
		}

		var ver int64
		var loadBalancerVersion cloudlets.LoadBalancerVersion
		for _, version := range versions {
			if version.Version > ver {
				ver = version.Version
				loadBalancerVersion = version
			}
		}
		if ver > 0 {
			origin, err := client.GetOrigin(ctx, cloudlets.GetOriginRequest{
				OriginID: originID,
			})
			if err != nil {
				return nil, err
			}

			loadBalancers = append(loadBalancers, LoadBalancerVersion{
				LoadBalancerVersion: loadBalancerVersion,
				OriginDescription:   origin.Description,
			})
		}
	}
	return loadBalancers, nil
}

func getOriginIDs(rules cloudlets.MatchRules) ([]string, error) {
	// the same originID can be assigned to multiple rules, so we need to deduplicate it
	originIDs := map[string]struct{}{}
	for _, rule := range rules {
		ruleALB, ok := rule.(*cloudlets.MatchRuleALB)
		if !ok {
			return nil, fmt.Errorf("match rule type is not a MatchRuleALB: %T", rule)
		}
		originID := ruleALB.ForwardSettings.OriginID
		if originID != "" {
			originIDs[originID] = struct{}{}
		}
	}
	result := make([]string, 0, len(originIDs))
	for originID := range originIDs {
		result = append(result, originID)
	}
	return result, nil
}

func getApplicationLoadBalancerActivation(ctx context.Context, client cloudlets.Cloudlets, originID string, network cloudlets.LoadBalancerActivationNetwork) (*cloudlets.LoadBalancerActivation, error) {
	activations, err := client.ListLoadBalancerActivations(ctx, cloudlets.ListLoadBalancerActivationsRequest{OriginID: originID})
	filteredActivations := make([]cloudlets.LoadBalancerActivation, 0)
	if err != nil {
		return nil, err
	}

	for _, act := range activations {
		if act.Network == network {
			filteredActivations = append(filteredActivations, act)
		}
	}

	// The API is not providing any id to match the status of the activation request within the list of the activation statuses.
	// The recommended solution is to get the newest activation which is most likely the right one.
	// So we sort by ActivatedDate to get the newest activation.
	sort.Slice(filteredActivations, func(i, j int) bool {
		return activations[i].ActivatedDate > activations[j].ActivatedDate
	})

	if len(filteredActivations) > 0 {
		return &filteredActivations[0], nil
	}
	return nil, nil
}

func (strategy *v2ActivationStrategy) initializeWithPolicy(ctx context.Context, name string) error {
	pageSize, offset := 1000, 0
	for {
		policies, err := strategy.client.ListPolicies(ctx, cloudlets.ListPoliciesRequest{
			Offset:   offset,
			PageSize: &pageSize,
		})
		if err != nil {
			return err
		}
		for _, p := range policies {
			if p.Name == name {
				strategy.policy = &p
				return nil
			}
		}
		if len(policies) < pageSize {
			break
		}
		offset += pageSize
	}
	return errPolicyNotFound
}

func (strategy *v2ActivationStrategy) populateWithLatestPolicyVersion(ctx context.Context) error {
	var version *int64
	policyID := strategy.policy.PolicyID
	pageSize, offset := 1000, 0
	for {
		versions, err := strategy.client.ListPolicyVersions(ctx, cloudlets.ListPolicyVersionsRequest{
			PolicyID:     policyID,
			IncludeRules: false,
			PageSize:     &pageSize,
			Offset:       offset,
		})
		if err != nil {
			return err
		}

		if len(versions) == 0 {
			break
		}
		for _, v := range versions {
			v := v
			if version == nil || v.Version > *version {
				version = &v.Version
			}
		}
		if len(versions) < pageSize {
			break
		}
		offset += pageSize
	}
	if version == nil {
		return nil
	}
	policyVersion, err := strategy.client.GetPolicyVersion(ctx, cloudlets.GetPolicyVersionRequest{
		PolicyID: policyID,
		Version:  *version,
	})
	if err != nil {
		return err
	}
	strategy.policyVersion = policyVersion
	return nil
}

func getActiveVersionAndProperties(policy *cloudlets.Policy, network cloudlets.PolicyActivationNetwork) *TFPolicyActivationData {
	var version int64
	var associatedProperties []string
	for _, activation := range policy.Activations {
		if activation.Network != network {
			continue
		}
		version = activation.PolicyInfo.Version
		associatedProperties = append(associatedProperties, activation.PropertyInfo.Name)
	}
	if associatedProperties == nil {
		return nil
	}
	return &TFPolicyActivationData{
		PolicyID:   policy.PolicyID,
		Version:    version,
		Properties: associatedProperties,
	}
}

func (strategy *v2ActivationStrategy) validatePolicy() error {
	var supportedCloudlets = map[string]struct{}{
		"ALB": {},
		"AP":  {},
		"AS":  {},
		"CD":  {},
		"ER":  {},
		"FR":  {},
		"IG":  {},
		"VP":  {},
	}

	_, ok := supportedCloudlets[strategy.policy.CloudletCode]
	if !ok {
		return fmt.Errorf("%w: %s", ErrCloudletTypeNotSupported, strategy.policy.CloudletCode)
	}
	return nil
}

func (strategy *v3ActivationStrategy) initializeWithPolicy(ctx context.Context, policyName string) error {
	page, size := 0, 1000
	for {
		policies, err := strategy.client.ListPolicies(ctx, v3.ListPoliciesRequest{Page: page, Size: size})
		if err != nil {
			return err
		}
		for _, policy := range policies.Content {
			if policy.Name == policyName {
				strategy.policy = &policy
				return nil
			}
		}
		if len(policies.Content) < size {
			break
		}
		page++
	}
	return errPolicyNotFound
}

func (strategy *v3ActivationStrategy) validatePolicy() error {
	var supportedCloudlets = map[v3.CloudletType]struct{}{
		v3.CloudletTypeAP: {},
		v3.CloudletTypeAS: {},
		v3.CloudletTypeCD: {},
		v3.CloudletTypeER: {},
		v3.CloudletTypeFR: {},
		v3.CloudletTypeIG: {},
	}

	_, ok := supportedCloudlets[strategy.policy.CloudletType]
	if !ok {
		return fmt.Errorf("%w: %s", ErrCloudletTypeNotSupported, strategy.policy.CloudletType)
	}
	return nil
}

func (strategy *v3ActivationStrategy) populateWithLatestPolicyVersion(ctx context.Context) error {
	policyID := strategy.policy.ID

	versions, err := strategy.client.ListPolicyVersions(ctx, v3.ListPolicyVersionsRequest{PolicyID: policyID, Page: 0, Size: 10})
	if err != nil {
		return err
	}
	if len(versions.PolicyVersions) == 0 {
		return nil
	}

	policyVersion, err := strategy.client.GetPolicyVersion(ctx, v3.GetPolicyVersionRequest{
		PolicyID:      policyID,
		PolicyVersion: versions.PolicyVersions[0].PolicyVersion,
	})
	if err != nil {
		return err
	}
	strategy.policyVersion = policyVersion
	return nil
}

func (strategy *v3ActivationStrategy) getTFPolicyData(_ context.Context, section string) (*TFPolicyData, error) {
	tfPolicyData := TFPolicyData{
		Section:           section,
		Name:              strategy.policy.Name,
		CloudletCode:      string(strategy.policy.CloudletType),
		GroupID:           strategy.policy.GroupID,
		IsV3:              true,
		PolicyActivations: TFPolicyActivationsData{IsV3: true},
	}

	if strategy.policyVersion != nil && strategy.policyVersion.Description != nil {
		tfPolicyData.Description = *strategy.policyVersion.Description
	}
	if strategy.policyVersion != nil && strategy.policyVersion.MatchRules != nil {
		tfPolicyData.MatchRules = strategy.policyVersion.MatchRules
	}

	if activationStaging := strategy.getActivationDataForV3(v3.StagingNetwork); activationStaging != nil {
		tfPolicyData.PolicyActivations.Staging = activationStaging
	}
	if activationProd := strategy.getActivationDataForV3(v3.ProductionNetwork); activationProd != nil {
		tfPolicyData.PolicyActivations.Production = activationProd
	}

	return &tfPolicyData, nil
}

func (strategy *v3ActivationStrategy) getActivationDataForV3(network v3.Network) *TFPolicyActivationData {
	var activationToCheck *v3.PolicyActivation
	switch network {
	case v3.StagingNetwork:
		activationToCheck = strategy.policy.CurrentActivations.Staging.Effective
	case v3.ProductionNetwork:
		activationToCheck = strategy.policy.CurrentActivations.Production.Effective
	}
	if activationToCheck == nil || activationToCheck.Operation == v3.OperationDeactivation {
		return nil
	}
	return &TFPolicyActivationData{
		PolicyID: activationToCheck.PolicyID,
		Version:  activationToCheck.PolicyVersion,
		IsV3:     true,
	}
}
