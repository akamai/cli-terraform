package papi

import (
	"context"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/hapi"
	"github.com/stretchr/testify/mock"
)

type mockhapi struct {
	mock.Mock
}

func (m *mockhapi) DeleteEdgeHostname(ctx context.Context, r hapi.DeleteEdgeHostnameRequest) (*hapi.DeleteEdgeHostnameResponse, error) {
	args := m.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*hapi.DeleteEdgeHostnameResponse), args.Error(1)
}

func (m *mockhapi) GetEdgeHostname(ctx context.Context, i int) (*hapi.GetEdgeHostnameResponse, error) {
	args := m.Called(ctx, i)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*hapi.GetEdgeHostnameResponse), args.Error(1)
}

func (m *mockhapi) UpdateEdgeHostname(ctx context.Context, r hapi.UpdateEdgeHostnameRequest) (*hapi.UpdateEdgeHostnameResponse, error) {
	args := m.Called(ctx, r)

	if args.Get(0) == nil {
		return nil, args.Error(1)
	}

	return args.Get(0).(*hapi.UpdateEdgeHostnameResponse), args.Error(1)
}
