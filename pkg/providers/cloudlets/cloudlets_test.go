package cloudlets

import (
	"context"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/cloudlets"
	"github.com/stretchr/testify/mock"
)

type mockCloudlets struct {
	mock.Mock
}

func (m *mockCloudlets) CreateLoadBalancerVersion(ctx context.Context, req cloudlets.CreateLoadBalancerVersionRequest) (*cloudlets.LoadBalancerVersion, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cloudlets.LoadBalancerVersion), args.Error(1)
}

func (m *mockCloudlets) GetLoadBalancerVersion(ctx context.Context, req cloudlets.GetLoadBalancerVersionRequest) (*cloudlets.LoadBalancerVersion, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cloudlets.LoadBalancerVersion), args.Error(1)
}

func (m *mockCloudlets) UpdateLoadBalancerVersion(ctx context.Context, req cloudlets.UpdateLoadBalancerVersionRequest) (*cloudlets.LoadBalancerVersion, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cloudlets.LoadBalancerVersion), args.Error(1)
}

func (m *mockCloudlets) ListLoadBalancerActivations(ctx context.Context, req cloudlets.ListLoadBalancerActivationsRequest) ([]cloudlets.LoadBalancerActivation, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]cloudlets.LoadBalancerActivation), args.Error(1)
}

func (m *mockCloudlets) ActivateLoadBalancerVersion(ctx context.Context, req cloudlets.ActivateLoadBalancerVersionRequest) (*cloudlets.LoadBalancerActivation, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cloudlets.LoadBalancerActivation), args.Error(1)
}

func (m *mockCloudlets) ListPolicyActivations(ctx context.Context, req cloudlets.ListPolicyActivationsRequest) ([]cloudlets.PolicyActivation, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]cloudlets.PolicyActivation), args.Error(1)
}

func (m *mockCloudlets) ActivatePolicyVersion(ctx context.Context, req cloudlets.ActivatePolicyVersionRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *mockCloudlets) ListOrigins(ctx context.Context, req cloudlets.ListOriginsRequest) ([]cloudlets.OriginResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]cloudlets.OriginResponse), args.Error(1)
}

func (m *mockCloudlets) GetOrigin(ctx context.Context, req cloudlets.GetOriginRequest) (*cloudlets.Origin, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cloudlets.Origin), args.Error(1)
}

func (m *mockCloudlets) CreateOrigin(ctx context.Context, req cloudlets.CreateOriginRequest) (*cloudlets.Origin, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cloudlets.Origin), args.Error(1)
}

func (m *mockCloudlets) UpdateOrigin(ctx context.Context, req cloudlets.UpdateOriginRequest) (*cloudlets.Origin, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cloudlets.Origin), args.Error(1)
}

func (m *mockCloudlets) ListPolicies(ctx context.Context, request cloudlets.ListPoliciesRequest) ([]cloudlets.Policy, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]cloudlets.Policy), args.Error(1)
}

func (m *mockCloudlets) GetPolicy(ctx context.Context, policyID cloudlets.GetPolicyRequest) (*cloudlets.Policy, error) {
	args := m.Called(ctx, policyID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cloudlets.Policy), args.Error(1)
}

func (m *mockCloudlets) CreatePolicy(ctx context.Context, req cloudlets.CreatePolicyRequest) (*cloudlets.Policy, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cloudlets.Policy), args.Error(1)
}

func (m *mockCloudlets) RemovePolicy(ctx context.Context, policyID cloudlets.RemovePolicyRequest) error {
	args := m.Called(ctx, policyID)
	return args.Error(0)
}

func (m *mockCloudlets) UpdatePolicy(ctx context.Context, req cloudlets.UpdatePolicyRequest) (*cloudlets.Policy, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cloudlets.Policy), args.Error(1)
}

func (m *mockCloudlets) ListPolicyVersions(ctx context.Context, request cloudlets.ListPolicyVersionsRequest) ([]cloudlets.PolicyVersion, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]cloudlets.PolicyVersion), args.Error(1)
}

func (m *mockCloudlets) GetPolicyVersion(ctx context.Context, req cloudlets.GetPolicyVersionRequest) (*cloudlets.PolicyVersion, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cloudlets.PolicyVersion), args.Error(1)
}

func (m *mockCloudlets) CreatePolicyVersion(ctx context.Context, req cloudlets.CreatePolicyVersionRequest) (*cloudlets.PolicyVersion, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cloudlets.PolicyVersion), args.Error(1)
}

func (m *mockCloudlets) DeletePolicyVersion(ctx context.Context, req cloudlets.DeletePolicyVersionRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *mockCloudlets) UpdatePolicyVersion(ctx context.Context, req cloudlets.UpdatePolicyVersionRequest) (*cloudlets.PolicyVersion, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*cloudlets.PolicyVersion), args.Error(1)
}

func (m *mockCloudlets) GetPolicyProperties(ctx context.Context, req cloudlets.GetPolicyPropertiesRequest) (map[string]cloudlets.PolicyProperty, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(map[string]cloudlets.PolicyProperty), args.Error(1)
}

func (m *mockCloudlets) ListLoadBalancerVersions(ctx context.Context, req cloudlets.ListLoadBalancerVersionsRequest) ([]cloudlets.LoadBalancerVersion, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]cloudlets.LoadBalancerVersion), args.Error(1)
}

func (m *mockCloudlets) DeletePolicyProperty(ctx context.Context, req cloudlets.DeletePolicyPropertyRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}
