package reportinggroups

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v13/pkg/reportinggroups"
	"github.com/akamai/cli-terraform/v2/pkg/templates"
	"github.com/akamai/cli/v2/pkg/terminal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type groupMockData struct {
	id               int64
	name             string
	accessContractID string
	contractID       string
	cpCodes          []reportinggroups.CPCode
	listResponse     []reportinggroups.ReportingGroup
}

func defaultGroupMockData() groupMockData {
	return groupMockData{
		id:               42,
		name:             "test reporting group",
		accessContractID: "1-ACCGRP",
		contractID:       "1-CNTR",
		cpCodes: []reportinggroups.CPCode{
			{CPCodeID: 12345, CPCodeName: "My CP Code"},
		},
	}
}

func (d *groupMockData) buildListResponse() []reportinggroups.ReportingGroup {
	if d.listResponse != nil {
		return d.listResponse
	}
	return []reportinggroups.ReportingGroup{
		{
			ReportingGroupID:   d.id,
			ReportingGroupName: d.name,
			AccessGroup: reportinggroups.AccessGroup{
				ContractID: d.accessContractID,
			},
			Contracts: []reportinggroups.Contract{
				{
					ContractID: d.contractID,
					CPCodes:    d.cpCodes,
				},
			},
		},
	}
}

func (d *groupMockData) mockListReportingGroups(m *reportinggroups.Mock, err error) {
	req := reportinggroups.ListReportingGroupsRequest{
		ReportingGroupName: d.name,
	}
	resp := &reportinggroups.ListReportingGroupsResponse{
		Groups: d.buildListResponse(),
	}
	m.On("ListReportingGroups", mock.Anything, req).Return(resp, err).Once()
}

func TestProcessReportingGroupsTemplates(t *testing.T) {
	tests := map[string]struct {
		dir                    string
		init                   func(*groupMockData, *reportinggroups.Mock)
		edgercPath             string
		configSection          string
		withError              string
		buildTemplateProcessor func(t *testing.T) templates.TemplateProcessor
	}{
		"basic export": {
			dir: "basic",
			init: func(d *groupMockData, m *reportinggroups.Mock) {
				d.mockListReportingGroups(m, nil)
			},
		},
		"multiple cp codes": {
			dir: "multiple_cp_codes",
			init: func(d *groupMockData, m *reportinggroups.Mock) {
				d.cpCodes = []reportinggroups.CPCode{
					{CPCodeID: 11111, CPCodeName: "CP Code One"},
					{CPCodeID: 22222, CPCodeName: "CP Code Two"},
					{CPCodeID: 33333, CPCodeName: "CP Code Three"},
				}
				d.mockListReportingGroups(m, nil)
			},
		},
		"name begins with number": {
			dir: "name_begins_with_number",
			init: func(d *groupMockData, m *reportinggroups.Mock) {
				d.name = "42group"
				d.mockListReportingGroups(m, nil)
			},
		},
		"non-default edgerc path and section": {
			dir: "non_default_edgerc",
			init: func(d *groupMockData, m *reportinggroups.Mock) {
				d.mockListReportingGroups(m, nil)
			},
			edgercPath:    "/non/default/path/.edgerc",
			configSection: "other-section",
		},
		"error listing reporting groups": {
			init: func(d *groupMockData, m *reportinggroups.Mock) {
				d.mockListReportingGroups(m, fmt.Errorf("API error"))
			},
			withError: "error fetching reporting group: API error",
		},
		"error group not found": {
			init: func(d *groupMockData, m *reportinggroups.Mock) {
				d.listResponse = []reportinggroups.ReportingGroup{}
				d.mockListReportingGroups(m, nil)
			},
			withError: `error finding reporting group: no reporting group found with name "test reporting group"`,
		},
		"error no contracts in group": {
			init: func(d *groupMockData, m *reportinggroups.Mock) {
				d.listResponse = []reportinggroups.ReportingGroup{
					{
						ReportingGroupID:   d.id,
						ReportingGroupName: d.name,
						AccessGroup:        reportinggroups.AccessGroup{ContractID: d.accessContractID},
						Contracts:          []reportinggroups.Contract{},
					},
				}
				d.mockListReportingGroups(m, nil)
			},
			withError: "invalid reporting group: expected exactly one contract, got 0",
		},
		"error multiple contracts in group": {
			init: func(d *groupMockData, m *reportinggroups.Mock) {
				d.listResponse = []reportinggroups.ReportingGroup{
					{
						ReportingGroupID:   d.id,
						ReportingGroupName: d.name,
						AccessGroup:        reportinggroups.AccessGroup{ContractID: d.accessContractID},
						Contracts: []reportinggroups.Contract{
							{ContractID: d.contractID, CPCodes: d.cpCodes},
							{ContractID: "1-OTHER", CPCodes: d.cpCodes},
						},
					},
				}
				d.mockListReportingGroups(m, nil)
			},
			withError: "invalid reporting group: expected exactly one contract, got 2",
		},
		"error multiple groups with same name": {
			init: func(d *groupMockData, m *reportinggroups.Mock) {
				item := reportinggroups.ReportingGroup{
					ReportingGroupID:   d.id,
					ReportingGroupName: d.name,
					AccessGroup:        reportinggroups.AccessGroup{ContractID: d.accessContractID},
					Contracts: []reportinggroups.Contract{
						{ContractID: d.contractID, CPCodes: d.cpCodes},
					},
				}
				d.listResponse = []reportinggroups.ReportingGroup{item, item}
				d.mockListReportingGroups(m, nil)
			},
			withError: `error finding reporting group: multiple reporting groups found with name "test reporting group"`,
		},
		"templating error": {
			init: func(d *groupMockData, m *reportinggroups.Mock) {
				d.mockListReportingGroups(m, nil)
			},
			buildTemplateProcessor: func(t *testing.T) templates.TemplateProcessor {
				return templates.FSTemplateProcessor{
					TemplatesFS: templateFiles,
					TemplateTargets: map[string]string{
						"nosuchtemplate.tmpl": fmt.Sprintf("%s/nosuchtemplate.tf", t.TempDir()),
					},
				}
			},
			withError: "error saving terraform project files: no template file: nosuchtemplate.tmpl",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			m := &reportinggroups.Mock{}
			d := defaultGroupMockData()
			test.init(&d, m)

			if test.configSection == "" {
				test.configSection = "default"
			}
			if test.edgercPath == "" {
				test.edgercPath = "~/.edgerc"
			}

			params := createReportingGroupParams{
				name:          d.name,
				edgercPath:    test.edgercPath,
				configSection: test.configSection,
				client:        m,
			}

			var tempDir string
			if test.buildTemplateProcessor != nil {
				params.templateProcessor = test.buildTemplateProcessor(t)
			} else if test.dir != "" {
				tempDir = t.TempDir()
				params.templateProcessor = templates.FSTemplateProcessor{
					TemplatesFS: templateFiles,
					TemplateTargets: map[string]string{
						"reportinggroups.tmpl": fmt.Sprintf("%s/reportinggroups.tf", tempDir),
						"variables.tmpl":       fmt.Sprintf("%s/variables.tf", tempDir),
						"import.tmpl":          fmt.Sprintf("%s/import.sh", tempDir),
					},
				}
			}

			ctx := terminal.Context(context.Background(), terminal.New(terminal.DiscardWriter(), nil, terminal.DiscardWriter()))
			err := createReportingGroup(ctx, params)

			m.AssertExpectations(t)

			if test.withError != "" {
				assert.EqualError(t, err, test.withError)
				return
			}
			require.NoError(t, err)

			for _, f := range []string{"reportinggroups.tf", "import.sh", "variables.tf"} {
				expected, err := os.ReadFile(fmt.Sprintf("./testdata/%s/%s", test.dir, f))
				require.NoError(t, err)
				result, err := os.ReadFile(fmt.Sprintf("%s/%s", tempDir, f))
				require.NoError(t, err)
				assert.Equal(t, string(expected), string(result), "file %s does not match", f)
			}
		})
	}
}
