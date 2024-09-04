// Package iam contains code for exporting identity access manager configuration
package iam

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/iam"
	"github.com/akamai/cli-terraform/pkg/tools"
	"github.com/akamai/cli/pkg/terminal"
	"github.com/urfave/cli/v2"
)

type (
	// TFData represents the iam data used in templates
	TFData struct {
		TFUsers    []*TFUser
		TFRoles    []TFRole
		TFGroups   []TFGroup
		Section    string
		Subcommand string
	}

	// TFUser represents the user data used in templates
	TFUser struct {
		TFUserBasicInfo
		IsLocked   bool
		AuthGrants string
	}

	// TFUserBasicInfo represents user basic info data used in templates
	TFUserBasicInfo struct {
		ID                string
		FirstName         string
		LastName          string
		Email             string
		Country           string
		Phone             string
		TFAEnabled        bool
		ContactType       string
		JobTitle          string
		TimeZone          string
		SecondaryEmail    string
		MobilePhone       string
		Address           string
		City              string
		State             string
		ZipCode           string
		PreferredLanguage string
		SessionTimeOut    *int
	}

	// TFRole represents a role used in templates
	TFRole struct {
		RoleID          int64
		RoleName        string
		RoleDescription string
		GrantedRoles    []int
	}

	// TFGroup represents a group used in templates
	TFGroup struct {
		GroupID       int
		ParentGroupID int
		GroupName     string
	}
)

var (
	//go:embed templates/*
	templateFiles embed.FS

	additionalFunctions = tools.DecorateWithMultilineHandlingFunctions(map[string]any{})

	// ErrFetchingUsers is returned when fetching users fails
	ErrFetchingUsers = errors.New("unable to fetch users under this account")
)

// CmdCreateIAM is an entrypoint to create-iam command. This is only for action validation purpose
func CmdCreateIAM(_ *cli.Context) error {
	return nil
}

func getTFUsers(ctx context.Context, client iam.IAM, users []iam.UserListItem, term terminal.Terminal) ([]*TFUser, error) {
	res := make([]*TFUser, 0)
	for _, v := range users {
		user, err := client.GetUser(ctx, iam.GetUserRequest{
			IdentityID:    v.IdentityID,
			Actions:       true,
			AuthGrants:    true,
			Notifications: true,
		})
		if err != nil {
			_, err := term.Writeln(fmt.Sprintf("[WARN] Unable to fetch user of ID '%s' - skipping:\n%s", v.IdentityID, err))
			if err != nil {
				return nil, err
			}
			continue
		}
		tfUser, err := getTFUser(user)
		if err != nil {
			return nil, err
		}
		res = append(res, tfUser)
	}

	return res, nil
}

func getTFUser(user *iam.User) (*TFUser, error) {
	authGrants, err := getUserAuthGrants(user)
	if err != nil {
		return nil, fmt.Errorf("%w: %s for user with email %s", ErrMarshalUserAuthGrants, err, user.Email)
	}

	return &TFUser{
		TFUserBasicInfo: TFUserBasicInfo{
			ID:                user.IdentityID,
			FirstName:         user.FirstName,
			LastName:          user.LastName,
			Email:             user.Email,
			Country:           user.Country,
			Phone:             user.Phone,
			TFAEnabled:        user.TFAEnabled,
			ContactType:       user.ContactType,
			JobTitle:          user.JobTitle,
			TimeZone:          user.TimeZone,
			SecondaryEmail:    user.SecondaryEmail,
			MobilePhone:       user.MobilePhone,
			Address:           user.Address,
			City:              user.City,
			State:             user.State,
			ZipCode:           user.ZipCode,
			PreferredLanguage: user.PreferredLanguage,
			SessionTimeOut:    user.SessionTimeOut,
		},
		IsLocked:   user.IsLocked,
		AuthGrants: authGrants,
	}, nil
}

func getTFGroup(group *iam.Group) TFGroup {
	return TFGroup{
		GroupID:       int(group.GroupID),
		ParentGroupID: int(group.ParentGroupID),
		GroupName:     group.GroupName,
	}
}

func getUserAuthGrants(user *iam.User) (string, error) {
	var authGrantsJSON []byte
	var err error
	if len(user.AuthGrants) > 0 {
		authGrantsJSON, err = json.Marshal(user.AuthGrants)
		if err != nil {
			return "", err
		}

		var authGrants []iam.AuthGrantRequest
		err = json.Unmarshal(authGrantsJSON, &authGrants)
		if err != nil {
			return "", err
		}

		authGrantsJSON, err = json.Marshal(authGrants)
		if err != nil {
			return "", err
		}

	}
	return string(authGrantsJSON), nil
}

func getTFRoles(ctx context.Context, client iam.IAM, roles []iam.Role) ([]TFRole, error) {
	tfRoles := make([]TFRole, 0)
	for _, r := range roles {
		roleID := r.RoleID
		role, err := client.GetRole(ctx, iam.GetRoleRequest{
			ID:           roleID,
			GrantedRoles: true,
		})
		if err != nil {
			return nil, fmt.Errorf("could not get role of ID '%v': %w", roleID, err)
		}
		tfRoles = append(tfRoles, TFRole{
			RoleID:          role.RoleID,
			RoleName:        role.RoleName,
			RoleDescription: role.RoleDescription,
			GrantedRoles:    getGrantedRolesID(role.GrantedRoles),
		})
	}
	return tfRoles, nil
}

func getGrantedRolesID(grantedRoles []iam.RoleGrantedRole) []int {
	rolesIDs := make([]int, 0, len(grantedRoles))
	for _, v := range grantedRoles {
		rolesIDs = append(rolesIDs, int(v.RoleID))
	}
	return rolesIDs
}
