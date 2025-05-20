// Package iam contains code for exporting identity access manager configuration
package iam

import (
	"context"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v10/pkg/iam"
	"github.com/akamai/cli-terraform/v2/pkg/tools"
	"github.com/akamai/cli/v2/pkg/terminal"
	"github.com/urfave/cli/v2"
)

type (
	// TFData represents the iam data used in templates
	TFData struct {
		TFUsers     []*TFUser
		TFRoles     []TFRole
		TFGroups    []TFGroup
		TFAllowlist TFAllowlist
		TFClient    TFClient
		Section     string
		Subcommand  string
	}

	// TFAllowlist represents iam allowlist data used in templates
	TFAllowlist struct {
		CIDRBlocks []TFCIDRBlock
		Enabled    bool
	}

	// TFCIDRBlock represent iam cidr blocks data used in templates
	TFCIDRBlock struct {
		CIDRBlockID int64
		CIDRBlock   string
		Enabled     bool
		Comments    *string
	}

	// TFUser represents the user data used in templates
	TFUser struct {
		TFUserBasicInfo
		IsLocked          bool
		AuthGrants        string
		UserNotifications TFUserNotifications
	}

	// TFUserBasicInfo represents user basic info data used in templates
	TFUserBasicInfo struct {
		ID                       string
		FirstName                string
		LastName                 string
		Email                    string
		Country                  string
		Phone                    string
		TFAEnabled               bool
		ContactType              string
		JobTitle                 string
		TimeZone                 string
		SecondaryEmail           string
		MobilePhone              string
		Address                  string
		City                     string
		State                    string
		ZipCode                  string
		PreferredLanguage        string
		SessionTimeOut           *int
		AdditionalAuthentication string
	}

	// TFUserNotifications represents a user's notifications
	TFUserNotifications struct {
		EnableEmailNotifications              bool
		APIClientCredentialExpiryNotification bool
		NewUserNotification                   bool
		PasswordExpiry                        bool
		Proactive                             []string
		Upgrade                               []string
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

	// TFClient represents a client used in templates
	TFClient struct {
		AllowAccountSwitch      bool
		APIAccess               TFAPIAccessRequest
		AuthorizedUsers         []string
		CanCreateAutoCredential bool
		ClientDescription       string
		ClientName              string
		ClientType              iam.ClientType
		GroupAccess             TFGroupAccessRequest
		IPACL                   *TFIPACL
		NotificationEmails      []string
		PurgeOptions            *TFPurgeOptions
		Lock                    bool
		ClientID                string
		Credential              *TFCredential
	}

	// TFCredential represents a client credential used in templates.
	TFCredential struct {
		Description string
		ExpiresOn   string
		Status      string
	}

	// TFIPACL represents a client IPACL used in templates
	TFIPACL struct {
		CIDR   []string
		Enable bool
	}

	// TFPurgeOptions represents a client PurgeOptions used in templates
	TFPurgeOptions struct {
		CanPurgeByCacheTag bool
		CanPurgeByCPCode   bool
		CPCodeAccess       TFCPCodeAccess
	}

	// TFCPCodeAccess represents the CP codes used in templates which the API client can purge
	TFCPCodeAccess struct {
		AllCurrentAndNewCPCodes bool
		CPCodes                 []int64
	}

	// TFGroupAccessRequest specifies the API client's group access.
	TFGroupAccessRequest struct {
		CloneAuthorizedUserGroups bool
		Groups                    []TFClientGroupRequestItem
	}

	// TFClientGroupRequestItem represents a group the API client can access.
	TFClientGroupRequestItem struct {
		GroupID   int64
		RoleID    int64
		Subgroups []TFClientGroupRequestItem
	}

	// TFAPIAccessRequest represents the APIs the API client can access.
	TFAPIAccessRequest struct {
		AllAccessibleAPIs bool
		APIs              []TFAPIRequestItem
	}

	// TFAPIRequestItem represents single Application Programming Interface (API).
	TFAPIRequestItem struct {
		APIID       int64
		AccessLevel TFAccessLevel
	}

	// TFAccessLevel represents the access level for API.
	TFAccessLevel string
)

var (
	//go:embed templates/*
	templateFiles embed.FS

	additionalFunctions = tools.DecorateWithMultilineHandlingFunctions(map[string]any{
		"cidrName":     cidrName,
		"getLastIndex": tools.GetLastIndex,
	})

	// ErrFetchingUsers is returned when fetching users fails
	ErrFetchingUsers = errors.New("unable to fetch users under this account")
	// ErrFetchingGroups is returned when fetching groups fails
	ErrFetchingGroups = errors.New("unable to fetch groups under this account")
	// ErrFetchingRoles is returned when fetching roles fails
	ErrFetchingRoles = errors.New("unable to fetch roles under this account")
	// ErrFetchingCIDRBlocks is returned when fetching CIDR blocks fails
	ErrFetchingCIDRBlocks = errors.New("unable to fetch cidr blocks under this account")
	// ErrFetchingIPAllowlistStatus is returned when fetching IP allowlist status fails
	ErrFetchingIPAllowlistStatus = errors.New("unable to fetch ip allowlist status for this account")
	// ErrFetchingUsersWithinGroup is returned when fetching users within group fails
	ErrFetchingUsersWithinGroup = errors.New("unable to fetch users within group")
	// ErrFetchingRolesWithinGroup is returned when fetching roles within group fails
	ErrFetchingRolesWithinGroup = errors.New("unable to fetch roles within group")
	// ErrFetchingRole is returned when fetching role fails
	ErrFetchingRole = errors.New("unable to fetch role by role_id")
	// ErrFetchingUser is returned when fetching user fails
	ErrFetchingUser = errors.New("unable to fetch user by email")
	// ErrUserNotExist is returned when user does not exist
	ErrUserNotExist = errors.New("user does not exist with given email")
	// ErrMarshalUserAuthGrants is returned when marshal user auth grants failed
	ErrMarshalUserAuthGrants = errors.New("unable to marshal AuthGrants ")
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
			ID:                       user.IdentityID,
			FirstName:                user.FirstName,
			LastName:                 user.LastName,
			Email:                    user.Email,
			Country:                  user.Country,
			Phone:                    user.Phone,
			TFAEnabled:               user.TFAEnabled,
			ContactType:              user.ContactType,
			JobTitle:                 user.JobTitle,
			TimeZone:                 user.TimeZone,
			SecondaryEmail:           user.SecondaryEmail,
			MobilePhone:              user.MobilePhone,
			Address:                  user.Address,
			City:                     user.City,
			State:                    user.State,
			ZipCode:                  user.ZipCode,
			PreferredLanguage:        user.PreferredLanguage,
			SessionTimeOut:           user.SessionTimeOut,
			AdditionalAuthentication: string(user.AdditionalAuthentication),
		},
		IsLocked:          user.IsLocked,
		AuthGrants:        authGrants,
		UserNotifications: getUserNotifications(user),
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

func getUserNotifications(user *iam.User) TFUserNotifications {
	return TFUserNotifications{
		EnableEmailNotifications:              user.Notifications.EnableEmail,
		APIClientCredentialExpiryNotification: user.Notifications.Options.APIClientCredentialExpiry,
		NewUserNotification:                   user.Notifications.Options.NewUser,
		PasswordExpiry:                        user.Notifications.Options.PasswordExpiry,
		Proactive:                             user.Notifications.Options.Proactive,
		Upgrade:                               user.Notifications.Options.Upgrade,
	}
}

func getTFCIDRBlocks(ctx context.Context, client iam.IAM) ([]TFCIDRBlock, error) {

	cidrBlocks, err := client.ListCIDRBlocks(ctx, iam.ListCIDRBlocksRequest{
		Actions: true,
	})
	if err != nil {
		return nil, err
	}

	var tfCIDRBlocks []TFCIDRBlock
	for _, cidr := range cidrBlocks {
		tfCIDRBlocks = append(tfCIDRBlocks, TFCIDRBlock{
			CIDRBlockID: cidr.CIDRBlockID,
			CIDRBlock:   cidr.CIDRBlock,
			Enabled:     cidr.Enabled,
			Comments:    cidr.Comments,
		})
	}

	return tfCIDRBlocks, nil
}

func getTFClient(apiClient *iam.GetAPIClientResponse) TFClient {
	return TFClient{
		APIAccess:               getAPIAccess(apiClient),
		AuthorizedUsers:         apiClient.AuthorizedUsers,
		CanCreateAutoCredential: apiClient.CanAutoCreateCredential,
		AllowAccountSwitch:      apiClient.AllowAccountSwitch,
		ClientDescription:       apiClient.ClientDescription,
		ClientName:              apiClient.ClientName,
		ClientType:              apiClient.ClientType,
		GroupAccess:             getGroupAccess(apiClient),
		IPACL:                   getIPACL(apiClient),
		NotificationEmails:      apiClient.NotificationEmails,
		PurgeOptions:            getPurge(apiClient),
		Lock:                    apiClient.IsLocked,
		ClientID:                apiClient.ClientID,
		Credential:              getCredential(apiClient),
	}
}

func getCredential(c *iam.GetAPIClientResponse) *TFCredential {
	credentialsResponse := c.Credentials
	sort.Slice(credentialsResponse, func(i, j int) bool {
		return credentialsResponse[i].CreatedOn.Before(credentialsResponse[j].CreatedOn)
	})
	if len(credentialsResponse) > 0 {
		return &TFCredential{
			Description: credentialsResponse[0].Description,
			ExpiresOn:   credentialsResponse[0].ExpiresOn.Format(time.RFC3339Nano),
			Status:      string(credentialsResponse[0].Status),
		}
	}
	return &TFCredential{}
}

func getIPACL(apiClient *iam.GetAPIClientResponse) *TFIPACL {
	if apiClient.IPACL != nil {
		return &TFIPACL{
			CIDR:   apiClient.IPACL.CIDR,
			Enable: apiClient.IPACL.Enable,
		}
	}
	return nil
}

func getPurge(apiClient *iam.GetAPIClientResponse) *TFPurgeOptions {
	if apiClient.PurgeOptions != nil {
		return &TFPurgeOptions{
			CanPurgeByCacheTag: apiClient.PurgeOptions.CanPurgeByCacheTag,
			CanPurgeByCPCode:   apiClient.PurgeOptions.CanPurgeByCPCode,
			CPCodeAccess: TFCPCodeAccess{
				AllCurrentAndNewCPCodes: apiClient.PurgeOptions.CPCodeAccess.AllCurrentAndNewCPCodes,
				CPCodes:                 apiClient.PurgeOptions.CPCodeAccess.CPCodes,
			},
		}
	}
	return nil
}

func getAPIAccess(apiClient *iam.GetAPIClientResponse) TFAPIAccessRequest {
	apiAccess := TFAPIAccessRequest{
		AllAccessibleAPIs: apiClient.APIAccess.AllAccessibleAPIs,
	}

	for _, api := range apiClient.APIAccess.APIs {
		apiAccess.APIs = append(apiAccess.APIs, TFAPIRequestItem{
			APIID:       api.APIID,
			AccessLevel: TFAccessLevel(api.AccessLevel),
		})
	}

	return apiAccess
}

func getGroupAccess(apiClient *iam.GetAPIClientResponse) TFGroupAccessRequest {
	return TFGroupAccessRequest{
		CloneAuthorizedUserGroups: apiClient.GroupAccess.CloneAuthorizedUserGroups,
		Groups:                    buildGroupAccess(apiClient.GroupAccess.Groups),
	}
}

func buildGroupAccess(groups []iam.ClientGroup) []TFClientGroupRequestItem {
	var groupList []TFClientGroupRequestItem
	for _, subgroup := range groups {
		groupList = append(groupList, TFClientGroupRequestItem{
			GroupID:   subgroup.GroupID,
			RoleID:    subgroup.RoleID,
			Subgroups: buildGroupAccess(subgroup.Subgroups),
		})
	}
	return groupList
}

func cidrName(cidr string) string {
	cidr = strings.Replace(cidr, ".", "_", -1)
	cidr = strings.Replace(cidr, ":", "_", -1)
	cidr = strings.Replace(cidr, "/", "-", -1)

	return fmt.Sprintf("cidr_%s", cidr)
}
