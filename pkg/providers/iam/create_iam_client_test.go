package iam

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/iam"
	"github.com/akamai/cli-terraform/v2/pkg/templates"
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

	clientTfIPACL = IPACL{
		Enable: true,
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

	clientTfPurgeOptions = PurgeOptions{
		CanPurgeByCPCode:   true,
		CanPurgeByCacheTag: true,
		CPCodeAccess: CPCodeAccess{
			AllCurrentAndNewCPCodes: false,
			CPCodes:                 []int64{101},
		},
	}

	getAPIClientResponse = func(IPACL *iam.IPACL, purgeOptions *iam.PurgeOptions) iam.GetAPIClientResponse {
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
			CreatedDate:             time.Time{},
			GroupAccess: iam.GroupAccess{
				CloneAuthorizedUserGroups: false,
				Groups:                    singleGroup,
			},
			IPACL:              IPACL,
			IsLocked:           false,
			NotificationEmails: []string{"mw+2@example.com"},
			PurgeOptions:       purgeOptions,
		}
	}

	expectGetAPIClient = func(client *iam.Mock, res iam.GetAPIClientResponse) {
		req := iam.GetAPIClientRequest{
			ClientID:    "1a2b3",
			GroupAccess: true,
			APIAccess:   true,
			IPACL:       true,
		}
		client.On("GetAPIClient", mock.Anything, req).Return(&res, nil).Once()
	}

	expectGetSelfAPIClient = func(client *iam.Mock, res iam.GetAPIClientResponse) {
		req := iam.GetAPIClientRequest{
			GroupAccess: true,
			APIAccess:   true,
			IPACL:       true,
		}
		client.On("GetAPIClient", mock.Anything, req).Return(&res, nil).Once()
	}

	getTfClient = func(ipAcl *IPACL, purgeOptions *PurgeOptions) TFClient {
		return TFClient{
			ClientID:           "1a2b3",
			AuthorizedUsers:    []string{"mw+2"},
			ClientType:         "CLIENT",
			ClientName:         "mw+2_1",
			NotificationEmails: []string{"mw+2@example.com"},
			ClientDescription:  "Test API Client",
			Lock:               false,
			GroupAccess: iam.GroupAccess{
				CloneAuthorizedUserGroups: false,
				Groups: []iam.ClientGroup{
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
				},
			},
			IPACL: ipAcl,
			APIAccess: iam.APIAccess{
				AllAccessibleAPIs: false,
				APIs: []iam.API{
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
		init     func(*iam.Mock, *templates.MockProcessor)
		clientID string
	}{
		"fetch API client with client ID": {
			init: func(i *iam.Mock, p *templates.MockProcessor) {
				expectGetAPIClient(i, getAPIClientResponse(&clientIPACL, &clientPurgeOptions))
				expectClientProcessTemplates(p, section, getTfClient(&clientTfIPACL, &clientTfPurgeOptions))
			},
			clientID: "1a2b3",
		},
		"fetch self API client": {
			init: func(i *iam.Mock, p *templates.MockProcessor) {
				expectGetSelfAPIClient(i, getAPIClientResponse(&clientIPACL, &clientPurgeOptions))
				expectClientProcessTemplates(p, section, getTfClient(&clientTfIPACL, &clientTfPurgeOptions))
			},
			clientID: "",
		},
		"fetch API client no IPACL": {
			init: func(i *iam.Mock, p *templates.MockProcessor) {
				expectGetAPIClient(i, getAPIClientResponse(nil, &clientPurgeOptions))
				expectClientProcessTemplates(p, section, getTfClient(nil, &clientTfPurgeOptions))
			},
			clientID: "1a2b3",
		},
		"fetch API client no PurgeOptions": {
			init: func(i *iam.Mock, p *templates.MockProcessor) {
				expectGetAPIClient(i, getAPIClientResponse(&clientIPACL, nil))
				expectClientProcessTemplates(p, section, getTfClient(&clientTfIPACL, nil))
			},
			clientID: "1a2b3",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mi := new(iam.Mock)
			mp := new(templates.MockProcessor)
			test.init(mi, mp)
			ctx := terminal.Context(context.Background(), terminal.New(terminal.DiscardWriter(), nil, terminal.DiscardWriter()))
			err := createIAMAPIClient(ctx, test.clientID, section, mi, mp)
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
					GroupAccess: iam.GroupAccess{
						CloneAuthorizedUserGroups: false,
						Groups: []iam.ClientGroup{
							{
								GroupID: 123,
								RoleID:  340,
								Subgroups: []iam.ClientGroup{
									{
										GroupID: 333,
										RoleID:  540,
									},
								},
							},
						},
					},
					IPACL: &IPACL{
						Enable: true,
						CIDR:   []string{"128.5.6.5/24"},
					},
					APIAccess: iam.APIAccess{
						AllAccessibleAPIs: false,
						APIs: []iam.API{
							{
								APIID:            5580,
								APIName:          "Search Data Feed",
								Description:      "Search Data Feed",
								Endpoint:         "/search-portal-data-feed-api/",
								DocumentationURL: "/",
								AccessLevel:      "READ-ONLY",
							},
							{
								APIID:            5801,
								APIName:          "EdgeWorkers",
								Description:      "EdgeWorkers",
								Endpoint:         "/edgeworkers/",
								DocumentationURL: "https://developer.akamai.com/api/web_performance/edgeworkers/v1.html",
								AccessLevel:      "READ-WRITE",
							},
						},
					},
					PurgeOptions: &PurgeOptions{
						CanPurgeByCPCode:   true,
						CanPurgeByCacheTag: true,
						CPCodeAccess: CPCodeAccess{
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
		"client with a few groups": {
			givenData: TFData{
				TFClient: TFClient{
					ClientID:           "1a2b3",
					AuthorizedUsers:    []string{"mw+2"},
					ClientType:         "CLIENT",
					ClientName:         "mw+2_1",
					NotificationEmails: []string{"mw+2@example.com"},
					ClientDescription:  "Test API Client",
					Lock:               false,
					GroupAccess: iam.GroupAccess{
						CloneAuthorizedUserGroups: false,
						Groups: []iam.ClientGroup{
							{
								GroupID: 123,
								RoleID:  340,
								Subgroups: []iam.ClientGroup{
									{
										GroupID: 333,
										RoleID:  540,
									},
								},
							},
							{
								GroupID: 345,
								RoleID:  340,
								Subgroups: []iam.ClientGroup{
									{
										GroupID: 333,
										RoleID:  540,
									},
								},
							},
						},
					},
					IPACL: &IPACL{
						Enable: true,
						CIDR:   []string{"128.5.6.5/24"},
					},
					APIAccess: iam.APIAccess{
						AllAccessibleAPIs: false,
						APIs: []iam.API{
							{
								APIID:            5580,
								APIName:          "Search Data Feed",
								Description:      "Search Data Feed",
								Endpoint:         "/search-portal-data-feed-api/",
								DocumentationURL: "/",
								AccessLevel:      "READ-ONLY",
							},
							{
								APIID:            5801,
								APIName:          "EdgeWorkers",
								Description:      "EdgeWorkers",
								Endpoint:         "/edgeworkers/",
								DocumentationURL: "https://developer.akamai.com/api/web_performance/edgeworkers/v1.html",
								AccessLevel:      "READ-WRITE",
							},
						},
					},
					PurgeOptions: &PurgeOptions{
						CanPurgeByCPCode:   true,
						CanPurgeByCacheTag: true,
						CPCodeAccess: CPCodeAccess{
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
					GroupAccess: iam.GroupAccess{
						CloneAuthorizedUserGroups: false,
						Groups: []iam.ClientGroup{
							{
								GroupID: 123,
								RoleID:  340,
							},
						},
					},
					IPACL: &IPACL{
						Enable: true,
						CIDR:   []string{"128.5.6.5/24"},
					},
					APIAccess: iam.APIAccess{
						AllAccessibleAPIs: false,
						APIs: []iam.API{
							{
								APIID:            5580,
								APIName:          "Search Data Feed",
								Description:      "Search Data Feed",
								Endpoint:         "/search-portal-data-feed-api/",
								DocumentationURL: "/",
								AccessLevel:      "READ-ONLY",
							},
							{
								APIID:            5801,
								APIName:          "EdgeWorkers",
								Description:      "EdgeWorkers",
								Endpoint:         "/edgeworkers/",
								DocumentationURL: "https://developer.akamai.com/api/web_performance/edgeworkers/v1.html",
								AccessLevel:      "READ-WRITE",
							},
						},
					},
					PurgeOptions: &PurgeOptions{
						CanPurgeByCPCode:   true,
						CanPurgeByCacheTag: true,
						CPCodeAccess: CPCodeAccess{
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
					GroupAccess: iam.GroupAccess{
						CloneAuthorizedUserGroups: false,
						Groups: []iam.ClientGroup{
							{
								GroupID: 123,
								RoleID:  340,
								Subgroups: []iam.ClientGroup{
									{
										GroupID: 333,
										RoleID:  540,
									},
								},
							},
						},
					},
					IPACL: &IPACL{
						Enable: true,
						CIDR:   []string{"128.5.6.5/24"},
					},
					APIAccess: iam.APIAccess{
						AllAccessibleAPIs: false,
						APIs: []iam.API{
							{
								APIID:            5580,
								APIName:          "Search Data Feed",
								Description:      "Search Data Feed",
								Endpoint:         "/search-portal-data-feed-api/",
								DocumentationURL: "/",
								AccessLevel:      "READ-ONLY",
							},
							{
								APIID:            5801,
								APIName:          "EdgeWorkers",
								Description:      "EdgeWorkers",
								Endpoint:         "/edgeworkers/",
								DocumentationURL: "https://developer.akamai.com/api/web_performance/edgeworkers/v1.html",
								AccessLevel:      "READ-WRITE",
							},
						},
					},
					PurgeOptions: &PurgeOptions{
						CanPurgeByCPCode:   true,
						CanPurgeByCacheTag: true,
						CPCodeAccess: CPCodeAccess{
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
					GroupAccess: iam.GroupAccess{
						CloneAuthorizedUserGroups: false,
						Groups: []iam.ClientGroup{
							{
								GroupID: 123,
								RoleID:  340,
								Subgroups: []iam.ClientGroup{
									{
										GroupID: 333,
										RoleID:  540,
									},
								},
							},
						},
					},
					IPACL: &IPACL{
						Enable: true,
					},
					APIAccess: iam.APIAccess{
						AllAccessibleAPIs: false,
						APIs: []iam.API{
							{
								APIID:            5580,
								APIName:          "Search Data Feed",
								Description:      "Search Data Feed",
								Endpoint:         "/search-portal-data-feed-api/",
								DocumentationURL: "/",
								AccessLevel:      "READ-ONLY",
							},
							{
								APIID:            5801,
								APIName:          "EdgeWorkers",
								Description:      "EdgeWorkers",
								Endpoint:         "/edgeworkers/",
								DocumentationURL: "https://developer.akamai.com/api/web_performance/edgeworkers/v1.html",
								AccessLevel:      "READ-WRITE",
							},
						},
					},
					PurgeOptions: &PurgeOptions{
						CanPurgeByCPCode:   true,
						CanPurgeByCacheTag: true,
						CPCodeAccess: CPCodeAccess{
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
					GroupAccess: iam.GroupAccess{
						CloneAuthorizedUserGroups: false,
						Groups: []iam.ClientGroup{
							{
								GroupID: 123,
								RoleID:  340,
								Subgroups: []iam.ClientGroup{
									{
										GroupID: 333,
										RoleID:  540,
										Subgroups: []iam.ClientGroup{
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
					IPACL: &IPACL{
						Enable: true,
						CIDR:   []string{"128.5.6.5/24"},
					},
					APIAccess: iam.APIAccess{
						AllAccessibleAPIs: false,
						APIs: []iam.API{
							{
								APIID:            5580,
								APIName:          "Search Data Feed",
								Description:      "Search Data Feed",
								Endpoint:         "/search-portal-data-feed-api/",
								DocumentationURL: "/",
								AccessLevel:      "READ-ONLY",
							},
							{
								APIID:            5801,
								APIName:          "EdgeWorkers",
								Description:      "EdgeWorkers",
								Endpoint:         "/edgeworkers/",
								DocumentationURL: "https://developer.akamai.com/api/web_performance/edgeworkers/v1.html",
								AccessLevel:      "READ-WRITE",
							},
						},
					},
					PurgeOptions: &PurgeOptions{
						CanPurgeByCPCode:   true,
						CanPurgeByCacheTag: true,
						CPCodeAccess: CPCodeAccess{
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
					GroupAccess: iam.GroupAccess{
						CloneAuthorizedUserGroups: true,
						Groups: []iam.ClientGroup{
							{
								GroupID: 123,
								RoleID:  340,
								Subgroups: []iam.ClientGroup{
									{
										GroupID: 333,
										RoleID:  540,
									},
								},
							},
						},
					},
					IPACL: &IPACL{
						Enable: true,
						CIDR:   []string{"128.5.6.5/24"},
					},
					APIAccess: iam.APIAccess{
						AllAccessibleAPIs: false,
						APIs: []iam.API{
							{
								APIID:            5580,
								APIName:          "Search Data Feed",
								Description:      "Search Data Feed",
								Endpoint:         "/search-portal-data-feed-api/",
								DocumentationURL: "/",
								AccessLevel:      "READ-ONLY",
							},
							{
								APIID:            5801,
								APIName:          "EdgeWorkers",
								Description:      "EdgeWorkers",
								Endpoint:         "/edgeworkers/",
								DocumentationURL: "https://developer.akamai.com/api/web_performance/edgeworkers/v1.html",
								AccessLevel:      "READ-WRITE",
							},
						},
					},
					PurgeOptions: &PurgeOptions{
						CanPurgeByCPCode:   true,
						CanPurgeByCacheTag: true,
						CPCodeAccess: CPCodeAccess{
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
					GroupAccess: iam.GroupAccess{
						CloneAuthorizedUserGroups: false,
						Groups: []iam.ClientGroup{
							{
								GroupID: 123,
								RoleID:  340,
								Subgroups: []iam.ClientGroup{
									{
										GroupID: 333,
										RoleID:  540,
									},
								},
							},
						},
					},
					IPACL: &IPACL{
						Enable: true,
						CIDR:   []string{"128.5.6.5/24"},
					},
					APIAccess: iam.APIAccess{
						AllAccessibleAPIs: true,
						APIs: []iam.API{
							{
								APIID:            5580,
								APIName:          "Search Data Feed",
								Description:      "Search Data Feed",
								Endpoint:         "/search-portal-data-feed-api/",
								DocumentationURL: "/",
								AccessLevel:      "READ-ONLY",
							},
							{
								APIID:            5801,
								APIName:          "EdgeWorkers",
								Description:      "EdgeWorkers",
								Endpoint:         "/edgeworkers/",
								DocumentationURL: "https://developer.akamai.com/api/web_performance/edgeworkers/v1.html",
								AccessLevel:      "READ-WRITE",
							},
						},
					},
					PurgeOptions: &PurgeOptions{
						CanPurgeByCPCode:   true,
						CanPurgeByCacheTag: true,
						CPCodeAccess: CPCodeAccess{
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
					GroupAccess: iam.GroupAccess{
						CloneAuthorizedUserGroups: false,
						Groups: []iam.ClientGroup{
							{
								GroupID: 123,
								RoleID:  340,
								Subgroups: []iam.ClientGroup{
									{
										GroupID: 333,
										RoleID:  540,
									},
								},
							},
						},
					},
					IPACL: &IPACL{
						Enable: true,
						CIDR:   []string{"128.5.6.5/24"},
					},
					APIAccess: iam.APIAccess{
						AllAccessibleAPIs: false,
						APIs: []iam.API{
							{
								APIID:            5580,
								APIName:          "Search Data Feed",
								Description:      "Search Data Feed",
								Endpoint:         "/search-portal-data-feed-api/",
								DocumentationURL: "/",
								AccessLevel:      "READ-ONLY",
							},
							{
								APIID:            5801,
								APIName:          "EdgeWorkers",
								Description:      "EdgeWorkers",
								Endpoint:         "/edgeworkers/",
								DocumentationURL: "https://developer.akamai.com/api/web_performance/edgeworkers/v1.html",
								AccessLevel:      "READ-WRITE",
							},
						},
					},
					PurgeOptions: &PurgeOptions{
						CanPurgeByCPCode:   true,
						CanPurgeByCacheTag: true,
						CPCodeAccess: CPCodeAccess{
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
					GroupAccess: iam.GroupAccess{
						CloneAuthorizedUserGroups: false,
						Groups: []iam.ClientGroup{
							{
								GroupID: 123,
								RoleID:  340,
								Subgroups: []iam.ClientGroup{
									{
										GroupID: 333,
										RoleID:  540,
									},
								},
							},
						},
					},
					APIAccess: iam.APIAccess{
						AllAccessibleAPIs: false,
						APIs: []iam.API{
							{
								APIID:            5580,
								APIName:          "Search Data Feed",
								Description:      "Search Data Feed",
								Endpoint:         "/search-portal-data-feed-api/",
								DocumentationURL: "/",
								AccessLevel:      "READ-ONLY",
							},
							{
								APIID:            5801,
								APIName:          "EdgeWorkers",
								Description:      "EdgeWorkers",
								Endpoint:         "/edgeworkers/",
								DocumentationURL: "https://developer.akamai.com/api/web_performance/edgeworkers/v1.html",
								AccessLevel:      "READ-WRITE",
							},
						},
					},
					PurgeOptions: &PurgeOptions{
						CanPurgeByCPCode:   true,
						CanPurgeByCacheTag: true,
						CPCodeAccess: CPCodeAccess{
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
					GroupAccess: iam.GroupAccess{
						CloneAuthorizedUserGroups: false,
						Groups: []iam.ClientGroup{
							{
								GroupID: 123,
								RoleID:  340,
								Subgroups: []iam.ClientGroup{
									{
										GroupID: 333,
										RoleID:  540,
									},
								},
							},
						},
					},
					IPACL: &IPACL{
						Enable: true,
						CIDR:   []string{"128.5.6.5/24"},
					},
					APIAccess: iam.APIAccess{
						AllAccessibleAPIs: false,
						APIs: []iam.API{
							{
								APIID:            5580,
								APIName:          "Search Data Feed",
								Description:      "Search Data Feed",
								Endpoint:         "/search-portal-data-feed-api/",
								DocumentationURL: "/",
								AccessLevel:      "READ-ONLY",
							},
							{
								APIID:            5801,
								APIName:          "EdgeWorkers",
								Description:      "EdgeWorkers",
								Endpoint:         "/edgeworkers/",
								DocumentationURL: "https://developer.akamai.com/api/web_performance/edgeworkers/v1.html",
								AccessLevel:      "READ-WRITE",
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
