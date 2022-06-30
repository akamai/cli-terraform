package iam

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v2/pkg/iam"
	"github.com/akamai/cli-terraform/pkg/templates"
	"github.com/akamai/cli-terraform/pkg/tools"
	"github.com/akamai/cli/pkg/terminal"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

var (
	expectListAllUsers = func(client *mockiam) {
		listUserReq := iam.ListUsersRequest{}

		users := []iam.UserListItem{
			{
				IdentityID: "001",
				Email:      "001@akamai.com",
			},
			{
				IdentityID: "002",
				Email:      "002@akamai.com",
			},
		}

		client.On("ListUsers", mock.Anything, listUserReq).Return(users, nil).Once()
	}

	expectGetUser001 = func(client *mockiam) {
		getUserReq := iam.GetUserRequest{
			IdentityID:    "001",
			Actions:       true,
			AuthGrants:    true,
			Notifications: true,
		}

		user := iam.User{
			UserBasicInfo: getUserBasicInfo(),
			IdentityID:    "001",
			IsLocked:      false,
			AuthGrants: []iam.AuthGrant{
				{
					RoleID:          tools.IntPtr(201),
					RoleName:        "role_201",
					RoleDescription: "role 201 description",
					GroupID:         101,
					GroupName:       "grp_101",
				},
			},
		}

		client.On("GetUser", mock.Anything, getUserReq).Return(&user, nil).Once()
	}
	expectGetUser002 = func(client *mockiam) {
		getUserReq := iam.GetUserRequest{
			IdentityID:    "002",
			Actions:       true,
			AuthGrants:    true,
			Notifications: true,
		}

		user := iam.User{
			UserBasicInfo: getUserBasicInfo(),
			IdentityID:    "002",
			IsLocked:      false,
			AuthGrants:    []iam.AuthGrant{},
		}

		client.On("GetUser", mock.Anything, getUserReq).Return(&user, nil).Once()
	}

	expectListAllGroups = func(client *mockiam) {
		listGroupsReq := iam.ListGroupsRequest{Actions: true}

		groups := []iam.Group{
			{
				GroupID:       101,
				GroupName:     "grp_101",
				ParentGroupID: 111,
			},
			{
				GroupID:       102,
				GroupName:     "grp_102",
				ParentGroupID: 111,
			},
		}

		client.On("ListGroups", mock.Anything, listGroupsReq).Return(groups, nil).Once()
	}

	expectListAllRoles = func(client *mockiam) {
		listRolesReq := iam.ListRolesRequest{
			Actions:       true,
			IgnoreContext: true,
			Users:         true,
		}

		roles := []iam.Role{
			{
				RoleID:          201,
				RoleName:        "role_201",
				RoleDescription: "role 201 description",
			},
			{
				RoleID:          202,
				RoleName:        "role_202",
				RoleDescription: "role 202 description",
			},
		}

		client.On("ListRoles", mock.Anything, listRolesReq).Return(roles, nil).Once()
	}

	expectAllProcessTemplates = func(p *mockProcessor, section string) *mock.Call {

		call := p.On(
			"ProcessTemplates",
			getTestData(section),
		)
		return call.Return(nil)
	}
)

func TestCreateIAMAll(t *testing.T) {
	section := "test_section"

	tests := map[string]struct {
		init func(*mockiam, *mockProcessor)
		err  error
	}{
		"fetch user": {
			init: func(i *mockiam, p *mockProcessor) {
				expectListAllUsers(i)
				expectGetUser001(i)
				expectGetUser002(i)
				expectListAllGroups(i)
				expectListAllRoles(i)
				expectAllProcessTemplates(p, section)
			},
		},

		"fail list users": {
			init: func(i *mockiam, _ *mockProcessor) {
				i.On("ListUsers", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("oops")).Once()
			},
			err: ErrFetchingUsers,
		},

		"fail get user": {
			init: func(i *mockiam, _ *mockProcessor) {
				expectListAllUsers(i)
				i.On("GetUser", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("oops")).Once()
			},
			err: ErrFetchingUserByID,
		},

		"fail list groups": {
			init: func(i *mockiam, _ *mockProcessor) {
				expectListAllUsers(i)
				expectGetUser001(i)
				expectGetUser002(i)
				i.On("ListGroups", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("oops")).Once()
			},
			err: ErrFetchingGroups,
		},

		"fail list roles": {
			init: func(i *mockiam, _ *mockProcessor) {
				expectListAllUsers(i)
				expectGetUser001(i)
				expectGetUser002(i)
				expectListAllGroups(i)
				i.On("ListRoles", mock.Anything, mock.Anything).Return(nil, fmt.Errorf("oops")).Once()
			},
			err: ErrFetchingGroups,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mi := new(mockiam)
			mp := new(mockProcessor)
			test.init(mi, mp)
			ctx := terminal.Context(context.Background(), terminal.New(terminal.DiscardWriter(), nil, terminal.DiscardWriter()))
			err := createIAMAll(ctx, section, mi, mp)
			if test.err != nil {
				errors.Is(err, test.err)
			} else {
				require.NoError(t, err)
			}
			mi.AssertExpectations(t)
			mp.AssertExpectations(t)
		})
	}
}

func TestProcessIAMAllTemplates(t *testing.T) {
	section := "test_section"

	tests := map[string]struct {
		givenData    TFData
		dir          string
		filesToCheck []string
	}{
		"one used, one not": {
			givenData:    getTestData(section),
			dir:          "iam_all",
			filesToCheck: []string{"users.tf", "variables.tf", "import.sh", "roles.tf", "groups.tf"},
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
					"users.tmpl":     fmt.Sprintf("./testdata/res/%s/users.tf", test.dir),
					"variables.tmpl": fmt.Sprintf("./testdata/res/%s/variables.tf", test.dir),
				},
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

func getTestData(section string) TFData {
	tfData := TFData{
		TFUsers: []*TFUser{
			{
				IsLocked:        false,
				AuthGrants:      "[{\"groupId\":101,\"isBlocked\":false,\"roleId\":201}]",
				TFUserBasicInfo: getTFUserBasicInfo(),
			},
			{
				IsLocked:        false,
				AuthGrants:      "",
				TFUserBasicInfo: getTFUserBasicInfo(),
			},
		},
		TFGroups: []TFGroup{
			{
				GroupID:       101,
				ParentGroupID: 111,
				GroupName:     "grp_101",
			},
			{
				GroupID:       102,
				ParentGroupID: 111,
				GroupName:     "grp_102",
			},
		},
		TFRoles: []TFRole{
			{
				RoleID:          201,
				RoleName:        "role_201",
				RoleDescription: "role 201 description",
				GrantedRoles:    []int{},
			},
			{
				RoleID:          202,
				RoleName:        "role_202",
				RoleDescription: "role 202 description",
				GrantedRoles:    []int{},
			},
		},
		Section:    section,
		Subcommand: "all",
	}
	tfData.TFUsers[0].ID = "001"
	tfData.TFUsers[1].ID = "002"
	return tfData
}
