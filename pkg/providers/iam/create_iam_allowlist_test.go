package iam

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/iam"
	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/ptr"
	"github.com/akamai/cli-terraform/pkg/templates"
	"github.com/akamai/cli/pkg/terminal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	expectGetIPAllowlistStatus = func(client *iam.Mock, enabled bool) {
		client.On("GetIPAllowlistStatus", mock.Anything).Return(&iam.GetIPAllowlistStatusResponse{
			Enabled: enabled,
		}, nil).Once()
	}

	expectListCIDRBlocks = func(client *iam.Mock, cidrs iam.ListCIDRBlocksResponse) {
		req := iam.ListCIDRBlocksRequest{
			Actions: true,
		}
		client.On("ListCIDRBlocks", mock.Anything, req).Return(cidrs, nil).Once()
	}

	expectAllowlistProcessTemplates = func(p *templates.MockProcessor, section string, response iam.ListCIDRBlocksResponse, enabled bool) *mock.Call {
		tfData := TFData{
			TFAllowlist: TFAllowlist{},
			Section:     section,
			Subcommand:  "allowlist",
		}

		for _, cidr := range response {
			tfData.TFAllowlist.CIDRBlocks = append(tfData.TFAllowlist.CIDRBlocks, TFCIDRBlock{
				CIDRBlockID: cidr.CIDRBlockID,
				CIDRBlock:   cidr.CIDRBlock,
				Comments:    cidr.Comments,
				Enabled:     cidr.Enabled,
			})
		}
		tfData.TFAllowlist.Enabled = enabled

		call := p.On(
			"ProcessTemplates",
			tfData,
		)
		return call.Return(nil)
	}
)

func TestCreateIAMAllowlist(t *testing.T) {
	section := "test_section"

	tests := map[string]struct {
		init func(*iam.Mock, *templates.MockProcessor)
	}{
		"fetch ip allowlist status and CIDR blocks": {
			init: func(i *iam.Mock, p *templates.MockProcessor) {
				expectGetIPAllowlistStatus(i, true)
				expectListCIDRBlocks(i, cidrs)
				expectAllowlistProcessTemplates(p, section, cidrs, true)
			},
		},
		"fetch ip allowlist status and no CIDR blocks": {
			init: func(i *iam.Mock, p *templates.MockProcessor) {
				noCIDRBlocks := iam.ListCIDRBlocksResponse{}
				expectGetIPAllowlistStatus(i, false)
				expectListCIDRBlocks(i, noCIDRBlocks)
				expectAllowlistProcessTemplates(p, section, noCIDRBlocks, false)
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mi := new(iam.Mock)
			mp := new(templates.MockProcessor)
			test.init(mi, mp)
			ctx := terminal.Context(context.Background(), terminal.New(terminal.DiscardWriter(), nil, terminal.DiscardWriter()))
			err := createIAMAllowlist(ctx, section, mi, mp)
			require.NoError(t, err)
			mi.AssertExpectations(t)
			mp.AssertExpectations(t)
		})
	}
}

func TestProcessIAMAllowlistTemplates(t *testing.T) {
	section := "test_section"

	tests := map[string]struct {
		givenData    TFData
		dir          string
		filesToCheck []string
	}{
		"single cidr block with all fields, allowlist enabled": {
			givenData: TFData{
				TFAllowlist: TFAllowlist{
					CIDRBlocks: []TFCIDRBlock{
						{
							CIDRBlockID: 1,
							CIDRBlock:   "1.1.1.1/1",
							Enabled:     true,
							Comments:    ptr.To("comment"),
						},
					},
					Enabled: true,
				},
				Section:    section,
				Subcommand: "allowlist",
			},
			dir:          "iam_allowlist/single",
			filesToCheck: []string{"variables.tf", "import.sh", "allowlist.tf"},
		},
		"multiple cidr blocks, allowlist disabled": {
			givenData: TFData{
				TFAllowlist: TFAllowlist{
					CIDRBlocks: []TFCIDRBlock{
						{
							CIDRBlockID: 1,
							CIDRBlock:   "1.1.1.1/1",
							Enabled:     true,
							Comments:    ptr.To("comment"),
						},
						{
							CIDRBlockID: 2,
							CIDRBlock:   "2.2.2.2/2",
							Enabled:     false,
						},
					},
					Enabled: false,
				},
				Section:    section,
				Subcommand: "allowlist",
			},
			dir:          "iam_allowlist/multiple",
			filesToCheck: []string{"variables.tf", "import.sh", "allowlist.tf"},
		},
		"no cidr blocks, allowlist disabled": {
			givenData: TFData{
				TFAllowlist: TFAllowlist{
					CIDRBlocks: []TFCIDRBlock{},
					Enabled:    false,
				},
				Section:    section,
				Subcommand: "allowlist",
			},
			dir:          "iam_allowlist/no_cidr_blocks",
			filesToCheck: []string{"variables.tf", "import.sh", "allowlist.tf"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			require.NoError(t, os.MkdirAll(fmt.Sprintf("./testdata/res/%s", test.dir), 0755))
			processor := templates.FSTemplateProcessor{
				TemplatesFS: templateFiles,
				TemplateTargets: map[string]string{
					"allowlist.tmpl": fmt.Sprintf("./testdata/res/%s/allowlist.tf", test.dir),
					"imports.tmpl":   fmt.Sprintf("./testdata/res/%s/import.sh", test.dir),
					"variables.tmpl": fmt.Sprintf("./testdata/res/%s/variables.tf", test.dir),
				},
				AdditionalFuncs: additionalFunctions,
			}
			require.NoError(t, processor.ProcessTemplates(test.givenData))

			for _, f := range test.filesToCheck {
				expected, err := os.ReadFile(fmt.Sprintf("./testdata/%s/%s", test.dir, f))
				require.NoError(t, err)
				result, err := os.ReadFile(fmt.Sprintf("./testdata/res/%s/%s", test.dir, f))
				require.NoError(t, err)
				assert.Equal(t, string(expected), string(result))
			}
		})
	}
}
