package iam

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v12/pkg/iam"
	"github.com/akamai/cli-terraform/v2/pkg/templates"
	"github.com/akamai/cli-terraform/v2/pkg/tools"
	"github.com/akamai/cli/v2/pkg/terminal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	apisGet = []iam.API{
		{
			APIID:            5801,
			APIName:          "EdgeWorkers",
			Description:      "EdgeWorkers",
			Endpoint:         "/edgeworkers/",
			DocumentationURL: "https://developer.akamai.com/api/web_performance/edgeworkers/v1.html",
			AccessLevel:      "READ-WRITE",
		},
		{
			APIID:            5580,
			APIName:          "Search Data Feed",
			Description:      "Search Data Feed",
			Endpoint:         "/search-portal-data-feed-api/",
			DocumentationURL: "/",
			AccessLevel:      "READ-ONLY",
		},
	}

	singleGroup = []iam.ClientGroup{
		{
			GroupID:         123,
			GroupName:       "group2",
			IsBlocked:       false,
			ParentGroupID:   0,
			RoleDescription: "group description",
			RoleID:          340,
			RoleName:        "role",
			Subgroups: []iam.ClientGroup{
				{
					GroupID: 333,
					RoleID:  540,
					Subgroups: []iam.ClientGroup{
						{
							GroupID: 444,
							RoleID:  640,
						},
					},
				},
			},
		},
	}

	clientIPACL = iam.IPACL{
		Enable: true,
		CIDR:   []string{"128.5.6.5/24"},
	}

	clientIPACLDisabled = iam.IPACL{
		Enable: false,
		CIDR:   []string{"128.5.6.5/24"},
	}

	clientTFIPACL = TFIPACL{
		Enable: true,
		CIDR:   []string{"128.5.6.5/24"},
	}

	clientTFIPACLDisabled = TFIPACL{
		Enable: false,
		CIDR:   []string{"128.5.6.5/24"},
	}

	clientPurgeOptions = iam.PurgeOptions{
		CanPurgeByCacheTag: true,
		CanPurgeByCPCode:   true,
		CPCodeAccess: iam.CPCodeAccess{
			AllCurrentAndNewCPCodes: false,
			CPCodes:                 []int64{101},
		},
	}

	clientTFPurgeOptions = TFPurgeOptions{
		CanPurgeByCPCode:   true,
		CanPurgeByCacheTag: true,
		CPCodeAccess: TFCPCodeAccess{
			AllCurrentAndNewCPCodes: false,
			CPCodes:                 []int64{101},
		},
	}

	getAPIClientResponse = func(ipACL *iam.IPACL, purgeOptions *iam.PurgeOptions) iam.GetAPIClientResponse {
		return iam.GetAPIClientResponse{
			AccessToken:           "access_token",
			ActiveCredentialCount: 1,
			AllowAccountSwitch:    false,
			APIAccess: iam.APIAccess{
				AllAccessibleAPIs: false,
				APIs:              apisGet,
			},
			AuthorizedUsers:         []string{"mw+2"},
			BaseURL:                 "base_url",
			CanAutoCreateCredential: false,
			ClientDescription:       "Test API Client",
			ClientID:                "1a2b3",
			ClientName:              "mw+2_1",
			ClientType:              "CLIENT",
			CreatedBy:               "someuser",
			CreatedDate:             tools.ParseRFC3339("2023-06-13T14:48:08.000Z"),
			GroupAccess: iam.GroupAccess{
				CloneAuthorizedUserGroups: false,
				Groups:                    singleGroup,
			},
			IPACL:              ipACL,
			IsLocked:           false,
			NotificationEmails: []string{"mw+2@example.com"},
			PurgeOptions:       purgeOptions,
			Credentials: []iam.APIClientCredential{
				{
					Description: "Test API Client Credential 1",
					Status:      "ACTIVE",
					ExpiresOn:   tools.ParseRFC3339("2025-06-13T14:48:08.000Z"),
					CreatedOn:   tools.ParseRFC3339("2023-06-13T14:48:08.000Z"),
				},
				{
					Description: "Test API Client Credential 2",
					Status:      "ACTIVE",
					ExpiresOn:   tools.ParseRFC3339("2025-06-13T14:48:08.000Z"),
					CreatedOn:   tools.ParseRFC3339("2022-06-13T14:48:08.000Z"),
				},
			},
		}
	}

	getAPIClientResponseNoCredentials = iam.GetAPIClientResponse{
		AccessToken:           "access_token",
		ActiveCredentialCount: 1,
		AllowAccountSwitch:    false,
		APIAccess: iam.APIAccess{
			AllAccessibleAPIs: false,
			APIs:              apisGet,
		},
		AuthorizedUsers:         []string{"mw+2"},
		BaseURL:                 "base_url",
		CanAutoCreateCredential: false,
		ClientDescription:       "Test API Client",
		ClientID:                "1a2b3",
		ClientName:              "mw+2_1",
		ClientType:              "CLIENT",
		CreatedBy:               "someuser",
		CreatedDate:             tools.ParseRFC3339("2023-06-13T14:48:08.000Z"),
		GroupAccess: iam.GroupAccess{
			CloneAuthorizedUserGroups: false,
			Groups:                    singleGroup,
		},
		IPACL: &iam.IPACL{
			Enable: true,
			CIDR:   []string{"128.5.6.5/24"},
		},
		IsLocked:           false,
		NotificationEmails: []string{"mw+2@example.com"},
		Credentials:        []iam.APIClientCredential{},
	}

	expectGetAPIClient = func(client *iam.Mock, res iam.GetAPIClientResponse) {
		req := iam.GetAPIClientRequest{
			ClientID:    "1a2b3",
			GroupAccess: true,
			APIAccess:   true,
			IPACL:       true,
			Credentials: true,
		}
		client.On("GetAPIClient", mock.Anything, req).Return(&res, nil).Once()
	}

	expectGetSelfAPIClient = func(client *iam.Mock, res iam.GetAPIClientResponse) {
		req := iam.GetAPIClientRequest{
			GroupAccess: true,
			APIAccess:   true,
			IPACL:       true,
			Credentials: true,
		}
		client.On("GetAPIClient", mock.Anything, req).Return(&res, nil).Once()
	}

	getTfClient = func(ipACL *TFIPACL, purgeOptions *TFPurgeOptions) TFClient {
		return TFClient{
			ClientID:           "1a2b3",
			AuthorizedUsers:    []string{"mw+2"},
			ClientType:         "CLIENT",
			ClientName:         "mw+2_1",
			NotificationEmails: []string{"mw+2@example.com"},
			ClientDescription:  "Test API Client",
			Lock:               false,
			Credential: &TFCredential{
				Description: "Test API Client Credential 2",
				Status:      "ACTIVE",
				ExpiresOn:   "2025-06-13T14:48:08Z",
			},
			GroupAccess: TFGroupAccessRequest{
				CloneAuthorizedUserGroups: false,
				Groups: []TFClientGroupRequestItem{
					{
						GroupID: 123,
						RoleID:  340,
						Subgroups: []TFClientGroupRequestItem{
							{
								GroupID: 333,
								RoleID:  540,
								Subgroups: []TFClientGroupRequestItem{
									{
										GroupID: 444,
										RoleID:  640,
									},
								},
							},
						},
					},
				},
			},
			IPACL: ipACL,
			APIAccess: TFAPIAccessRequest{
				AllAccessibleAPIs: false,
				APIs: []TFAPIRequestItem{
					{
						APIID:       5801,
						AccessLevel: "READ-WRITE",
					},
					{
						APIID:       5580,
						AccessLevel: "READ-ONLY",
					},
				},
			},
			PurgeOptions: purgeOptions,
		}
	}

	expectClientProcessTemplates = func(p *templates.MockProcessor, section string, data TFClient) *mock.Call {
		tfData := TFData{
			TFClient:   data,
			Section:    section,
			Subcommand: "client",
		}

		call := p.On("ProcessTemplates", tfData)
		return call.Return(nil)
	}
)

func TestCreateIAMClient(t *testing.T) {
	section := "test_section"

	tests := map[string]struct {
		init      func(*iam.Mock, *templates.MockProcessor)
		clientID  string
		withError error
	}{
		"fetch API client with client ID": {
			init: func(i *iam.Mock, p *templates.MockProcessor) {
				expectGetAPIClient(i, getAPIClientResponse(&clientIPACL, &clientPurgeOptions))
				expectClientProcessTemplates(p, section, getTfClient(&clientTFIPACL, &clientTFPurgeOptions))
			},
			clientID: "1a2b3",
		},
		"fetch self API client": {
			init: func(i *iam.Mock, p *templates.MockProcessor) {
				expectGetSelfAPIClient(i, getAPIClientResponse(&clientIPACL, &clientPurgeOptions))
				expectClientProcessTemplates(p, section, getTfClient(&clientTFIPACL, &clientTFPurgeOptions))
			},
			clientID: "",
		},
		"fetch API client no IPACL": {
			init: func(i *iam.Mock, p *templates.MockProcessor) {
				expectGetAPIClient(i, getAPIClientResponse(nil, &clientPurgeOptions))
				expectClientProcessTemplates(p, section, getTfClient(nil, &clientTFPurgeOptions))
			},
			clientID: "1a2b3",
		},
		"fetch API client with disabled IPACL": {
			init: func(i *iam.Mock, p *templates.MockProcessor) {
				expectGetAPIClient(i, getAPIClientResponse(&clientIPACLDisabled, &clientPurgeOptions))
				expectClientProcessTemplates(p, section, getTfClient(&clientTFIPACLDisabled, &clientTFPurgeOptions))
			},
			clientID: "1a2b3",
		},
		"fetch API client no PurgeOptions": {
			init: func(i *iam.Mock, p *templates.MockProcessor) {
				expectGetAPIClient(i, getAPIClientResponse(&clientIPACL, nil))
				expectClientProcessTemplates(p, section, getTfClient(&clientTFIPACL, nil))
			},
			clientID: "1a2b3",
		},
		"fetch API client no credentials - expect error": {
			init: func(i *iam.Mock, _ *templates.MockProcessor) {
				expectGetAPIClient(i, getAPIClientResponseNoCredentials)
			},
			clientID:  "1a2b3",
			withError: fmt.Errorf("It's impossible to manage API Client with no credential via Terraform"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mi := new(iam.Mock)
			mp := new(templates.MockProcessor)
			test.init(mi, mp)
			ctx := terminal.Context(context.Background(), terminal.New(terminal.DiscardWriter(), nil, terminal.DiscardWriter()))
			err := createIAMAPIClient(ctx, test.clientID, section, mi, mp)
			if test.withError != nil {
				require.ErrorContains(t, err, test.withError.Error())
				return
			}
			require.NoError(t, err)
			mi.AssertExpectations(t)
			mp.AssertExpectations(t)
		})
	}
}

func TestProcessIAMClientTemplates(t *testing.T) {
	section := "test_section"

	tests := map[string]struct {
		givenData    TFData
		dir          string
		filesToCheck []string
	}{
		"basic client": {
			givenData: TFData{
				TFClient: TFClient{
					ClientID:           "1a2b3",
					AuthorizedUsers:    []string{"mw+2"},
					ClientType:         "CLIENT",
					ClientName:         "mw+2_1",
					NotificationEmails: []string{"mw+2@example.com"},
					ClientDescription:  "Test API Client",
					Lock:               false,
					Credential: &TFCredential{
						Description: "Test API Client Credential 1",
						Status:      "ACTIVE",
						ExpiresOn:   "2027-04-09T12:34:13Z",
					},
					GroupAccess: TFGroupAccessRequest{
						CloneAuthorizedUserGroups: false,
						Groups: []TFClientGroupRequestItem{
							{
								GroupID: 123,
								RoleID:  340,
								Subgroups: []TFClientGroupRequestItem{
									{
										GroupID: 333,
										RoleID:  540,
									},
								},
							},
						},
					},
					IPACL: &TFIPACL{
						Enable: true,
						CIDR:   []string{"128.5.6.5/24"},
					},
					APIAccess: TFAPIAccessRequest{
						AllAccessibleAPIs: false,
						APIs: []TFAPIRequestItem{
							{
								APIID:       5580,
								AccessLevel: "READ-ONLY",
							},
							{
								APIID:       5801,
								AccessLevel: "READ-WRITE",
							},
						},
					},
					PurgeOptions: &TFPurgeOptions{
						CanPurgeByCPCode:   true,
						CanPurgeByCacheTag: true,
						CPCodeAccess: TFCPCodeAccess{
							AllCurrentAndNewCPCodes: false,
							CPCodes:                 []int64{101},
						},
					},
				},
				Section:    section,
				Subcommand: "client",
			},
			dir:          "iam_client",
			filesToCheck: []string{"client.tf", "import.sh", "variables.tf"},
		},
		"client with a few groups and no credential description": {
			givenData: TFData{
				TFClient: TFClient{
					ClientID:           "1a2b3",
					AuthorizedUsers:    []string{"mw+2"},
					ClientType:         "CLIENT",
					ClientName:         "mw+2_1",
					NotificationEmails: []string{"mw+2@example.com"},
					ClientDescription:  "Test API Client",
					Lock:               false,
					Credential: &TFCredential{
						Description: "",
						Status:      "ACTIVE",
						ExpiresOn:   "2027-04-09T12:34:13Z",
					},
					GroupAccess: TFGroupAccessRequest{
						CloneAuthorizedUserGroups: false,
						Groups: []TFClientGroupRequestItem{
							{
								GroupID: 123,
								RoleID:  340,
								Subgroups: []TFClientGroupRequestItem{
									{
										GroupID: 333,
										RoleID:  540,
									},
								},
							},
							{
								GroupID: 345,
								RoleID:  340,
								Subgroups: []TFClientGroupRequestItem{
									{
										GroupID: 333,
										RoleID:  540,
									},
								},
							},
						},
					},
					IPACL: &TFIPACL{
						Enable: true,
						CIDR:   []string{"128.5.6.5/24"},
					},
					APIAccess: TFAPIAccessRequest{
						AllAccessibleAPIs: false,
						APIs: []TFAPIRequestItem{
							{
								APIID:       5580,
								AccessLevel: "READ-ONLY",
							},
							{
								APIID:       5801,
								AccessLevel: "READ-WRITE",
							},
						},
					},
					PurgeOptions: &TFPurgeOptions{
						CanPurgeByCPCode:   true,
						CanPurgeByCacheTag: true,
						CPCodeAccess: TFCPCodeAccess{
							AllCurrentAndNewCPCodes: false,
							CPCodes:                 []int64{101},
						},
					},
				},
				Section:    section,
				Subcommand: "client",
			},
			dir:          "iam_client_with_groups",
			filesToCheck: []string{"client.tf", "import.sh", "variables.tf"},
		},
		"client no subgroups": {
			givenData: TFData{
				TFClient: TFClient{
					ClientID:           "1a2b3",
					AuthorizedUsers:    []string{"mw+2"},
					ClientType:         "CLIENT",
					ClientName:         "mw+2_1",
					NotificationEmails: []string{"mw+2@example.com"},
					ClientDescription:  "Test API Client",
					Lock:               false,
					Credential: &TFCredential{
						Description: "",
						Status:      "INACTIVE",
						ExpiresOn:   "2027-04-09T12:34:13Z",
					},
					GroupAccess: TFGroupAccessRequest{
						CloneAuthorizedUserGroups: false,
						Groups: []TFClientGroupRequestItem{
							{
								GroupID: 123,
								RoleID:  340,
							},
						},
					},
					IPACL: &TFIPACL{
						Enable: true,
						CIDR:   []string{"128.5.6.5/24"},
					},
					APIAccess: TFAPIAccessRequest{
						AllAccessibleAPIs: false,
						APIs: []TFAPIRequestItem{
							{
								APIID:       5580,
								AccessLevel: "READ-ONLY",
							},
							{
								APIID:       5801,
								AccessLevel: "READ-WRITE",
							},
						},
					},
					PurgeOptions: &TFPurgeOptions{
						CanPurgeByCPCode:   true,
						CanPurgeByCacheTag: true,
						CPCodeAccess: TFCPCodeAccess{
							AllCurrentAndNewCPCodes: false,
							CPCodes:                 []int64{101},
						},
					},
				},
				Section:    section,
				Subcommand: "client",
			},
			dir:          "iam_client_no_subgroups",
			filesToCheck: []string{"client.tf", "import.sh", "variables.tf"},
		},
		"client no cp codes": {
			givenData: TFData{
				TFClient: TFClient{
					ClientID:           "1a2b3",
					AuthorizedUsers:    []string{"mw+2"},
					ClientType:         "CLIENT",
					ClientName:         "mw+2_1",
					NotificationEmails: []string{"mw+2@example.com"},
					ClientDescription:  "Test API Client",
					Lock:               false,
					Credential: &TFCredential{
						Description: "",
						Status:      "DELETED",
						ExpiresOn:   "2027-04-09T12:34:13Z",
					},
					GroupAccess: TFGroupAccessRequest{
						CloneAuthorizedUserGroups: false,
						Groups: []TFClientGroupRequestItem{
							{
								GroupID: 123,
								RoleID:  340,
								Subgroups: []TFClientGroupRequestItem{
									{
										GroupID: 333,
										RoleID:  540,
									},
								},
							},
						},
					},
					IPACL: &TFIPACL{
						Enable: true,
						CIDR:   []string{"128.5.6.5/24"},
					},
					APIAccess: TFAPIAccessRequest{
						AllAccessibleAPIs: false,
						APIs: []TFAPIRequestItem{
							{
								APIID:       5580,
								AccessLevel: "READ-ONLY",
							},
							{
								APIID:       5801,
								AccessLevel: "READ-WRITE",
							},
						},
					},
					PurgeOptions: &TFPurgeOptions{
						CanPurgeByCPCode:   true,
						CanPurgeByCacheTag: true,
						CPCodeAccess: TFCPCodeAccess{
							AllCurrentAndNewCPCodes: false,
						},
					},
				},
				Section:    section,
				Subcommand: "client",
			},
			dir:          "iam_client_empty_cp_codes",
			filesToCheck: []string{"client.tf", "import.sh", "variables.tf"},
		},
		"client no cidr": {
			givenData: TFData{
				TFClient: TFClient{
					ClientID:           "1a2b3",
					AuthorizedUsers:    []string{"mw+2"},
					ClientType:         "CLIENT",
					ClientName:         "mw+2_1",
					NotificationEmails: []string{"mw+2@example.com"},
					ClientDescription:  "Test API Client",
					Lock:               false,
					Credential: &TFCredential{
						Description: "Test API Client Credential 1",
						Status:      "ACTIVE",
						ExpiresOn:   "2027-04-09T12:34:13Z",
					},
					GroupAccess: TFGroupAccessRequest{
						CloneAuthorizedUserGroups: false,
						Groups: []TFClientGroupRequestItem{
							{
								GroupID: 123,
								RoleID:  340,
								Subgroups: []TFClientGroupRequestItem{
									{
										GroupID: 333,
										RoleID:  540,
									},
								},
							},
						},
					},
					IPACL: &TFIPACL{
						Enable: true,
					},
					APIAccess: TFAPIAccessRequest{
						AllAccessibleAPIs: false,
						APIs: []TFAPIRequestItem{
							{
								APIID:       5580,
								AccessLevel: "READ-ONLY",
							},
							{
								APIID:       5801,
								AccessLevel: "READ-WRITE",
							},
						},
					},
					PurgeOptions: &TFPurgeOptions{
						CanPurgeByCPCode:   true,
						CanPurgeByCacheTag: true,
						CPCodeAccess: TFCPCodeAccess{
							AllCurrentAndNewCPCodes: false,
							CPCodes:                 []int64{101},
						},
					},
				},
				Section:    section,
				Subcommand: "client",
			},
			dir:          "iam_client_no_cidr",
			filesToCheck: []string{"client.tf", "import.sh", "variables.tf"},
		},
		"client with recursion": {
			givenData: TFData{
				TFClient: TFClient{
					ClientID:           "1a2b3",
					AuthorizedUsers:    []string{"mw+2"},
					ClientType:         "CLIENT",
					ClientName:         "mw+2_1",
					NotificationEmails: []string{"mw+2@example.com"},
					ClientDescription:  "Test API Client",
					Lock:               false,
					Credential: &TFCredential{
						Description: "Test API Client Credential 1",
						Status:      "ACTIVE",
						ExpiresOn:   "2027-04-09T12:34:13Z",
					},
					GroupAccess: TFGroupAccessRequest{
						CloneAuthorizedUserGroups: false,
						Groups: []TFClientGroupRequestItem{
							{
								GroupID: 123,
								RoleID:  340,
								Subgroups: []TFClientGroupRequestItem{
									{
										GroupID: 333,
										RoleID:  540,
										Subgroups: []TFClientGroupRequestItem{
											{
												GroupID: 555,
												RoleID:  640,
											},
										},
									},
								},
							},
						},
					},
					IPACL: &TFIPACL{
						Enable: true,
						CIDR:   []string{"128.5.6.5/24"},
					},
					APIAccess: TFAPIAccessRequest{
						AllAccessibleAPIs: false,
						APIs: []TFAPIRequestItem{
							{
								APIID:       5580,
								AccessLevel: "READ-ONLY",
							},
							{
								APIID:       5801,
								AccessLevel: "READ-WRITE",
							},
						},
					},
					PurgeOptions: &TFPurgeOptions{
						CanPurgeByCPCode:   true,
						CanPurgeByCacheTag: true,
						CPCodeAccess: TFCPCodeAccess{
							AllCurrentAndNewCPCodes: false,
							CPCodes:                 []int64{101},
						},
					},
				},
				Section:    section,
				Subcommand: "client",
			},
			dir:          "iam_client_recursion",
			filesToCheck: []string{"client.tf", "import.sh", "variables.tf"},
		},
		"client with clone_authorized_user_groups true": {
			givenData: TFData{
				TFClient: TFClient{
					ClientID:           "1a2b3",
					AuthorizedUsers:    []string{"mw+2"},
					ClientType:         "CLIENT",
					ClientName:         "mw+2_1",
					NotificationEmails: []string{"mw+2@example.com"},
					ClientDescription:  "Test API Client",
					Lock:               false,
					Credential: &TFCredential{
						Description: "Test API Client Credential 1",
						Status:      "ACTIVE",
						ExpiresOn:   "2027-04-09T12:34:13Z",
					},
					GroupAccess: TFGroupAccessRequest{
						CloneAuthorizedUserGroups: true,
						Groups: []TFClientGroupRequestItem{
							{
								GroupID: 123,
								RoleID:  340,
								Subgroups: []TFClientGroupRequestItem{
									{
										GroupID: 333,
										RoleID:  540,
									},
								},
							},
						},
					},
					IPACL: &TFIPACL{
						Enable: true,
						CIDR:   []string{"128.5.6.5/24"},
					},
					APIAccess: TFAPIAccessRequest{
						AllAccessibleAPIs: false,
						APIs: []TFAPIRequestItem{
							{
								APIID:       5580,
								AccessLevel: "READ-ONLY",
							},
							{
								APIID:       5801,
								AccessLevel: "READ-WRITE",
							},
						},
					},
					PurgeOptions: &TFPurgeOptions{
						CanPurgeByCPCode:   true,
						CanPurgeByCacheTag: true,
						CPCodeAccess: TFCPCodeAccess{
							AllCurrentAndNewCPCodes: false,
							CPCodes:                 []int64{101},
						},
					},
				},
				Section:    section,
				Subcommand: "client",
			},
			dir:          "iam_client_clone_authorized_user_groups_is_true",
			filesToCheck: []string{"client.tf", "import.sh", "variables.tf"},
		},
		"client with all_accessible_apis true": {
			givenData: TFData{
				TFClient: TFClient{
					ClientID:           "1a2b3",
					AuthorizedUsers:    []string{"mw+2"},
					ClientType:         "CLIENT",
					ClientName:         "mw+2_1",
					NotificationEmails: []string{"mw+2@example.com"},
					ClientDescription:  "Test API Client",
					Lock:               false,
					Credential: &TFCredential{
						Description: "Test API Client Credential 1",
						Status:      "ACTIVE",
						ExpiresOn:   "2027-04-09T12:34:13Z",
					},
					GroupAccess: TFGroupAccessRequest{
						CloneAuthorizedUserGroups: false,
						Groups: []TFClientGroupRequestItem{
							{
								GroupID: 123,
								RoleID:  340,
								Subgroups: []TFClientGroupRequestItem{
									{
										GroupID: 333,
										RoleID:  540,
									},
								},
							},
						},
					},
					IPACL: &TFIPACL{
						Enable: true,
						CIDR:   []string{"128.5.6.5/24"},
					},
					APIAccess: TFAPIAccessRequest{
						AllAccessibleAPIs: true,
						APIs: []TFAPIRequestItem{
							{
								APIID:       5580,
								AccessLevel: "READ-ONLY",
							},
							{
								APIID:       5801,
								AccessLevel: "READ-WRITE",
							},
						},
					},
					PurgeOptions: &TFPurgeOptions{
						CanPurgeByCPCode:   true,
						CanPurgeByCacheTag: true,
						CPCodeAccess: TFCPCodeAccess{
							AllCurrentAndNewCPCodes: false,
							CPCodes:                 []int64{101},
						},
					},
				},
				Section:    section,
				Subcommand: "client",
			},
			dir:          "iam_client_all_accessible_apis_is_true",
			filesToCheck: []string{"client.tf", "import.sh", "variables.tf"},
		},
		"client with all_current_and_new_cp_codes true": {
			givenData: TFData{
				TFClient: TFClient{
					ClientID:           "1a2b3",
					AuthorizedUsers:    []string{"mw+2"},
					ClientType:         "CLIENT",
					ClientName:         "mw+2_1",
					NotificationEmails: []string{"mw+2@example.com"},
					ClientDescription:  "Test API Client",
					Lock:               false,
					Credential: &TFCredential{
						Description: "Test API Client Credential 1",
						Status:      "ACTIVE",
						ExpiresOn:   "2027-04-09T12:34:13Z",
					},
					GroupAccess: TFGroupAccessRequest{
						CloneAuthorizedUserGroups: false,
						Groups: []TFClientGroupRequestItem{
							{
								GroupID: 123,
								RoleID:  340,
								Subgroups: []TFClientGroupRequestItem{
									{
										GroupID: 333,
										RoleID:  540,
									},
								},
							},
						},
					},
					IPACL: &TFIPACL{
						Enable: true,
						CIDR:   []string{"128.5.6.5/24"},
					},
					APIAccess: TFAPIAccessRequest{
						AllAccessibleAPIs: false,
						APIs: []TFAPIRequestItem{
							{
								APIID:       5580,
								AccessLevel: "READ-ONLY",
							},
							{
								APIID:       5801,
								AccessLevel: "READ-WRITE",
							},
						},
					},
					PurgeOptions: &TFPurgeOptions{
						CanPurgeByCPCode:   true,
						CanPurgeByCacheTag: true,
						CPCodeAccess: TFCPCodeAccess{
							AllCurrentAndNewCPCodes: true,
						},
					},
				},
				Section:    section,
				Subcommand: "client",
			},
			dir:          "iam_client_all_current_and_new_cp_codes_is_true",
			filesToCheck: []string{"client.tf", "import.sh", "variables.tf"},
		},
		"client no IPACL": {
			givenData: TFData{
				TFClient: TFClient{
					ClientID:           "1a2b3",
					AuthorizedUsers:    []string{"mw+2"},
					ClientType:         "CLIENT",
					ClientName:         "mw+2_1",
					NotificationEmails: []string{"mw+2@example.com"},
					ClientDescription:  "Test API Client",
					Lock:               false,
					Credential: &TFCredential{
						Description: "Test API Client Credential 1",
						Status:      "ACTIVE",
						ExpiresOn:   "2027-04-09T12:34:13Z",
					},
					GroupAccess: TFGroupAccessRequest{
						CloneAuthorizedUserGroups: false,
						Groups: []TFClientGroupRequestItem{
							{
								GroupID: 123,
								RoleID:  340,
								Subgroups: []TFClientGroupRequestItem{
									{
										GroupID: 333,
										RoleID:  540,
									},
								},
							},
						},
					},
					APIAccess: TFAPIAccessRequest{
						AllAccessibleAPIs: false,
						APIs: []TFAPIRequestItem{
							{
								APIID:       5580,
								AccessLevel: "READ-ONLY",
							},
							{
								APIID:       5801,
								AccessLevel: "READ-WRITE",
							},
						},
					},
					PurgeOptions: &TFPurgeOptions{
						CanPurgeByCPCode:   true,
						CanPurgeByCacheTag: true,
						CPCodeAccess: TFCPCodeAccess{
							AllCurrentAndNewCPCodes: false,
							CPCodes:                 []int64{101},
						},
					},
				},
				Section:    section,
				Subcommand: "client",
			},
			dir:          "iam_client_no_ipacl",
			filesToCheck: []string{"client.tf", "import.sh", "variables.tf"},
		},
		"client with disabled IPACL": {
			givenData: TFData{
				TFClient: TFClient{
					ClientID:           "1a2b3",
					AuthorizedUsers:    []string{"mw+2"},
					ClientType:         "CLIENT",
					ClientName:         "mw+2_1",
					NotificationEmails: []string{"mw+2@example.com"},
					ClientDescription:  "Test API Client",
					Lock:               false,
					Credential: &TFCredential{
						Description: "Test API Client Credential 1",
						Status:      "ACTIVE",
						ExpiresOn:   "2027-04-09T12:34:13Z",
					},
					IPACL: &TFIPACL{
						Enable: false,
						CIDR:   []string{"1.2.3.4/24"},
					},
					GroupAccess: TFGroupAccessRequest{
						CloneAuthorizedUserGroups: false,
						Groups: []TFClientGroupRequestItem{
							{
								GroupID: 123,
								RoleID:  340,
								Subgroups: []TFClientGroupRequestItem{
									{
										GroupID: 333,
										RoleID:  540,
									},
								},
							},
						},
					},
					APIAccess: TFAPIAccessRequest{
						AllAccessibleAPIs: false,
						APIs: []TFAPIRequestItem{
							{
								APIID:       5580,
								AccessLevel: "READ-ONLY",
							},
							{
								APIID:       5801,
								AccessLevel: "READ-WRITE",
							},
						},
					},
					PurgeOptions: &TFPurgeOptions{
						CanPurgeByCPCode:   true,
						CanPurgeByCacheTag: true,
						CPCodeAccess: TFCPCodeAccess{
							AllCurrentAndNewCPCodes: false,
							CPCodes:                 []int64{101},
						},
					},
				},
				Section:    section,
				Subcommand: "client",
			},
			dir:          "iam_client_disabled_ipacl",
			filesToCheck: []string{"client.tf", "import.sh", "variables.tf"},
		},
		"client no PurgeOptions": {
			givenData: TFData{
				TFClient: TFClient{
					ClientID:           "1a2b3",
					AuthorizedUsers:    []string{"mw+2"},
					ClientType:         "CLIENT",
					ClientName:         "mw+2_1",
					NotificationEmails: []string{"mw+2@example.com"},
					ClientDescription:  "Test API Client",
					Lock:               false,
					Credential: &TFCredential{
						Description: "Test API Client Credential 1",
						Status:      "ACTIVE",
						ExpiresOn:   "2027-04-09T12:34:13Z",
					},
					GroupAccess: TFGroupAccessRequest{
						CloneAuthorizedUserGroups: false,
						Groups: []TFClientGroupRequestItem{
							{
								GroupID: 123,
								RoleID:  340,
								Subgroups: []TFClientGroupRequestItem{
									{
										GroupID: 333,
										RoleID:  540,
									},
								},
							},
						},
					},
					IPACL: &TFIPACL{
						Enable: true,
						CIDR:   []string{"128.5.6.5/24"},
					},
					APIAccess: TFAPIAccessRequest{
						AllAccessibleAPIs: false,
						APIs: []TFAPIRequestItem{
							{
								APIID:       5580,
								AccessLevel: "READ-ONLY",
							},
							{
								APIID:       5801,
								AccessLevel: "READ-WRITE",
							},
						},
					},
				},
				Section:    section,
				Subcommand: "client",
			},
			dir:          "iam_client_no_purge_options",
			filesToCheck: []string{"client.tf", "import.sh", "variables.tf"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			require.NoError(t, os.MkdirAll(fmt.Sprintf("./testdata/res/%s", test.dir), 0755))
			processor := templates.FSTemplateProcessor{
				TemplatesFS: templateFiles,
				TemplateTargets: map[string]string{
					"client.tmpl":    fmt.Sprintf("./testdata/res/%s/client.tf", test.dir),
					"variables.tmpl": fmt.Sprintf("./testdata/res/%s/variables.tf", test.dir),
					"imports.tmpl":   fmt.Sprintf("./testdata/res/%s/import.sh", test.dir),
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
