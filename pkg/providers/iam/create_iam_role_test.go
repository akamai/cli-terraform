package iam

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/iam"
	"github.com/akamai/cli-terraform/pkg/templates"
	"github.com/akamai/cli-terraform/pkg/tools"
	"github.com/akamai/cli/pkg/terminal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	group1 = iam.Group{
		GroupID:       56789,
		ParentGroupID: 6789,
		GroupName:     "Custom group 1",
	}

	group2 = iam.Group{
		GroupID:       98765,
		ParentGroupID: 6789,
		GroupName:     "Custom group 2",
	}

	user1 = iam.User{
		UserBasicInfo: iam.UserBasicInfo{
			FirstName:         "John",
			LastName:          "Smith",
			Email:             "terraform@akamai.com",
			Country:           "Canada",
			Phone:             "(617) 444-4649",
			TFAEnabled:        true,
			ContactType:       "Technical Decision Maker",
			JobTitle:          "job title ",
			TimeZone:          "GMT",
			SecondaryEmail:    "secondary-email-a@akamai.net",
			MobilePhone:       "(617) 444-4649",
			Address:           "123 A Street",
			City:              "A-Town",
			State:             "TBD",
			ZipCode:           "34567",
			PreferredLanguage: "English",
			SessionTimeOut:    tools.IntPtr(900),
		},
		IdentityID: "123",
		IsLocked:   false,
		AuthGrants: []iam.AuthGrant{
			{
				RoleID:          tools.IntPtr(12345),
				RoleName:        "Custom role",
				RoleDescription: "Custom role description",
				GroupID:         group1.GroupID,
				GroupName:       group1.GroupName,
			},
		},
	}

	user2 = iam.User{
		UserBasicInfo: iam.UserBasicInfo{
			FirstName:         "Steve",
			LastName:          "Smith",
			Email:             "terraform_1@akamai.com",
			Country:           "Canada",
			Phone:             "(617) 444-4650",
			TFAEnabled:        true,
			ContactType:       "Technical Decision Maker",
			JobTitle:          "job title ",
			TimeZone:          "GMT",
			SecondaryEmail:    "secondary-email-a@akamai.net",
			MobilePhone:       "(617) 444-4649",
			Address:           "123 A Street",
			City:              "A-Town",
			State:             "TBD",
			ZipCode:           "34567",
			PreferredLanguage: "English",
			SessionTimeOut:    tools.IntPtr(900),
		},
		IdentityID: "321",
		IsLocked:   false,
		AuthGrants: []iam.AuthGrant{
			{
				RoleID:          tools.IntPtr(12345),
				RoleName:        "Custom role",
				RoleDescription: "Custom role description",
				GroupID:         group1.GroupID,
				GroupName:       group1.GroupName,
			},
			{
				RoleID:          tools.IntPtr(54321),
				RoleName:        "Other custom role",
				RoleDescription: "Other custom role description",
				GroupID:         group2.GroupID,
				GroupName:       group2.GroupName,
			},
		},
	}

	role = iam.Role{
		RoleID:          12345,
		RoleName:        "Custom role",
		RoleDescription: "Custom role description",
		Users: []iam.RoleUser{
			{
				UIIdentityID: user1.IdentityID,
				FirstName:    user1.FirstName,
				LastName:     user1.LastName,
				Email:        user1.Email,
			},
			{
				UIIdentityID: user2.IdentityID,
				FirstName:    user2.FirstName,
				LastName:     user2.LastName,
				Email:        user2.Email,
			},
		},
	}

	expectGetRoleWithUsers = func(client *iam.Mock) {
		getRoleReq := iam.GetRoleRequest{
			ID:           role.RoleID,
			GrantedRoles: true,
			Users:        true,
		}

		client.On("GetRole", mock.Anything, getRoleReq).Return(&role, nil).Once()
	}

	expectRoleGetUser = func(client *iam.Mock, user iam.User, err error) {
		getUserReq := iam.GetUserRequest{
			IdentityID:    user.IdentityID,
			Actions:       true,
			AuthGrants:    true,
			Notifications: true,
		}

		if err != nil {
			client.On("GetUser", mock.Anything, getUserReq).Return(nil, err).Once()
		}

		client.On("GetUser", mock.Anything, getUserReq).Return(&user, err).Once()
	}

	expectRoleGetGroup = func(client *iam.Mock, group iam.Group) {
		getGroupReq := iam.GetGroupRequest{
			GroupID: group.GroupID,
		}

		client.On("GetGroup", mock.Anything, getGroupReq).Return(&group, nil).Once()
	}

	expectRoleProcessTemplates = func(p *templates.MockProcessor, section string) *mock.Call {
		tfData := TFData{
			TFUsers: []*TFUser{
				{
					IsLocked:   false,
					AuthGrants: `[{"groupId":56789,"isBlocked":false,"roleId":12345}]`,
					TFUserBasicInfo: TFUserBasicInfo{
						ID:                "123",
						FirstName:         "John",
						LastName:          "Smith",
						Email:             "terraform@akamai.com",
						Country:           "Canada",
						Phone:             "(617) 444-4649",
						TFAEnabled:        true,
						ContactType:       "Technical Decision Maker",
						JobTitle:          "job title ",
						TimeZone:          "GMT",
						SecondaryEmail:    "secondary-email-a@akamai.net",
						MobilePhone:       "(617) 444-4649",
						Address:           "123 A Street",
						City:              "A-Town",
						State:             "TBD",
						ZipCode:           "34567",
						PreferredLanguage: "English",
						SessionTimeOut:    tools.IntPtr(900),
					},
				},
				{
					IsLocked:   false,
					AuthGrants: `[{"groupId":56789,"isBlocked":false,"roleId":12345},{"groupId":98765,"isBlocked":false,"roleId":54321}]`,
					TFUserBasicInfo: TFUserBasicInfo{
						ID:                "321",
						FirstName:         "Steve",
						LastName:          "Smith",
						Email:             "terraform_1@akamai.com",
						Country:           "Canada",
						Phone:             "(617) 444-4650",
						TFAEnabled:        true,
						ContactType:       "Technical Decision Maker",
						JobTitle:          "job title ",
						TimeZone:          "GMT",
						SecondaryEmail:    "secondary-email-a@akamai.net",
						MobilePhone:       "(617) 444-4649",
						Address:           "123 A Street",
						City:              "A-Town",
						State:             "TBD",
						ZipCode:           "34567",
						PreferredLanguage: "English",
						SessionTimeOut:    tools.IntPtr(900),
					},
				},
			},
			TFRoles: []TFRole{
				{
					RoleID:          12345,
					RoleName:        "Custom role",
					RoleDescription: "Custom role description",
					GrantedRoles:    []int{},
				},
			},
			TFGroups: []TFGroup{
				{
					GroupID:       56789,
					ParentGroupID: 6789,
					GroupName:     "Custom group 1",
				},
				{
					GroupID:       98765,
					ParentGroupID: 6789,
					GroupName:     "Custom group 2",
				},
			},
			Section:    section,
			Subcommand: "role",
		}
		call := p.On(
			"ProcessTemplates",
			tfData,
		)
		return call.Return(nil)
	}
)

func TestGetUsersByRole(t *testing.T) {
	tests := map[string]struct {
		roleUsers    []iam.RoleUser
		init         func(*iam.Mock, *terminal.Mock)
		expectResult []*iam.User
		withError    error
	}{
		"ok no errors": {
			roleUsers: []iam.RoleUser{
				{
					UIIdentityID: "a",
				},
			},
			expectResult: []*iam.User{
				{IdentityID: "a", AuthGrants: []iam.AuthGrant{
					{
						GroupID: 12345,
						RoleID:  tools.IntPtr(54321),
					},
				}},
			},
			init: func(m *iam.Mock, t *terminal.Mock) {
				m.On("GetUser", mock.Anything, iam.GetUserRequest{
					IdentityID:    "a",
					Actions:       true,
					AuthGrants:    true,
					Notifications: true,
				}).Return(&iam.User{IdentityID: "a", AuthGrants: []iam.AuthGrant{
					{
						GroupID: 12345,
						RoleID:  tools.IntPtr(54321),
					},
				}}, nil).Once()
			},
		},
		"fail first": {
			roleUsers: []iam.RoleUser{
				{
					UIIdentityID: "a",
				},
				{
					UIIdentityID: "b",
				},
			},
			expectResult: []*iam.User{
				{IdentityID: "b", AuthGrants: []iam.AuthGrant{
					{
						GroupID: 12345,
						RoleID:  tools.IntPtr(54321),
					},
				}},
			},
			init: func(m *iam.Mock, t *terminal.Mock) {
				m.On("GetUser", mock.Anything, iam.GetUserRequest{
					IdentityID:    "a",
					Actions:       true,
					AuthGrants:    true,
					Notifications: true,
				}).Return(nil, fmt.Errorf("an error")).Once()
				t.On("Writeln", []interface{}{"[WARN] Unable to fetch user of ID 'a' - skipping:\nan error"}).Return(0, nil).Once()
				m.On("GetUser", mock.Anything, iam.GetUserRequest{
					IdentityID:    "b",
					Actions:       true,
					AuthGrants:    true,
					Notifications: true,
				}).Return(&iam.User{IdentityID: "b", AuthGrants: []iam.AuthGrant{
					{
						GroupID: 12345,
						RoleID:  tools.IntPtr(54321),
					},
				}}, nil).Once()
			},
		},
		"fail all": {
			roleUsers: []iam.RoleUser{
				{
					UIIdentityID: "a",
				},
				{
					UIIdentityID: "b",
				},
			},
			expectResult: []*iam.User{},
			init: func(m *iam.Mock, t *terminal.Mock) {
				m.On("GetUser", mock.Anything, iam.GetUserRequest{
					IdentityID:    "a",
					Actions:       true,
					AuthGrants:    true,
					Notifications: true,
				}).Return(nil, fmt.Errorf("an error")).Once()
				t.On("Writeln", []interface{}{"[WARN] Unable to fetch user of ID 'a' - skipping:\nan error"}).Return(0, nil).Once()
				m.On("GetUser", mock.Anything, iam.GetUserRequest{
					IdentityID:    "b",
					Actions:       true,
					AuthGrants:    true,
					Notifications: true,
				}).Return(nil, fmt.Errorf("another error")).Once()
				t.On("Writeln", []interface{}{"[WARN] Unable to fetch user of ID 'b' - skipping:\nanother error"}).Return(0, nil).Once()
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			client := iam.Mock{}
			term := terminal.Mock{}
			test.init(&client, &term)

			result, err := getUsersByRole(context.Background(), &term, test.roleUsers, &client)
			if test.withError != nil {
				require.Error(t, err)
				assert.True(t, errors.Is(err, test.withError))
				return
			}
			require.NoError(t, err)
			assert.Equal(t, test.expectResult, result)
		})
	}
}

func TestCreateIAMRole(t *testing.T) {
	section := "test_section"

	tests := map[string]struct {
		init func(*iam.Mock, *templates.MockProcessor)
	}{
		"fetch role": {
			init: func(i *iam.Mock, p *templates.MockProcessor) {
				expectGetRoleWithUsers(i)
				expectRoleGetUser(i, user1, nil)
				expectRoleGetUser(i, user2, nil)
				// user1's groups
				expectRoleGetGroup(i, group1)
				// user2's groups
				expectRoleGetGroup(i, group1)
				expectRoleGetGroup(i, group2)
				expectRoleProcessTemplates(p, section)
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mi := new(iam.Mock)
			mp := new(templates.MockProcessor)
			test.init(mi, mp)
			ctx := terminal.Context(context.Background(), terminal.New(terminal.DiscardWriter(), nil, terminal.DiscardWriter()))
			err := createIAMRoleByID(ctx, 12345, section, mi, mp)
			require.NoError(t, err)
			mi.AssertExpectations(t)
			mp.AssertExpectations(t)
		})
	}
}

func TestProcessIAMRoleTemplates(t *testing.T) {
	section := "test_section"

	tests := map[string]struct {
		givenData    TFData
		dir          string
		filesToCheck []string
	}{
		"basic role": {
			givenData: TFData{
				TFUsers: []*TFUser{
					{
						IsLocked:   false,
						AuthGrants: `[{"groupId":56789,"isBlocked":false,"roleId":12345}]`,
						TFUserBasicInfo: TFUserBasicInfo{
							ID:                "123",
							FirstName:         "John",
							LastName:          "Smith",
							Email:             "terraform@akamai.com",
							Country:           "Canada",
							Phone:             "(617) 444-4649",
							TFAEnabled:        true,
							ContactType:       "Technical Decision Maker",
							JobTitle:          "job title ",
							TimeZone:          "GMT",
							SecondaryEmail:    "secondary-email-a@akamai.net",
							MobilePhone:       "(617) 444-4649",
							Address:           "123 A Street",
							City:              "A-Town",
							State:             "TBD",
							ZipCode:           "34567",
							PreferredLanguage: "English",
							SessionTimeOut:    tools.IntPtr(900),
						},
					},
				},
				TFRoles: []TFRole{
					{
						RoleID:          12345,
						RoleName:        "Custom role",
						RoleDescription: "Custom role\ndescription",
						GrantedRoles:    []int{992, 707, 452, 677, 726, 296, 457, 987},
					},
				},
				TFGroups: []TFGroup{
					{
						GroupID:       56789,
						ParentGroupID: 98765,
						GroupName:     "Custom group 1",
					},
				},
				Section:    section,
				Subcommand: "role",
			},
			dir:          "iam_role_by_id_basic",
			filesToCheck: []string{"role.tf", "variables.tf", "import.sh", "users.tf", "groups.tf"},
		},
		"role with multiple users and groups": {
			givenData: TFData{
				TFUsers: []*TFUser{
					{
						IsLocked:   false,
						AuthGrants: `[{"groupId":56789,"isBlocked":false,"roleId":12345}]`,
						TFUserBasicInfo: TFUserBasicInfo{
							ID:                "123",
							FirstName:         "John",
							LastName:          "Smith",
							Email:             "terraform@akamai.com",
							Country:           "Canada",
							Phone:             "(617) 444-4649",
							TFAEnabled:        true,
							ContactType:       "Technical Decision Maker",
							JobTitle:          "job title ",
							TimeZone:          "GMT",
							SecondaryEmail:    "secondary-email-a@akamai.net",
							MobilePhone:       "(617) 444-4649",
							Address:           "123 A Street",
							City:              "A-Town",
							State:             "TBD",
							ZipCode:           "34567",
							PreferredLanguage: "English",
							SessionTimeOut:    tools.IntPtr(900),
						},
					},
					{
						IsLocked:   false,
						AuthGrants: `[{"groupId":56789,"isBlocked":false,"roleId":12345},{"groupId":98765,"isBlocked":false,"roleId":54321}]`,
						TFUserBasicInfo: TFUserBasicInfo{
							ID:                "321",
							FirstName:         "Steve",
							LastName:          "Smith",
							Email:             "terraform_1@akamai.com",
							Country:           "Canada",
							Phone:             "(617) 444-4650",
							TFAEnabled:        true,
							ContactType:       "Technical Decision Maker",
							JobTitle:          "job title ",
							TimeZone:          "GMT",
							SecondaryEmail:    "secondary-email-a@akamai.net",
							MobilePhone:       "(617) 444-4649",
							Address:           "123 A Street",
							City:              "A-Town",
							State:             "TBD",
							ZipCode:           "34567",
							PreferredLanguage: "English",
							SessionTimeOut:    tools.IntPtr(900),
						},
					},
				},
				TFRoles: []TFRole{
					{
						RoleID:          12345,
						RoleName:        "Custom role",
						RoleDescription: "Custom role description",
						GrantedRoles:    []int{992, 707, 452, 677, 726, 296, 457, 987},
					},
				},
				TFGroups: []TFGroup{
					{
						GroupID:       56789,
						ParentGroupID: 12345,
						GroupName:     "Custom group 1",
					},
					{
						GroupID:       98765,
						ParentGroupID: 12345,
						GroupName:     "Custom group 2",
					},
				},
				Section:    section,
				Subcommand: "role",
			},
			dir:          "iam_role_by_id_multiple",
			filesToCheck: []string{"role.tf", "variables.tf", "import.sh", "users.tf", "groups.tf"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			require.NoError(t, os.MkdirAll(fmt.Sprintf("./testdata/res/%s", test.dir), 0755))
			processor := templates.FSTemplateProcessor{
				TemplatesFS: templateFiles,
				TemplateTargets: map[string]string{
					"groups.tmpl":    fmt.Sprintf("./testdata/res/%s/groups.tf", test.dir),
					"imports.tmpl":   fmt.Sprintf("./testdata/res/%s/import.sh", test.dir),
					"roles.tmpl":     fmt.Sprintf("./testdata/res/%s/role.tf", test.dir),
					"users.tmpl":     fmt.Sprintf("./testdata/res/%s/users.tf", test.dir),
					"variables.tmpl": fmt.Sprintf("./testdata/res/%s/variables.tf", test.dir),
				},
				AdditionalFuncs: additionalFunctions,
			}
			require.NoError(t, processor.ProcessTemplates(test.givenData))

			for _, f := range test.filesToCheck {
				expected, err := ioutil.ReadFile(fmt.Sprintf("./testdata/%s/%s", test.dir, f))
				require.NoError(t, err)
				result, err := ioutil.ReadFile(fmt.Sprintf("./testdata/res/%s/%s", test.dir, f))
				require.NoError(t, err)
				assert.Equal(t, string(expected), string(result))
			}
		})
	}
}
