package iam

import (
	"embed"
)

type (
	// TFUserData represents the user data used in templates
	TFUserData struct {
		TFUserBasicInfo
		IsLocked   bool
		AuthGrants string
		Section    string
		Roles      []TFRole
		Groups     []TFGroup
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
)
