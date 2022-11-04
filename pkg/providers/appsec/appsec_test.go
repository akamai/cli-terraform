//nolint:revive
package appsec

import (
	"context"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v3/pkg/appsec"
	"github.com/stretchr/testify/mock"
)

type mockAppsec struct {
	mock.Mock
}

func (m *mockAppsec) UpdateWAPSelectedHostnames(ctx context.Context, req appsec.UpdateWAPSelectedHostnamesRequest) (*appsec.UpdateWAPSelectedHostnamesResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.UpdateWAPSelectedHostnamesResponse), args.Error(1)
}

func (m *mockAppsec) UpdateWAFProtection(ctx context.Context, req appsec.UpdateWAFProtectionRequest) (*appsec.UpdateWAFProtectionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.UpdateWAFProtectionResponse), args.Error(1)
}

func (m *mockAppsec) UpdateWAFMode(ctx context.Context, req appsec.UpdateWAFModeRequest) (*appsec.UpdateWAFModeResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.UpdateWAFModeResponse), args.Error(1)
}

func (m *mockAppsec) UpdateVersionNotes(ctx context.Context, req appsec.UpdateVersionNotesRequest) (*appsec.UpdateVersionNotesResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.UpdateVersionNotesResponse), args.Error(1)
}

func (p *mockAppsec) GetRuleRecommendations(ctx context.Context, params appsec.GetRuleRecommendationsRequest) (*appsec.GetRuleRecommendationsResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetRuleRecommendationsResponse), args.Error(1)
}

func (m *mockAppsec) UpdateThreatIntel(ctx context.Context, req appsec.UpdateThreatIntelRequest) (*appsec.UpdateThreatIntelResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.UpdateThreatIntelResponse), args.Error(1)
}

func (m *mockAppsec) UpdateSlowPostProtectionSetting(ctx context.Context, req appsec.UpdateSlowPostProtectionSettingRequest) (*appsec.UpdateSlowPostProtectionSettingResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.UpdateSlowPostProtectionSettingResponse), args.Error(1)
}

func (m *mockAppsec) UpdateSlowPostProtection(ctx context.Context, req appsec.UpdateSlowPostProtectionRequest) (*appsec.UpdateSlowPostProtectionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.UpdateSlowPostProtectionResponse), args.Error(1)
}

func (m *mockAppsec) UpdateSiemSettings(ctx context.Context, req appsec.UpdateSiemSettingsRequest) (*appsec.UpdateSiemSettingsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.UpdateSiemSettingsResponse), args.Error(1)
}

func (m *mockAppsec) UpdateSelectedHostnames(ctx context.Context, req appsec.UpdateSelectedHostnamesRequest) (*appsec.UpdateSelectedHostnamesResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.UpdateSelectedHostnamesResponse), args.Error(1)
}

func (m *mockAppsec) UpdateSelectedHostname(ctx context.Context, req appsec.UpdateSelectedHostnameRequest) (*appsec.UpdateSelectedHostnameResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.UpdateSelectedHostnameResponse), args.Error(1)
}

func (m *mockAppsec) UpdateSecurityPolicy(ctx context.Context, req appsec.UpdateSecurityPolicyRequest) (*appsec.UpdateSecurityPolicyResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.UpdateSecurityPolicyResponse), args.Error(1)
}

func (m *mockAppsec) UpdateRuleUpgrade(ctx context.Context, req appsec.UpdateRuleUpgradeRequest) (*appsec.UpdateRuleUpgradeResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.UpdateRuleUpgradeResponse), args.Error(1)
}

func (m *mockAppsec) UpdateRuleConditionException(ctx context.Context, req appsec.UpdateConditionExceptionRequest) (*appsec.UpdateConditionExceptionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.UpdateConditionExceptionResponse), args.Error(1)
}

func (m *mockAppsec) UpdateRule(ctx context.Context, req appsec.UpdateRuleRequest) (*appsec.UpdateRuleResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.UpdateRuleResponse), args.Error(1)
}

func (m *mockAppsec) UpdateReputationProtection(ctx context.Context, req appsec.UpdateReputationProtectionRequest) (*appsec.UpdateReputationProtectionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.UpdateReputationProtectionResponse), args.Error(1)
}

func (m *mockAppsec) UpdateReputationProfileAction(ctx context.Context, req appsec.UpdateReputationProfileActionRequest) (*appsec.UpdateReputationProfileActionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.UpdateReputationProfileActionResponse), args.Error(1)
}

func (m *mockAppsec) UpdateReputationProfile(ctx context.Context, req appsec.UpdateReputationProfileRequest) (*appsec.UpdateReputationProfileResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.UpdateReputationProfileResponse), args.Error(1)
}

func (m *mockAppsec) UpdateReputationAnalysis(ctx context.Context, req appsec.UpdateReputationAnalysisRequest) (*appsec.UpdateReputationAnalysisResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.UpdateReputationAnalysisResponse), args.Error(1)
}

func (m *mockAppsec) UpdateRateProtection(ctx context.Context, req appsec.UpdateRateProtectionRequest) (*appsec.UpdateRateProtectionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.UpdateRateProtectionResponse), args.Error(1)
}

func (m *mockAppsec) UpdateRatePolicyAction(ctx context.Context, req appsec.UpdateRatePolicyActionRequest) (*appsec.UpdateRatePolicyActionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.UpdateRatePolicyActionResponse), args.Error(1)
}

func (m *mockAppsec) UpdateRatePolicy(ctx context.Context, req appsec.UpdateRatePolicyRequest) (*appsec.UpdateRatePolicyResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.UpdateRatePolicyResponse), args.Error(1)
}

func (m *mockAppsec) UpdatePolicyProtections(ctx context.Context, req appsec.UpdatePolicyProtectionsRequest) (*appsec.PolicyProtectionsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.PolicyProtectionsResponse), args.Error(1)
}

func (m *mockAppsec) UpdatePenaltyBox(ctx context.Context, req appsec.UpdatePenaltyBoxRequest) (*appsec.UpdatePenaltyBoxResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.UpdatePenaltyBoxResponse), args.Error(1)
}

func (m *mockAppsec) UpdateNetworkLayerProtection(ctx context.Context, req appsec.UpdateNetworkLayerProtectionRequest) (*appsec.UpdateNetworkLayerProtectionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.UpdateNetworkLayerProtectionResponse), args.Error(1)
}

func (m *mockAppsec) UpdateMatchTargetSequence(ctx context.Context, req appsec.UpdateMatchTargetSequenceRequest) (*appsec.UpdateMatchTargetSequenceResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.UpdateMatchTargetSequenceResponse), args.Error(1)
}

func (m *mockAppsec) UpdateMatchTarget(ctx context.Context, req appsec.UpdateMatchTargetRequest) (*appsec.UpdateMatchTargetResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.UpdateMatchTargetResponse), args.Error(1)
}

func (m *mockAppsec) UpdateIPGeoProtection(ctx context.Context, req appsec.UpdateIPGeoProtectionRequest) (*appsec.UpdateIPGeoProtectionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.UpdateIPGeoProtectionResponse), args.Error(1)
}

func (m *mockAppsec) UpdateIPGeo(ctx context.Context, req appsec.UpdateIPGeoRequest) (*appsec.UpdateIPGeoResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.UpdateIPGeoResponse), args.Error(1)
}

func (m *mockAppsec) UpdateEvalRule(ctx context.Context, req appsec.UpdateEvalRuleRequest) (*appsec.UpdateEvalRuleResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.UpdateEvalRuleResponse), args.Error(1)
}

func (m *mockAppsec) UpdateEvalProtectHost(ctx context.Context, req appsec.UpdateEvalProtectHostRequest) (*appsec.UpdateEvalProtectHostResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.UpdateEvalProtectHostResponse), args.Error(1)
}

func (m *mockAppsec) UpdateEvalHost(ctx context.Context, req appsec.UpdateEvalHostRequest) (*appsec.UpdateEvalHostResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.UpdateEvalHostResponse), args.Error(1)
}

func (m *mockAppsec) UpdateEvalGroup(ctx context.Context, req appsec.UpdateAttackGroupRequest) (*appsec.UpdateAttackGroupResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.UpdateAttackGroupResponse), args.Error(1)
}

func (m *mockAppsec) UpdateEval(ctx context.Context, req appsec.UpdateEvalRequest) (*appsec.UpdateEvalResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.UpdateEvalResponse), args.Error(1)
}

func (m *mockAppsec) UpdateCustomRuleAction(ctx context.Context, req appsec.UpdateCustomRuleActionRequest) (*appsec.UpdateCustomRuleActionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.UpdateCustomRuleActionResponse), args.Error(1)
}

func (m *mockAppsec) UpdateCustomRule(ctx context.Context, req appsec.UpdateCustomRuleRequest) (*appsec.UpdateCustomRuleResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.UpdateCustomRuleResponse), args.Error(1)
}

func (m *mockAppsec) UpdateCustomDeny(ctx context.Context, req appsec.UpdateCustomDenyRequest) (*appsec.UpdateCustomDenyResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.UpdateCustomDenyResponse), args.Error(1)
}

func (m *mockAppsec) UpdateConfiguration(ctx context.Context, req appsec.UpdateConfigurationRequest) (*appsec.UpdateConfigurationResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.UpdateConfigurationResponse), args.Error(1)
}

func (m *mockAppsec) UpdateBypassNetworkLists(ctx context.Context, req appsec.UpdateBypassNetworkListsRequest) (*appsec.UpdateBypassNetworkListsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.UpdateBypassNetworkListsResponse), args.Error(1)
}

func (m *mockAppsec) UpdateAttackGroup(ctx context.Context, req appsec.UpdateAttackGroupRequest) (*appsec.UpdateAttackGroupResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.UpdateAttackGroupResponse), args.Error(1)
}

func (m *mockAppsec) UpdateApiRequestConstraints(ctx context.Context, req appsec.UpdateApiRequestConstraintsRequest) (*appsec.UpdateApiRequestConstraintsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.UpdateApiRequestConstraintsResponse), args.Error(1)
}

func (m *mockAppsec) UpdateAdvancedSettingsPrefetch(ctx context.Context, req appsec.UpdateAdvancedSettingsPrefetchRequest) (*appsec.UpdateAdvancedSettingsPrefetchResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.UpdateAdvancedSettingsPrefetchResponse), args.Error(1)
}

func (m *mockAppsec) UpdateAdvancedSettingsPragma(ctx context.Context, req appsec.UpdateAdvancedSettingsPragmaRequest) (*appsec.UpdateAdvancedSettingsPragmaResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.UpdateAdvancedSettingsPragmaResponse), args.Error(1)
}

func (m *mockAppsec) UpdateAdvancedSettingsLogging(ctx context.Context, req appsec.UpdateAdvancedSettingsLoggingRequest) (*appsec.UpdateAdvancedSettingsLoggingResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.UpdateAdvancedSettingsLoggingResponse), args.Error(1)
}

func (m *mockAppsec) UpdateAdvancedSettingsEvasivePathMatch(ctx context.Context, req appsec.UpdateAdvancedSettingsEvasivePathMatchRequest) (*appsec.UpdateAdvancedSettingsEvasivePathMatchResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.UpdateAdvancedSettingsEvasivePathMatchResponse), args.Error(1)
}

func (m *mockAppsec) UpdateAPIConstraintsProtection(ctx context.Context, req appsec.UpdateAPIConstraintsProtectionRequest) (*appsec.UpdateAPIConstraintsProtectionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.UpdateAPIConstraintsProtectionResponse), args.Error(1)
}

func (m *mockAppsec) RemoveSiemSettings(ctx context.Context, req appsec.RemoveSiemSettingsRequest) (*appsec.RemoveSiemSettingsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.RemoveSiemSettingsResponse), args.Error(1)
}

func (m *mockAppsec) RemoveSecurityPolicy(ctx context.Context, req appsec.RemoveSecurityPolicyRequest) (*appsec.RemoveSecurityPolicyResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.RemoveSecurityPolicyResponse), args.Error(1)
}

func (m *mockAppsec) RemoveReputationProtection(ctx context.Context, req appsec.RemoveReputationProtectionRequest) (*appsec.RemoveReputationProtectionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.RemoveReputationProtectionResponse), args.Error(1)
}

func (m *mockAppsec) RemoveReputationProfile(ctx context.Context, req appsec.RemoveReputationProfileRequest) (*appsec.RemoveReputationProfileResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.RemoveReputationProfileResponse), args.Error(1)
}

func (m *mockAppsec) RemoveReputationAnalysis(ctx context.Context, req appsec.RemoveReputationAnalysisRequest) (*appsec.RemoveReputationAnalysisResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.RemoveReputationAnalysisResponse), args.Error(1)
}

func (m *mockAppsec) RemoveRatePolicy(ctx context.Context, req appsec.RemoveRatePolicyRequest) (*appsec.RemoveRatePolicyResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.RemoveRatePolicyResponse), args.Error(1)
}

func (m *mockAppsec) RemovePolicyProtections(ctx context.Context, req appsec.UpdatePolicyProtectionsRequest) (*appsec.PolicyProtectionsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.PolicyProtectionsResponse), args.Error(1)
}

func (m *mockAppsec) RemoveNetworkLayerProtection(ctx context.Context, req appsec.RemoveNetworkLayerProtectionRequest) (*appsec.RemoveNetworkLayerProtectionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.RemoveNetworkLayerProtectionResponse), args.Error(1)
}

func (m *mockAppsec) RemoveMatchTarget(ctx context.Context, req appsec.RemoveMatchTargetRequest) (*appsec.RemoveMatchTargetResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.RemoveMatchTargetResponse), args.Error(1)
}

func (m *mockAppsec) RemoveEval(ctx context.Context, req appsec.RemoveEvalRequest) (*appsec.RemoveEvalResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.RemoveEvalResponse), args.Error(1)
}

func (m *mockAppsec) RemoveEvalHost(ctx context.Context, req appsec.RemoveEvalHostRequest) (*appsec.RemoveEvalHostResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.RemoveEvalHostResponse), args.Error(1)
}

func (m *mockAppsec) RemoveCustomRule(ctx context.Context, req appsec.RemoveCustomRuleRequest) (*appsec.RemoveCustomRuleResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.RemoveCustomRuleResponse), args.Error(1)
}

func (m *mockAppsec) RemoveCustomDeny(ctx context.Context, req appsec.RemoveCustomDenyRequest) (*appsec.RemoveCustomDenyResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.RemoveCustomDenyResponse), args.Error(1)
}

func (m *mockAppsec) RemoveConfigurationVersionClone(ctx context.Context, req appsec.RemoveConfigurationVersionCloneRequest) (*appsec.RemoveConfigurationVersionCloneResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.RemoveConfigurationVersionCloneResponse), args.Error(1)
}

func (m *mockAppsec) RemoveConfiguration(ctx context.Context, req appsec.RemoveConfigurationRequest) (*appsec.RemoveConfigurationResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.RemoveConfigurationResponse), args.Error(1)
}

func (m *mockAppsec) RemoveBypassNetworkLists(ctx context.Context, req appsec.RemoveBypassNetworkListsRequest) (*appsec.RemoveBypassNetworkListsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.RemoveBypassNetworkListsResponse), args.Error(1)
}

func (m *mockAppsec) RemoveApiRequestConstraints(ctx context.Context, req appsec.RemoveApiRequestConstraintsRequest) (*appsec.RemoveApiRequestConstraintsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.RemoveApiRequestConstraintsResponse), args.Error(1)
}

func (m *mockAppsec) RemoveAdvancedSettingsLogging(ctx context.Context, req appsec.RemoveAdvancedSettingsLoggingRequest) (*appsec.RemoveAdvancedSettingsLoggingResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.RemoveAdvancedSettingsLoggingResponse), args.Error(1)
}

func (m *mockAppsec) RemoveAdvancedSettingsEvasivePathMatch(ctx context.Context, req appsec.RemoveAdvancedSettingsEvasivePathMatchRequest) (*appsec.RemoveAdvancedSettingsEvasivePathMatchResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.RemoveAdvancedSettingsEvasivePathMatchResponse), args.Error(1)
}

func (m *mockAppsec) RemoveActivations(ctx context.Context, req appsec.RemoveActivationsRequest) (*appsec.RemoveActivationsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.RemoveActivationsResponse), args.Error(1)
}

func (m *mockAppsec) GetWAPSelectedHostnames(ctx context.Context, req appsec.GetWAPSelectedHostnamesRequest) (*appsec.GetWAPSelectedHostnamesResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetWAPSelectedHostnamesResponse), args.Error(1)
}

func (m *mockAppsec) GetWAFProtections(ctx context.Context, req appsec.GetWAFProtectionsRequest) (*appsec.GetWAFProtectionsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetWAFProtectionsResponse), args.Error(1)
}

func (m *mockAppsec) GetWAFProtection(ctx context.Context, req appsec.GetWAFProtectionRequest) (*appsec.GetWAFProtectionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetWAFProtectionResponse), args.Error(1)
}

func (m *mockAppsec) GetWAFModes(ctx context.Context, req appsec.GetWAFModesRequest) (*appsec.GetWAFModesResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetWAFModesResponse), args.Error(1)
}

func (m *mockAppsec) GetWAFMode(ctx context.Context, req appsec.GetWAFModeRequest) (*appsec.GetWAFModeResponse, error) {

	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetWAFModeResponse), args.Error(1)
}

func (m *mockAppsec) GetVersionNotes(ctx context.Context, req appsec.GetVersionNotesRequest) (*appsec.GetVersionNotesResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetVersionNotesResponse), args.Error(1)
}

func (m *mockAppsec) GetTuningRecommendations(ctx context.Context, req appsec.GetTuningRecommendationsRequest) (*appsec.GetTuningRecommendationsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetTuningRecommendationsResponse), args.Error(1)
}

func (m *mockAppsec) GetThreatIntel(ctx context.Context, req appsec.GetThreatIntelRequest) (*appsec.GetThreatIntelResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetThreatIntelResponse), args.Error(1)
}

func (m *mockAppsec) GetSlowPostProtections(ctx context.Context, req appsec.GetSlowPostProtectionsRequest) (*appsec.GetSlowPostProtectionsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetSlowPostProtectionsResponse), args.Error(1)
}

func (m *mockAppsec) GetSlowPostProtectionSettings(ctx context.Context, req appsec.GetSlowPostProtectionSettingsRequest) (*appsec.GetSlowPostProtectionSettingsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetSlowPostProtectionSettingsResponse), args.Error(1)
}

func (m *mockAppsec) GetSlowPostProtectionSetting(ctx context.Context, req appsec.GetSlowPostProtectionSettingRequest) (*appsec.GetSlowPostProtectionSettingResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetSlowPostProtectionSettingResponse), args.Error(1)
}

func (m *mockAppsec) GetSlowPostProtection(ctx context.Context, req appsec.GetSlowPostProtectionRequest) (*appsec.GetSlowPostProtectionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetSlowPostProtectionResponse), args.Error(1)
}

func (m *mockAppsec) GetSiemSettings(ctx context.Context, req appsec.GetSiemSettingsRequest) (*appsec.GetSiemSettingsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetSiemSettingsResponse), args.Error(1)
}

func (m *mockAppsec) GetSiemDefinitions(ctx context.Context, req appsec.GetSiemDefinitionsRequest) (*appsec.GetSiemDefinitionsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetSiemDefinitionsResponse), args.Error(1)
}

func (m *mockAppsec) GetSelectedHostnames(ctx context.Context, req appsec.GetSelectedHostnamesRequest) (*appsec.GetSelectedHostnamesResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetSelectedHostnamesResponse), args.Error(1)
}

func (m *mockAppsec) GetSelectedHostname(ctx context.Context, req appsec.GetSelectedHostnameRequest) (*appsec.GetSelectedHostnameResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetSelectedHostnameResponse), args.Error(1)
}

func (m *mockAppsec) GetSelectableHostnames(ctx context.Context, req appsec.GetSelectableHostnamesRequest) (*appsec.GetSelectableHostnamesResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetSelectableHostnamesResponse), args.Error(1)
}

func (m *mockAppsec) GetSecurityPolicyClones(ctx context.Context, req appsec.GetSecurityPolicyClonesRequest) (*appsec.GetSecurityPolicyClonesResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetSecurityPolicyClonesResponse), args.Error(1)
}

func (m *mockAppsec) GetSecurityPolicyClone(ctx context.Context, req appsec.GetSecurityPolicyCloneRequest) (*appsec.GetSecurityPolicyCloneResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetSecurityPolicyCloneResponse), args.Error(1)
}

func (m *mockAppsec) GetSecurityPolicy(ctx context.Context, req appsec.GetSecurityPolicyRequest) (*appsec.GetSecurityPolicyResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetSecurityPolicyResponse), args.Error(1)
}

func (m *mockAppsec) GetRuleUpgrade(ctx context.Context, req appsec.GetRuleUpgradeRequest) (*appsec.GetRuleUpgradeResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetRuleUpgradeResponse), args.Error(1)
}

func (m *mockAppsec) GetRules(ctx context.Context, req appsec.GetRulesRequest) (*appsec.GetRulesResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetRulesResponse), args.Error(1)
}

func (m *mockAppsec) GetRule(ctx context.Context, req appsec.GetRuleRequest) (*appsec.GetRuleResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetRuleResponse), args.Error(1)
}

func (m *mockAppsec) GetReputationProtections(ctx context.Context, req appsec.GetReputationProtectionsRequest) (*appsec.GetReputationProtectionsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetReputationProtectionsResponse), args.Error(1)
}

func (m *mockAppsec) GetReputationProtection(ctx context.Context, req appsec.GetReputationProtectionRequest) (*appsec.GetReputationProtectionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetReputationProtectionResponse), args.Error(1)
}

func (m *mockAppsec) GetReputationProfiles(ctx context.Context, req appsec.GetReputationProfilesRequest) (*appsec.GetReputationProfilesResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetReputationProfilesResponse), args.Error(1)
}

func (m *mockAppsec) GetReputationProfileActions(ctx context.Context, req appsec.GetReputationProfileActionsRequest) (*appsec.GetReputationProfileActionsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetReputationProfileActionsResponse), args.Error(1)
}

func (m *mockAppsec) GetReputationProfileAction(ctx context.Context, req appsec.GetReputationProfileActionRequest) (*appsec.GetReputationProfileActionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetReputationProfileActionResponse), args.Error(1)
}

func (m *mockAppsec) GetReputationProfile(ctx context.Context, req appsec.GetReputationProfileRequest) (*appsec.GetReputationProfileResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetReputationProfileResponse), args.Error(1)
}

func (m *mockAppsec) GetReputationAnalysis(ctx context.Context, req appsec.GetReputationAnalysisRequest) (*appsec.GetReputationAnalysisResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetReputationAnalysisResponse), args.Error(1)
}

func (m *mockAppsec) GetRateProtections(ctx context.Context, req appsec.GetRateProtectionsRequest) (*appsec.GetRateProtectionsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetRateProtectionsResponse), args.Error(1)
}

func (m *mockAppsec) GetRateProtection(ctx context.Context, req appsec.GetRateProtectionRequest) (*appsec.GetRateProtectionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetRateProtectionResponse), args.Error(1)
}

func (m *mockAppsec) GetRatePolicyActions(ctx context.Context, req appsec.GetRatePolicyActionsRequest) (*appsec.GetRatePolicyActionsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetRatePolicyActionsResponse), args.Error(1)
}

func (m *mockAppsec) GetRatePolicyAction(ctx context.Context, req appsec.GetRatePolicyActionRequest) (*appsec.GetRatePolicyActionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetRatePolicyActionResponse), args.Error(1)
}

func (m *mockAppsec) GetRatePolicy(ctx context.Context, req appsec.GetRatePolicyRequest) (*appsec.GetRatePolicyResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetRatePolicyResponse), args.Error(1)
}

func (m *mockAppsec) GetRatePolicies(ctx context.Context, req appsec.GetRatePoliciesRequest) (*appsec.GetRatePoliciesResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetRatePoliciesResponse), args.Error(1)
}

func (m *mockAppsec) GetSecurityPolicies(ctx context.Context, req appsec.GetSecurityPoliciesRequest) (*appsec.GetSecurityPoliciesResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetSecurityPoliciesResponse), args.Error(1)
}

func (m *mockAppsec) GetPolicyProtections(ctx context.Context, req appsec.GetPolicyProtectionsRequest) (*appsec.PolicyProtectionsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.PolicyProtectionsResponse), args.Error(1)
}

func (m *mockAppsec) GetPenaltyBoxes(ctx context.Context, req appsec.GetPenaltyBoxesRequest) (*appsec.GetPenaltyBoxesResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetPenaltyBoxesResponse), args.Error(1)
}

func (m *mockAppsec) GetPenaltyBox(ctx context.Context, req appsec.GetPenaltyBoxRequest) (*appsec.GetPenaltyBoxResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetPenaltyBoxResponse), args.Error(1)
}

func (m *mockAppsec) GetEvalPenaltyBox(ctx context.Context, params appsec.GetPenaltyBoxRequest) (*appsec.GetPenaltyBoxResponse, error) {
	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetPenaltyBoxResponse), args.Error(1)
}

func (m *mockAppsec) UpdateEvalPenaltyBox(ctx context.Context, params appsec.UpdatePenaltyBoxRequest) (*appsec.UpdatePenaltyBoxResponse, error) {
	args := m.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.UpdatePenaltyBoxResponse), args.Error(1)
}

func (m *mockAppsec) GetNetworkLayerProtections(ctx context.Context, req appsec.GetNetworkLayerProtectionsRequest) (*appsec.GetNetworkLayerProtectionsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetNetworkLayerProtectionsResponse), args.Error(1)
}

func (m *mockAppsec) GetNetworkLayerProtection(ctx context.Context, req appsec.GetNetworkLayerProtectionRequest) (*appsec.GetNetworkLayerProtectionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetNetworkLayerProtectionResponse), args.Error(1)
}

func (m *mockAppsec) GetMatchTargets(ctx context.Context, req appsec.GetMatchTargetsRequest) (*appsec.GetMatchTargetsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetMatchTargetsResponse), args.Error(1)
}

func (m *mockAppsec) GetMatchTargetSequence(ctx context.Context, req appsec.GetMatchTargetSequenceRequest) (*appsec.GetMatchTargetSequenceResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetMatchTargetSequenceResponse), args.Error(1)
}

func (m *mockAppsec) GetMatchTarget(ctx context.Context, req appsec.GetMatchTargetRequest) (*appsec.GetMatchTargetResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetMatchTargetResponse), args.Error(1)
}

func (m *mockAppsec) GetIPGeoProtection(ctx context.Context, req appsec.GetIPGeoProtectionRequest) (*appsec.GetIPGeoProtectionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetIPGeoProtectionResponse), args.Error(1)
}

func (m *mockAppsec) GetIPGeoProtections(ctx context.Context, req appsec.GetIPGeoProtectionsRequest) (*appsec.GetIPGeoProtectionsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetIPGeoProtectionsResponse), args.Error(1)
}

func (m *mockAppsec) GetIPGeo(ctx context.Context, req appsec.GetIPGeoRequest) (*appsec.GetIPGeoResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetIPGeoResponse), args.Error(1)
}

func (m *mockAppsec) GetFailoverHostnames(ctx context.Context, req appsec.GetFailoverHostnamesRequest) (*appsec.GetFailoverHostnamesResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetFailoverHostnamesResponse), args.Error(1)
}

func (m *mockAppsec) GetExportConfiguration(ctx context.Context, req appsec.GetExportConfigurationRequest) (*appsec.GetExportConfigurationResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetExportConfigurationResponse), args.Error(1)
}

func (m *mockAppsec) GetExportConfigurations(ctx context.Context, req appsec.GetExportConfigurationsRequest) (*appsec.GetExportConfigurationsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetExportConfigurationsResponse), args.Error(1)
}

func (m *mockAppsec) GetEvals(ctx context.Context, req appsec.GetEvalsRequest) (*appsec.GetEvalsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetEvalsResponse), args.Error(1)
}

func (m *mockAppsec) GetEvalRules(ctx context.Context, req appsec.GetEvalRulesRequest) (*appsec.GetEvalRulesResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetEvalRulesResponse), args.Error(1)
}

func (m *mockAppsec) GetEvalRule(ctx context.Context, req appsec.GetEvalRuleRequest) (*appsec.GetEvalRuleResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetEvalRuleResponse), args.Error(1)
}

func (m *mockAppsec) GetEvalProtectHosts(ctx context.Context, req appsec.GetEvalProtectHostsRequest) (*appsec.GetEvalProtectHostsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetEvalProtectHostsResponse), args.Error(1)
}

func (m *mockAppsec) GetEvalProtectHost(ctx context.Context, req appsec.GetEvalProtectHostRequest) (*appsec.GetEvalProtectHostResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetEvalProtectHostResponse), args.Error(1)
}

func (m *mockAppsec) GetEvalHosts(ctx context.Context, req appsec.GetEvalHostsRequest) (*appsec.GetEvalHostsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetEvalHostsResponse), args.Error(1)
}

func (m *mockAppsec) GetEvalHost(ctx context.Context, req appsec.GetEvalHostRequest) (*appsec.GetEvalHostResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetEvalHostResponse), args.Error(1)
}

func (m *mockAppsec) GetEvalGroups(ctx context.Context, req appsec.GetAttackGroupsRequest) (*appsec.GetAttackGroupsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetAttackGroupsResponse), args.Error(1)
}

func (m *mockAppsec) GetEvalGroup(ctx context.Context, req appsec.GetAttackGroupRequest) (*appsec.GetAttackGroupResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetAttackGroupResponse), args.Error(1)
}

func (m *mockAppsec) GetEval(ctx context.Context, req appsec.GetEvalRequest) (*appsec.GetEvalResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetEvalResponse), args.Error(1)
}

func (m *mockAppsec) GetCustomRules(ctx context.Context, req appsec.GetCustomRulesRequest) (*appsec.GetCustomRulesResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetCustomRulesResponse), args.Error(1)
}

func (m *mockAppsec) GetCustomRuleActions(ctx context.Context, req appsec.GetCustomRuleActionsRequest) (*appsec.GetCustomRuleActionsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetCustomRuleActionsResponse), args.Error(1)
}

func (m *mockAppsec) GetCustomRuleAction(ctx context.Context, req appsec.GetCustomRuleActionRequest) (*appsec.GetCustomRuleActionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetCustomRuleActionResponse), args.Error(1)
}

func (m *mockAppsec) GetCustomRule(ctx context.Context, req appsec.GetCustomRuleRequest) (*appsec.GetCustomRuleResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetCustomRuleResponse), args.Error(1)
}

func (m *mockAppsec) GetCustomDenyList(ctx context.Context, req appsec.GetCustomDenyListRequest) (*appsec.GetCustomDenyListResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetCustomDenyListResponse), args.Error(1)
}

func (m *mockAppsec) GetCustomDeny(ctx context.Context, req appsec.GetCustomDenyRequest) (*appsec.GetCustomDenyResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetCustomDenyResponse), args.Error(1)
}

func (m *mockAppsec) GetContractsGroups(ctx context.Context, req appsec.GetContractsGroupsRequest) (*appsec.GetContractsGroupsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetContractsGroupsResponse), args.Error(1)
}

func (m *mockAppsec) GetConfigurations(ctx context.Context, req appsec.GetConfigurationsRequest) (*appsec.GetConfigurationsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetConfigurationsResponse), args.Error(1)
}

func (m *mockAppsec) GetConfigurationVersions(ctx context.Context, req appsec.GetConfigurationVersionsRequest) (*appsec.GetConfigurationVersionsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetConfigurationVersionsResponse), args.Error(1)
}

func (m *mockAppsec) GetConfigurationVersionClone(ctx context.Context, req appsec.GetConfigurationVersionCloneRequest) (*appsec.GetConfigurationVersionCloneResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetConfigurationVersionCloneResponse), args.Error(1)
}

func (m *mockAppsec) GetConfigurationClone(ctx context.Context, req appsec.GetConfigurationCloneRequest) (*appsec.GetConfigurationCloneResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetConfigurationCloneResponse), args.Error(1)
}

func (m *mockAppsec) GetBypassNetworkLists(ctx context.Context, req appsec.GetBypassNetworkListsRequest) (*appsec.GetBypassNetworkListsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetBypassNetworkListsResponse), args.Error(1)
}

func (m *mockAppsec) GetAttackGroups(ctx context.Context, req appsec.GetAttackGroupsRequest) (*appsec.GetAttackGroupsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetAttackGroupsResponse), args.Error(1)
}

func (m *mockAppsec) GetAttackGroupRecommendations(ctx context.Context, req appsec.GetAttackGroupRecommendationsRequest) (*appsec.GetAttackGroupRecommendationsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetAttackGroupRecommendationsResponse), args.Error(1)
}

func (m *mockAppsec) GetAttackGroup(ctx context.Context, req appsec.GetAttackGroupRequest) (*appsec.GetAttackGroupResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetAttackGroupResponse), args.Error(1)
}

func (m *mockAppsec) GetApiRequestConstraints(ctx context.Context, req appsec.GetApiRequestConstraintsRequest) (*appsec.GetApiRequestConstraintsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetApiRequestConstraintsResponse), args.Error(1)
}

func (m *mockAppsec) GetApiHostnameCoverageOverlapping(ctx context.Context, req appsec.GetApiHostnameCoverageOverlappingRequest) (*appsec.GetApiHostnameCoverageOverlappingResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetApiHostnameCoverageOverlappingResponse), args.Error(1)
}

func (m *mockAppsec) GetApiHostnameCoverageMatchTargets(ctx context.Context, req appsec.GetApiHostnameCoverageMatchTargetsRequest) (*appsec.GetApiHostnameCoverageMatchTargetsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetApiHostnameCoverageMatchTargetsResponse), args.Error(1)
}

func (m *mockAppsec) GetApiHostnameCoverage(ctx context.Context, req appsec.GetApiHostnameCoverageRequest) (*appsec.GetApiHostnameCoverageResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetApiHostnameCoverageResponse), args.Error(1)
}

func (m *mockAppsec) GetApiEndpoints(ctx context.Context, req appsec.GetApiEndpointsRequest) (*appsec.GetApiEndpointsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetApiEndpointsResponse), args.Error(1)
}

func (m *mockAppsec) GetAdvancedSettingsPragma(ctx context.Context, req appsec.GetAdvancedSettingsPragmaRequest) (*appsec.GetAdvancedSettingsPragmaResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetAdvancedSettingsPragmaResponse), args.Error(1)
}

func (m *mockAppsec) GetAdvancedSettingsPrefetch(ctx context.Context, req appsec.GetAdvancedSettingsPrefetchRequest) (*appsec.GetAdvancedSettingsPrefetchResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetAdvancedSettingsPrefetchResponse), args.Error(1)
}

func (m *mockAppsec) GetAdvancedSettingsLogging(ctx context.Context, req appsec.GetAdvancedSettingsLoggingRequest) (*appsec.GetAdvancedSettingsLoggingResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetAdvancedSettingsLoggingResponse), args.Error(1)
}

func (m *mockAppsec) GetAdvancedSettingsEvasivePathMatch(ctx context.Context, req appsec.GetAdvancedSettingsEvasivePathMatchRequest) (*appsec.GetAdvancedSettingsEvasivePathMatchResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetAdvancedSettingsEvasivePathMatchResponse), args.Error(1)
}

func (m *mockAppsec) GetActivations(ctx context.Context, req appsec.GetActivationsRequest) (*appsec.GetActivationsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetActivationsResponse), args.Error(1)
}

/*
func (m *mockAppsec) GetActivation(ctx context.Context, req appsec.GetActivationRequest) (*appsec.GetActivationResponse, error) {
        args := m.Called(ctx, req)
        if args.Get(0) == nil {
                return nil, args.Error(1)
        }
        return args.Get(0).(*appsec.GetActivationResponse), args.Error(1)
}
*/

func (m *mockAppsec) GetAPIConstraintsProtection(ctx context.Context, req appsec.GetAPIConstraintsProtectionRequest) (*appsec.GetAPIConstraintsProtectionResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetAPIConstraintsProtectionResponse), args.Error(1)
}

func (m *mockAppsec) CreateSecurityPolicyClone(ctx context.Context, req appsec.CreateSecurityPolicyCloneRequest) (*appsec.CreateSecurityPolicyCloneResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.CreateSecurityPolicyCloneResponse), args.Error(1)
}

func (m *mockAppsec) CreateSecurityPolicy(ctx context.Context, req appsec.CreateSecurityPolicyRequest) (*appsec.CreateSecurityPolicyResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.CreateSecurityPolicyResponse), args.Error(1)
}

func (m *mockAppsec) CreateReputationProfile(ctx context.Context, req appsec.CreateReputationProfileRequest) (*appsec.CreateReputationProfileResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.CreateReputationProfileResponse), args.Error(1)
}

func (m *mockAppsec) CreateRatePolicy(ctx context.Context, req appsec.CreateRatePolicyRequest) (*appsec.CreateRatePolicyResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.CreateRatePolicyResponse), args.Error(1)
}

func (m *mockAppsec) CreateMatchTarget(ctx context.Context, req appsec.CreateMatchTargetRequest) (*appsec.CreateMatchTargetResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.CreateMatchTargetResponse), args.Error(1)
}

func (m *mockAppsec) CreateCustomRule(ctx context.Context, req appsec.CreateCustomRuleRequest) (*appsec.CreateCustomRuleResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.CreateCustomRuleResponse), args.Error(1)
}

func (m *mockAppsec) CreateCustomDeny(ctx context.Context, req appsec.CreateCustomDenyRequest) (*appsec.CreateCustomDenyResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.CreateCustomDenyResponse), args.Error(1)
}

func (m *mockAppsec) CreateConfiguration(ctx context.Context, req appsec.CreateConfigurationRequest) (*appsec.CreateConfigurationResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.CreateConfigurationResponse), args.Error(1)
}

func (m *mockAppsec) CreateConfigurationVersionClone(ctx context.Context, req appsec.CreateConfigurationVersionCloneRequest) (*appsec.CreateConfigurationVersionCloneResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.CreateConfigurationVersionCloneResponse), args.Error(1)
}

func (m *mockAppsec) CreateConfigurationClone(ctx context.Context, req appsec.CreateConfigurationCloneRequest) (*appsec.CreateConfigurationCloneResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.CreateConfigurationCloneResponse), args.Error(1)
}

func (m *mockAppsec) CreateActivations(ctx context.Context, req appsec.CreateActivationsRequest, _ bool) (*appsec.CreateActivationsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.CreateActivationsResponse), args.Error(1)
}

func (m *mockAppsec) GetConfiguration(ctx context.Context, req appsec.GetConfigurationRequest) (*appsec.GetConfigurationResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetConfigurationResponse), args.Error(1)
}

func (m *mockAppsec) GetWAPBypassNetworkLists(ctx context.Context, req appsec.GetWAPBypassNetworkListsRequest) (*appsec.GetWAPBypassNetworkListsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetWAPBypassNetworkListsResponse), args.Error(1)
}

func (m *mockAppsec) RemoveWAPBypassNetworkLists(ctx context.Context, req appsec.RemoveWAPBypassNetworkListsRequest) (*appsec.RemoveWAPBypassNetworkListsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.RemoveWAPBypassNetworkListsResponse), args.Error(1)
}

func (m *mockAppsec) UpdateWAPBypassNetworkLists(ctx context.Context, req appsec.UpdateWAPBypassNetworkListsRequest) (*appsec.UpdateWAPBypassNetworkListsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.UpdateWAPBypassNetworkListsResponse), args.Error(1)
}

func (m *mockAppsec) GetActivationHistory(ctx context.Context, req appsec.GetActivationHistoryRequest) (*appsec.GetActivationHistoryResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*appsec.GetActivationHistoryResponse), args.Error(1)
}

func (p *mockAppsec) GetMalwareProtection(ctx context.Context, params appsec.GetMalwareProtectionRequest) (*appsec.GetMalwareProtectionResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetMalwareProtectionResponse), args.Error(1)
}

func (p *mockAppsec) GetMalwareProtections(ctx context.Context, params appsec.GetMalwareProtectionsRequest) (*appsec.GetMalwareProtectionsResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetMalwareProtectionsResponse), args.Error(1)
}

func (p *mockAppsec) UpdateMalwareProtection(ctx context.Context, params appsec.UpdateMalwareProtectionRequest) (*appsec.UpdateMalwareProtectionResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.UpdateMalwareProtectionResponse), args.Error(1)
}

func (p *mockAppsec) GetMalwareContentTypes(ctx context.Context, params appsec.GetMalwareContentTypesRequest) (*appsec.GetMalwareContentTypesResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetMalwareContentTypesResponse), args.Error(1)
}

func (p *mockAppsec) CreateMalwarePolicy(ctx context.Context, params appsec.CreateMalwarePolicyRequest) (*appsec.MalwarePolicyResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.MalwarePolicyResponse), args.Error(1)
}

func (p *mockAppsec) GetMalwarePolicy(ctx context.Context, params appsec.GetMalwarePolicyRequest) (*appsec.MalwarePolicyResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.MalwarePolicyResponse), args.Error(1)
}

func (p *mockAppsec) GetMalwarePolicies(ctx context.Context, params appsec.GetMalwarePoliciesRequest) (*appsec.MalwarePoliciesResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.MalwarePoliciesResponse), args.Error(1)
}

func (p *mockAppsec) UpdateMalwarePolicy(ctx context.Context, params appsec.UpdateMalwarePolicyRequest) (*appsec.MalwarePolicyResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.MalwarePolicyResponse), args.Error(1)
}

func (p *mockAppsec) RemoveMalwarePolicy(_ context.Context, _ appsec.RemoveMalwarePolicyRequest) error {
	return nil
}

func (p *mockAppsec) GetMalwarePolicyActions(ctx context.Context, params appsec.GetMalwarePolicyActionsRequest) (*appsec.GetMalwarePolicyActionsResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.GetMalwarePolicyActionsResponse), args.Error(1)
}

func (p *mockAppsec) UpdateMalwarePolicyAction(ctx context.Context, params appsec.UpdateMalwarePolicyActionRequest) (*appsec.UpdateMalwarePolicyActionResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.UpdateMalwarePolicyActionResponse), args.Error(1)
}

func (p *mockAppsec) UpdateMalwarePolicyActions(ctx context.Context, params appsec.UpdateMalwarePolicyActionsRequest) (*appsec.UpdateMalwarePolicyActionsResponse, error) {
	args := p.Called(ctx, params)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*appsec.UpdateMalwarePolicyActionsResponse), args.Error(1)
}
