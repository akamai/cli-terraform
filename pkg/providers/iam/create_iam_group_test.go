package iam

import (
	"context"
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
	groupID = 56789

	expectListUsersWithinGroup = func(client *mockiam) {
		listUserReq := iam.ListUsersRequest{
			Actions: true,
			GroupID: tools.IntPtr(groupID),
		}

		users := []iam.UserListItem{
			{
				IdentityID: "123",
				Email:      "terraform@akamai.com",
				Actions:    &iam.UserActions{EditProfile: true},
			},
		}

		client.On("ListUsers", mock.Anything, listUserReq).Return(users, nil).Once()
	}

	expectGetUserWithinGroup = func(client *mockiam) {
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

	expectListRolesWithinGroup = func(client *mockiam) {
		listRolesReq := iam.ListRolesRequest{
			GroupID: tools.Int64Ptr(int64(groupID)),
		}

		roles := []iam.Role{
			{
				RoleID:          12345,
				RoleName:        "Custom role",
				RoleDescription: "Custom role description",
			},
		}
		client.On("ListRoles", mock.Anything, listRolesReq).Return(roles, nil).Once()
	}

	expectGetGroupWithinRole = func(client *mockiam) {
		getGroupReq := iam.GetGroupRequest{
			GroupID: 56789,
		}
		group := iam.Group{
			GroupID:       56789,
			ParentGroupID: 98765,
			GroupName:     "Custom group",
		}
		client.On("GetGroup", mock.Anything, getGroupReq).Return(&group, nil).Once()
	}

	expectGroupByIDProcessTemplates = func(p *mockProcessor, section string) *mock.Call {
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
			},
			Section:    section,
			Subcommand: "group",
		}
		call := p.On(
			"ProcessTemplates",
			tfData,
		)
		return call.Return(nil)
	}
)

func TestCreateIAMGroupByID(t *testing.T) {
	section := "test_section"

	tests := map[string]struct {
		init func(*mockiam, *mockProcessor)
	}{
		"fetch user": {
			init: func(i *mockiam, p *mockProcessor) {
				expectListUsersWithinGroup(i)
				expectGetUserWithinGroup(i)
				expectListRolesWithinGroup(i)
				expectGetGroupWithinRole(i)
				expectGroupByIDProcessTemplates(p, section)
			},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			mi := new(mockiam)
			mp := new(mockProcessor)
			test.init(mi, mp)
			ctx := terminal.Context(context.Background(), terminal.New(terminal.DiscardWriter(), nil, terminal.DiscardWriter()))
			err := createIAMGroupByID(ctx, groupID, section, mi, mp)
			require.NoError(t, err)
			mi.AssertExpectations(t)
			mp.AssertExpectations(t)
		})
	}
}

func TestProcessIAMGroupTemplates(t *testing.T) {
	section := "test_section"

	tests := map[string]struct {
		givenData    TFData
		dir          string
		filesToCheck []string
	}{
		"basic group": {
			givenData: TFData{
				TFUsers: []*TFUser{
					{
						TFUserBasicInfo: getTFUserBasicInfo(),
						IsLocked:        false,
						AuthGrants:      "[{\"groupId\":56789,\"groupName\":\"Custom group\",\"isBlocked\":false,\"roleId\":12345}]",
					},
				},
				TFGroups: []TFGroup{
					{
						GroupID:       56789,
						ParentGroupID: 98765,
						GroupName:     "Custom group",
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
				Section:    section,
				Subcommand: "group",
			},
			dir:          "iam_group_by_id",
			filesToCheck: []string{"users.tf", "variables.tf", "import.sh", "roles.tf", "group.tf"},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			require.NoError(t, os.MkdirAll(fmt.Sprintf("./testdata/res/%s", test.dir), 0755))
			processor := templates.FSTemplateProcessor{
				TemplatesFS: templateFiles,
				TemplateTargets: map[string]string{
					"groups.tmpl":    fmt.Sprintf("./testdata/res/%s/group.tf", test.dir),
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
