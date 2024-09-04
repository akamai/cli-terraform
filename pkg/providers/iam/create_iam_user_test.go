package iam

import (
	"context"
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
	expectListUsers = func(client *iam.Mock) {
		listUserReq := iam.ListUsersRequest{}

		users := []iam.UserListItem{
			{
				IdentityID: "123",
				Email:      "terraform@akamai.com",
			},
			{
				IdentityID: "321",
				Email:      "terraform_1@akamai.com",
			},
		}

		client.On("ListUsers", mock.Anything, listUserReq).Return(users, nil).Once()
	}

	expectGetUser = func(client *iam.Mock) {
		getUserReq := iam.GetUserRequest{
			IdentityID:    "123",
			Actions:       true,
			AuthGrants:    true,
			Notifications: true,
		}

		user := iam.User{
			UserBasicInfo: getUserBasicInfo(),
			IdentityID:    "123",
			IsLocked:      false,
			AuthGrants: []iam.AuthGrant{
				{
					RoleID:          tools.IntPtr(12345),
					RoleName:        "Custom role",
					RoleDescription: "Custom role description",
					GroupID:         56789,
					GroupName:       "Custom group",
				},
			},
		}

		client.On("GetUser", mock.Anything, getUserReq).Return(&user, nil).Once()
	}

	expectGetUserRole = func(client *iam.Mock) {
		getRoleReq := iam.GetRoleRequest{
			ID:           12345,
			GrantedRoles: true,
		}
		role := iam.Role{
			RoleID:          12345,
			RoleName:        "Custom role",
			RoleDescription: "Custom role description",
		}

		client.On("GetRole", mock.Anything, getRoleReq).Return(&role, nil).Once()
	}

	expectGetUserGroup = func(client *iam.Mock) {
		getGroupReq1 := iam.GetGroupRequest{
			GroupID: 56789,
		}
		group1 := iam.Group{
			GroupID:       56789,
			ParentGroupID: 98765,
			GroupName:     "Custom group",
			SubGroups:     []iam.Group{{GroupID: 56473, GroupName: "the subgroup", ParentGroupID: 56789}},
		}
		client.On("GetGroup", mock.Anything, getGroupReq1).Return(&group1, nil).Once()

		getGroupReq2 := iam.GetGroupRequest{
			GroupID: 56473,
		}
		group2 := iam.Group{
			GroupID:       56473,
			ParentGroupID: 56789,
			GroupName:     "the subgroup",
		}
		client.On("GetGroup", mock.Anything, getGroupReq2).Return(&group2, nil).Once()

	}

	expectProcessTemplates = func(p *templates.MockProcessor, section string) *mock.Call {
		tfData := TFData{
			TFUsers: []*TFUser{
				{
					IsLocked:        false,
					AuthGrants:      "[{\"groupId\":56789,\"isBlocked\":false,\"roleId\":12345}]",
					TFUserBasicInfo: getTFUserBasicInfo(),
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
					ParentGroupID: 98765,
					GroupName:     "Custom group",
				},
				{
					GroupID:       56473,
					ParentGroupID: 56789,
					GroupName:     "the subgroup",
				},
			},
			Section:    section,
			Subcommand: "user",
		}
		call := p.On(
			"ProcessTemplates",
			tfData,
		)
		return call.Return(nil)
	}
)

func TestCreateIAMUserByEmail(t *testing.T) {
	section := "test_section"

	tests := map[string]struct {
		init func(*iam.Mock, *templates.MockProcessor)
	}{
		"fetch user": {
			init: func(i *iam.Mock, p *templates.MockProcessor) {
				expectListUsers(i)
				expectGetUser(i)
				expectGetUserRole(i)
				expectGetUserGroup(i)
				expectProcessTemplates(p, section)
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mi := new(iam.Mock)
			mp := new(templates.MockProcessor)
			test.init(mi, mp)
			ctx := terminal.Context(context.Background(), terminal.New(terminal.DiscardWriter(), nil, terminal.DiscardWriter()))
			err := createIAMUserByEmail(ctx, "terraform@akamai.com", section, mi, mp)
			require.NoError(t, err)
			mi.AssertExpectations(t)
			mp.AssertExpectations(t)
		})
	}
}

func TestProcessIAMUserTemplates(t *testing.T) {
	section := "test_section"

	tests := map[string]struct {
		givenData    TFData
		dir          string
		filesToCheck []string
	}{
		"basic user": {
			givenData: TFData{
				TFUsers: []*TFUser{
					{
						TFUserBasicInfo: getTFUserBasicInfo(),
						IsLocked:        false,
						AuthGrants:      "[{\"groupId\":56789,\"groupName\":\"Custom group\",\"isBlocked\":false,\"roleId\":12345}]",
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
						ParentGroupID: 98765,
						GroupName:     "Custom group",
					},
				},
				Section:    section,
				Subcommand: "user",
			},
			dir:          "iam_user_by_email_basic",
			filesToCheck: []string{"user.tf", "variables.tf", "import.sh", "roles.tf", "groups.tf"},
		},
		"basic user with newlines": {
			givenData: TFData{
				TFUsers: []*TFUser{
					{
						TFUserBasicInfo: getTFUserBasicInfoWithNewlines(),
						IsLocked:        false,
						AuthGrants:      "[{\"groupId\":56789,\"groupName\":\"Custom group\",\"isBlocked\":false,\"roleId\":12345}]",
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
						ParentGroupID: 98765,
						GroupName:     "Custom group",
					},
				},
				Section:    section,
				Subcommand: "user",
			},
			dir:          "iam_user_by_email_newlines",
			filesToCheck: []string{"user.tf", "variables.tf", "import.sh", "roles.tf", "groups.tf"},
		},
		"basic user with newline at the end of address": {
			givenData: TFData{
				TFUsers: []*TFUser{
					{
						TFUserBasicInfo: func() TFUserBasicInfo {
							info := getTFUserBasicInfo()
							info.Address = info.Address + "\n"
							return info

						}(),
						IsLocked:   false,
						AuthGrants: "[{\"groupId\":56789,\"groupName\":\"Custom group\",\"isBlocked\":false,\"roleId\":12345}]",
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
						ParentGroupID: 98765,
						GroupName:     "Custom group",
					},
				},
				Section:    section,
				Subcommand: "user",
			},
			dir:          "iam_user_by_email_end_newline",
			filesToCheck: []string{"user.tf", "variables.tf", "import.sh", "roles.tf", "groups.tf"},
		},
		"user with multiple auth grants": {
			givenData: TFData{
				TFUsers: []*TFUser{
					{
						TFUserBasicInfo: getTFUserBasicInfo(),
						IsLocked:        false,
						AuthGrants:      "[{\"groupId\":56789,\"isBlocked\":false,\"roleId\":12345},{\"groupId\":987,\"isBlocked\":false,\"roleId\":54321}]",
					},
				},
				TFRoles: []TFRole{
					{
						RoleID:          12345,
						RoleName:        "Custom role 12345",
						RoleDescription: "Custom role description",
						GrantedRoles:    []int{992, 707, 452, 677, 726, 296, 457, 987},
					},
					{
						RoleID:          54321,
						RoleName:        "Custom role 54321",
						RoleDescription: "Custom role description",
						GrantedRoles:    []int{992, 707, 452, 677, 726, 296, 457, 987},
					},
				},
				TFGroups: []TFGroup{
					{
						GroupID:       56789,
						ParentGroupID: 98765,
						GroupName:     "Custom group 56789",
					},
					{
						GroupID:       987,
						ParentGroupID: 98765,
						GroupName:     "Custom group 987",
					},
				},
				Section:    section,
				Subcommand: "user",
			},
			dir:          "iam_user_by_email_multiple_auth_grants",
			filesToCheck: []string{"groups.tf", "import.sh", "roles.tf", "user.tf", "variables.tf"},
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
					"roles.tmpl":     fmt.Sprintf("./testdata/res/%s/roles.tf", test.dir),
					"users.tmpl":     fmt.Sprintf("./testdata/res/%s/user.tf", test.dir),
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

func getUserBasicInfo() iam.UserBasicInfo {
	return iam.UserBasicInfo{
		FirstName:         "Terraform",
		LastName:          "Test",
		Email:             "terraform8@akamai.com",
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
	}
}

func getTFUserBasicInfo() TFUserBasicInfo {
	return TFUserBasicInfo{
		ID:                "123",
		FirstName:         "Terraform",
		LastName:          "Test",
		Email:             "terraform8@akamai.com",
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
	}
}

func getTFUserBasicInfoWithNewlines() TFUserBasicInfo {
	return TFUserBasicInfo{
		ID:                "123",
		FirstName:         "Terraform\n newline",
		LastName:          "Test\n newline",
		Email:             "terraform8@akamai.com",
		Country:           "Canada",
		Phone:             "(617) 444-4649",
		TFAEnabled:        true,
		ContactType:       "Technical Decision Maker",
		JobTitle:          "job\n title ",
		TimeZone:          "GMT",
		SecondaryEmail:    "secondary-email-a@akamai.net",
		MobilePhone:       "(617) 444-4649",
		Address:           "123\nA\nStreet",
		City:              "A-Town",
		State:             "TBD",
		ZipCode:           "34567",
		PreferredLanguage: "English",
		SessionTimeOut:    tools.IntPtr(900),
	}
}

func TestGetUserAuthGrants(t *testing.T) {
	tests := map[string]struct {
		user               iam.User
		expectedAuthGrants string
	}{
		"basic user": {
			user: iam.User{
				UserBasicInfo: getUserBasicInfo(),
				AuthGrants: []iam.AuthGrant{
					{
						GroupID:         12345,
						GroupName:       "Group name",
						IsBlocked:       false,
						RoleDescription: "Role description",
						RoleID:          tools.IntPtr(54321),
						RoleName:        "Custom role",
					},
				},
			},
			expectedAuthGrants: "[{\"groupId\":12345,\"isBlocked\":false,\"roleId\":54321}]",
		},
		"basic user with multiple auth grants": {
			user: iam.User{
				UserBasicInfo: getUserBasicInfo(),
				AuthGrants: []iam.AuthGrant{
					{
						GroupID:         56789,
						GroupName:       "Custom group 56789",
						IsBlocked:       false,
						RoleDescription: "Custom role description",
						RoleID:          tools.IntPtr(12345),
						RoleName:        "Custom role 12345",
					},
					{
						GroupID:         987,
						GroupName:       "Custom group 987",
						IsBlocked:       false,
						RoleDescription: "Custom role description",
						RoleID:          tools.IntPtr(54321),
						RoleName:        "Custom role 54321",
					},
				},
			},
			expectedAuthGrants: "[{\"groupId\":56789,\"isBlocked\":false,\"roleId\":12345}," +
				"{\"groupId\":987,\"isBlocked\":false,\"roleId\":54321}]",
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actualAuthGrants, err := getUserAuthGrants(&test.user)
			require.NoError(t, err)
			assert.Equal(t, test.expectedAuthGrants, actualAuthGrants)
		})
	}
}
