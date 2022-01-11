package edgeworkers

import (
	"context"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/edgeworkers"
	"github.com/stretchr/testify/mock"
)

type mockEdgeworkers struct {
	mock.Mock
}

// Activations

func (m *mockEdgeworkers) ListActivations(ctx context.Context, req edgeworkers.ListActivationsRequest) (*edgeworkers.ListActivationsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*edgeworkers.ListActivationsResponse), args.Error(1)
}

func (m *mockEdgeworkers) GetActivation(ctx context.Context, req edgeworkers.GetActivationRequest) (*edgeworkers.Activation, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*edgeworkers.Activation), args.Error(1)
}


func (m *mockEdgeworkers) ActivateVersion(ctx context.Context, req edgeworkers.ActivateVersionRequest) (*edgeworkers.Activation, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*edgeworkers.Activation), args.Error(1)
}

func (m *mockEdgeworkers) CancelPendingActivation(ctx context.Context, req edgeworkers.CancelActivationRequest) (*edgeworkers.Activation, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*edgeworkers.Activation), args.Error(1)
}

// Contracts

func (m *mockEdgeworkers) ListContracts(ctx context.Context) (*edgeworkers.ListContractsResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*edgeworkers.ListContractsResponse), args.Error(1)
}

// Deactivations

func (m *mockEdgeworkers) ListDeactivations(ctx context.Context, req edgeworkers.ListDeactivationsRequest) (*edgeworkers.ListDeactivationsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*edgeworkers.ListDeactivationsResponse), args.Error(1)
}

func (m *mockEdgeworkers) GetDeactivation(ctx context.Context, req edgeworkers.GetDeactivationRequest) (*edgeworkers.Deactivation, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*edgeworkers.Deactivation), args.Error(1)
}

func (m *mockEdgeworkers) DeactivateVersion(ctx context.Context, req edgeworkers.DeactivateVersionRequest) (*edgeworkers.Deactivation, error)  {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*edgeworkers.Deactivation), args.Error(1)
}

// EdgeKVAccessTokens

func (m *mockEdgeworkers) CreateEdgeKVAccessToken(ctx context.Context, req edgeworkers.CreateEdgeKVAccessTokenRequest) (*edgeworkers.CreateEdgeKVAccessTokenResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*edgeworkers.CreateEdgeKVAccessTokenResponse), args.Error(1)
}

func (m *mockEdgeworkers) GetEdgeKVAccessToken(ctx context.Context, req edgeworkers.GetEdgeKVAccessTokenRequest) (*edgeworkers.GetEdgeKVAccessTokenResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*edgeworkers.GetEdgeKVAccessTokenResponse), args.Error(1)
}

func (m *mockEdgeworkers) ListEdgeKVAccessTokens(ctx context.Context, req edgeworkers.ListEdgeKVAccessTokensRequest) (*edgeworkers.ListEdgeKVAccessTokensResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*edgeworkers.ListEdgeKVAccessTokensResponse), args.Error(1)
}

func (m *mockEdgeworkers) DeleteEdgeKVAccessToken(ctx context.Context, req edgeworkers.DeleteEdgeKVAccessTokenRequest) (*edgeworkers.DeleteEdgeKVAccessTokenResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*edgeworkers.DeleteEdgeKVAccessTokenResponse), args.Error(1)
}

// EdgeKVInitialize

func (m *mockEdgeworkers) InitializeEdgeKV(ctx context.Context) (*edgeworkers.EdgeKVInitializationStatus, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*edgeworkers.EdgeKVInitializationStatus), args.Error(1)
}

func (m *mockEdgeworkers) GetEdgeKVInitializationStatus(ctx context.Context) (*edgeworkers.EdgeKVInitializationStatus, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*edgeworkers.EdgeKVInitializationStatus), args.Error(1)
}

// EdgeKVItems

func (m *mockEdgeworkers) ListItems(ctx context.Context, req edgeworkers.ListItemsRequest) (*edgeworkers.ListItemsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*edgeworkers.ListItemsResponse), args.Error(1)
}

func (m *mockEdgeworkers) GetItem(ctx context.Context, req edgeworkers.GetItemRequest) (*edgeworkers.Item, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*edgeworkers.Item), args.Error(1)
}

func (m *mockEdgeworkers) UpsertItem(ctx context.Context, req edgeworkers.UpsertItemRequest) (*string, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*string), args.Error(1)
}

func (m *mockEdgeworkers) DeleteItem(ctx context.Context, req edgeworkers.DeleteItemRequest) (*string, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*string), args.Error(1)
}

// EdgeKVNamespaces

func (m *mockEdgeworkers) ListEdgeKVNamespaces(ctx context.Context, req edgeworkers.ListEdgeKVNamespacesRequest) (*edgeworkers.ListEdgeKVNamespacesResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*edgeworkers.ListEdgeKVNamespacesResponse), args.Error(1)
}

func (m *mockEdgeworkers) GetEdgeKVNamespace(ctx context.Context, req edgeworkers.GetEdgeKVNamespaceRequest) (*edgeworkers.Namespace, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*edgeworkers.Namespace), args.Error(1)
}

func (m *mockEdgeworkers) CreateEdgeKVNamespace(ctx context.Context, req edgeworkers.CreateEdgeKVNamespaceRequest) (*edgeworkers.Namespace, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*edgeworkers.Namespace), args.Error(1)
}

func (m *mockEdgeworkers) UpdateEdgeKVNamespace(ctx context.Context, req edgeworkers.UpdateEdgeKVNamespaceRequest) (*edgeworkers.Namespace, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*edgeworkers.Namespace), args.Error(1)
}

// EdgeWorkerID

func (m *mockEdgeworkers) GetEdgeWorkerID(ctx context.Context, req edgeworkers.GetEdgeWorkerIDRequest) (*edgeworkers.EdgeWorkerID, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*edgeworkers.EdgeWorkerID), args.Error(1)
}

func (m *mockEdgeworkers) ListEdgeWorkersID(ctx context.Context, req edgeworkers.ListEdgeWorkersIDRequest) (*edgeworkers.ListEdgeWorkersIDResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*edgeworkers.ListEdgeWorkersIDResponse), args.Error(1)
}

func (m *mockEdgeworkers) CreateEdgeWorkerID(ctx context.Context, req edgeworkers.CreateEdgeWorkerIDRequest) (*edgeworkers.EdgeWorkerID, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*edgeworkers.EdgeWorkerID), args.Error(1)
}

func (m *mockEdgeworkers) UpdateEdgeWorkerID(ctx context.Context, req edgeworkers.UpdateEdgeWorkerIDRequest) (*edgeworkers.EdgeWorkerID, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*edgeworkers.EdgeWorkerID), args.Error(1)
}

func (m *mockEdgeworkers) CloneEdgeWorkerID(ctx context.Context, req edgeworkers.CloneEdgeWorkerIDRequest) (*edgeworkers.EdgeWorkerID, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*edgeworkers.EdgeWorkerID), args.Error(1)
}

func (m *mockEdgeworkers) DeleteEdgeWorkerID(ctx context.Context, req edgeworkers.DeleteEdgeWorkerIDRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

// EdgeWorkerVersion

func (m *mockEdgeworkers) GetEdgeWorkerVersion(ctx context.Context, req edgeworkers.GetEdgeWorkerVersionRequest) (*edgeworkers.EdgeWorkerVersion, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*edgeworkers.EdgeWorkerVersion), args.Error(1)
}

func (m *mockEdgeworkers) ListEdgeWorkerVersions(ctx context.Context, req edgeworkers.ListEdgeWorkerVersionsRequest) (*edgeworkers.ListEdgeWorkerVersionsResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*edgeworkers.ListEdgeWorkerVersionsResponse), args.Error(1)
}

func (m *mockEdgeworkers) GetEdgeWorkerVersionContent(ctx context.Context, req edgeworkers.GetEdgeWorkerVersionContentRequest) (*edgeworkers.Bundle, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*edgeworkers.Bundle), args.Error(1)
}

func (m *mockEdgeworkers) CreateEdgeWorkerVersion(ctx context.Context, req edgeworkers.CreateEdgeWorkerVersionRequest) (*edgeworkers.EdgeWorkerVersion, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*edgeworkers.EdgeWorkerVersion), args.Error(1)
}

func (m *mockEdgeworkers) DeleteEdgeWorkerVersion(ctx context.Context, req edgeworkers.DeleteEdgeWorkerVersionRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

// PermissionGroups

func (m *mockEdgeworkers) GetPermissionGroup(ctx context.Context, req edgeworkers.GetPermissionGroupRequest) (*edgeworkers.PermissionGroup, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*edgeworkers.PermissionGroup), args.Error(1)
}

func (m *mockEdgeworkers) ListPermissionGroups(ctx context.Context) (*edgeworkers.ListPermissionGroupsResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*edgeworkers.ListPermissionGroupsResponse), args.Error(1)
}

// Properties

func (m *mockEdgeworkers) ListProperties(ctx context.Context, req edgeworkers.ListPropertiesRequest) (*edgeworkers.ListPropertiesResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*edgeworkers.ListPropertiesResponse), args.Error(1)
}

// ResourceTiers

func (m *mockEdgeworkers) ListResourceTiers(ctx context.Context, req edgeworkers.ListResourceTiersRequest) (*edgeworkers.ListResourceTiersResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*edgeworkers.ListResourceTiersResponse), args.Error(1)
}

func (m *mockEdgeworkers) GetResourceTier(ctx context.Context, req edgeworkers.GetResourceTierRequest) (*edgeworkers.ResourceTier, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*edgeworkers.ResourceTier), args.Error(1)
}

// SecureTokens

func (m *mockEdgeworkers) CreateSecureToken(ctx context.Context, req edgeworkers.CreateSecureTokenRequest) (*edgeworkers.CreateSecureTokenResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*edgeworkers.CreateSecureTokenResponse), args.Error(1)
}

// Validations

func (m *mockEdgeworkers) ValidateBundle(ctx context.Context, req edgeworkers.ValidateBundleRequest) (*edgeworkers.ValidateBundleResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*edgeworkers.ValidateBundleResponse), args.Error(1)
}
