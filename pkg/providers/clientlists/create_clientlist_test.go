package clientlists

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v9/pkg/clientlists"
	"github.com/akamai/cli-terraform/pkg/templates"
	"github.com/akamai/cli/pkg/terminal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	if err := os.MkdirAll("./testdata/res", 0755); err != nil {
		log.Fatal(err)
	}
	exitCode := m.Run()
	if err := os.RemoveAll("./testdata/res"); err != nil {
		log.Fatal(err)
	}
	os.Exit(exitCode)
}

func TestTemplateProcessing(t *testing.T) {
	listID := "123_ABC"
	emptyClientList := clientlists.GetClientListResponse{
		ListContent: clientlists.ListContent{
			Version:                    1,
			Type:                       "IP",
			Tags:                       []string{"tag1", "tag2"},
			Name:                       "Test Client List",
			Notes:                      "Some Notes",
			ItemsCount:                 0,
			StagingActivationStatus:    "INACTIVE",
			ProductionActivationStatus: "INACTIVE",
		},
		ContractID: "12_CA",
		GroupID:    12,
		Items:      []clientlists.ListItemContent{},
	}

	clientListItems := []clientlists.ListItemContent{
		{
			Value:          "1.1.1.1",
			Description:    "item 1",
			Tags:           []string{"t1"},
			ExpirationDate: "2026-12-26T01:00:00+00:00",
		},
		{
			Value:          "1.1.1.2",
			Description:    "",
			Tags:           []string{},
			ExpirationDate: "",
		},
		{
			Value: "1.1.1.3",
		},
	}

	clientListWithItems := clientlists.GetClientListResponse{
		ListContent: clientlists.ListContent{
			Version:                    1,
			Type:                       "IP",
			Tags:                       []string{"tag1", "tag2"},
			Name:                       "Test Client List",
			Notes:                      "Some Notes",
			ItemsCount:                 int64(len(clientListItems)),
			StagingActivationStatus:    "INACTIVE",
			ProductionActivationStatus: "INACTIVE",
		},
		ContractID: "12_CA",
		GroupID:    12,
		Items:      clientListItems,
	}

	clientListActive := clientlists.GetClientListResponse{
		ListContent: clientlists.ListContent{
			Version:                    1,
			Type:                       "IP",
			Tags:                       []string{"tag1", "tag2"},
			Name:                       "Test Client List",
			Notes:                      "Some Notes",
			ItemsCount:                 0,
			StagingActivationStatus:    "ACTIVE",
			ProductionActivationStatus: "ACTIVE",
		},
		ContractID: "12_CA",
		GroupID:    12,
	}

	tests := map[string]struct {
		init       func(*clientlists.Mock, string)
		dir        string
		checkItems bool
		withError  string
	}{
		"inactive empty client list": {
			init: func(m *clientlists.Mock, listID string) {
				mockGetClientList(m, &emptyClientList, listID, 1)
			},
			dir: "empty_list",
		},
		"inactive client list with items": {
			init: func(m *clientlists.Mock, listID string) {
				mockGetClientList(m, &clientListWithItems, listID, 1)
			},
			dir:        "list_with_items",
			checkItems: true,
		},
		"active client list": {
			init: func(m *clientlists.Mock, listID string) {
				mockGetClientList(m, &clientListActive, listID, 1)
				mockGetActivationStatus(m, listID, &clientlists.GetActivationStatusResponse{
					ActivationStatus:       "ACTIVE",
					Comments:               "Staging Activation",
					NotificationRecipients: []string{"a@b.com", "c@d.com"},
					SiebelTicketID:         "12_AB",
				}, clientlists.Staging, 1)
				mockGetActivationStatus(m, listID, &clientlists.GetActivationStatusResponse{
					ActivationStatus:       "ACTIVE",
					Comments:               "Production Activation",
					NotificationRecipients: []string{"1@2.com", "3@4.com"},
					SiebelTicketID:         "34_CD",
				}, clientlists.Production, 1)
			},
			dir: "active_list",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			resDir := fmt.Sprintf("./testdata/res/%s", test.dir)
			testDir := fmt.Sprintf("./testdata/%s", test.dir)
			mc := new(clientlists.Mock)
			processor := buildTemplateProcessor(resDir)
			ctx := terminal.Context(context.Background(), terminal.New(terminal.DiscardWriter(), nil, terminal.DiscardWriter()))
			test.init(mc, listID)
			require.NoError(t, os.MkdirAll(resDir, 0755))

			err := createClientList(ctx, listID, "edgerc_path", "test_section", resDir, mc, processor)
			require.NoError(t, err)

			for _, f := range []string{"client-list.tf", "variables.tf", "imports.sh"} {
				expected, err := os.ReadFile(fmt.Sprintf("%s/%s", testDir, f))
				require.NoError(t, err)
				result, err := os.ReadFile(fmt.Sprintf("%s/%s", resDir, f))
				require.NoError(t, err)
				assert.Equal(t, string(expected), string(result))
			}

			if test.checkItems {
				expected, err := os.ReadFile(fmt.Sprintf("%s/%s.json", testDir, listID))
				require.NoError(t, err)
				result, err := os.ReadFile(fmt.Sprintf("%s/%s.json", resDir, listID))
				require.NoError(t, err)
				assert.Equal(t, string(expected), string(result))
			}
			mc.AssertExpectations(t)
		})
	}
}

func TestFailures(t *testing.T) {

	tests := map[string]struct {
		init func(*clientlists.Mock, *templates.MockProcessor, string)
		err  error
	}{
		"get client list failed": {
			init: func(m *clientlists.Mock, p *templates.MockProcessor, listID string) {
				mockGetClientListErr(m, listID)
				mockTemplateProcessor(p, nil)
			},
			err: ErrClientListNotFound,
		},
		"templates processing failed": {
			init: func(m *clientlists.Mock, p *templates.MockProcessor, listID string) {
				mockGetClientList(m, &clientlists.GetClientListResponse{}, listID, 1)
				mockTemplateProcessor(p, errors.New("error"))
				mockGetActivationStatus(m, listID, &clientlists.GetActivationStatusResponse{ActivationStatus: "ACTIVE"}, clientlists.Staging, 1)
				mockGetActivationStatus(m, listID, &clientlists.GetActivationStatusResponse{ActivationStatus: "ACTIVE"}, clientlists.Production, 1)
			},
			err: ErrSavingFiles,
		},
		"get activation status failed": {
			init: func(m *clientlists.Mock, p *templates.MockProcessor, listID string) {
				mockGetClientList(m, &clientlists.GetClientListResponse{}, listID, 1)
				mockTemplateProcessor(p, nil)
				mockGetActivationStatusErr(m, listID, errors.New("error"), clientlists.Staging, 1)
				mockGetActivationStatusErr(m, listID, errors.New("error"), clientlists.Production, 1)
			},
			err: ErrActivationDetails,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			listID := "123_ABC"
			mc := new(clientlists.Mock)
			p := new(templates.MockProcessor)
			ctx := terminal.Context(context.Background(), terminal.New(terminal.DiscardWriter(), nil, terminal.DiscardWriter()))
			test.init(mc, p, listID)

			err := createClientList(ctx, listID, "edgerc_path", "test_section", "./", mc, p)
			assert.Error(t, err)
			assert.Contains(t, err.Error(), test.err.Error())
		})
	}
}

func mockGetClientListErr(m *clientlists.Mock, listID string) {
	m.On("GetClientList", mock.Anything, clientlists.GetClientListRequest{ListID: listID, IncludeItems: true}).Return(nil, errors.New("error")).Once()
}

func mockGetClientList(m *clientlists.Mock, res *clientlists.GetClientListResponse, listID string, times int) {
	res.ListID = listID
	m.On("GetClientList", mock.Anything, clientlists.GetClientListRequest{ListID: listID, IncludeItems: true}).Return(res, nil).Times(times)
}

func mockGetActivationStatus(m *clientlists.Mock, listID string, res *clientlists.GetActivationStatusResponse, network clientlists.ActivationNetwork, times int) {
	m.On("GetActivationStatus", mock.Anything, clientlists.GetActivationStatusRequest{
		ListID:  listID,
		Network: network,
	}).Return(res, nil).Times(times)
}

func mockGetActivationStatusErr(m *clientlists.Mock, listID string, err error, network clientlists.ActivationNetwork, times int) {
	m.On("GetActivationStatus", mock.Anything, clientlists.GetActivationStatusRequest{
		ListID:  listID,
		Network: network,
	}).Return(nil, err).Times(times)
}

func mockTemplateProcessor(p *templates.MockProcessor, err error) *templates.MockProcessor {
	p.On("ProcessTemplates", mock.Anything).Return(err).Once()
	return p
}

func buildTemplateProcessor(dir string) *templates.FSTemplateProcessor {
	templateToFile := map[string]string{
		"client-list.tmpl": fmt.Sprintf("%s/client-list.tf", dir),
		"variables.tmpl":   fmt.Sprintf("%s/variables.tf", dir),
		"imports.tmpl":     fmt.Sprintf("%s/imports.sh", dir),
	}

	return &templates.FSTemplateProcessor{
		TemplatesFS:     templateFiles,
		TemplateTargets: templateToFile,
	}
}
