package iam

import (
	"log"
	"os"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v7/pkg/iam"
	"github.com/stretchr/testify/assert"
)

func TestMain(m *testing.M) {
	if err := os.MkdirAll("./testdata/res", 0755); err != nil {
		log.Fatal(err)
	}
	exitCode := m.Run()
	if err := os.RemoveAll("./testdata/res"); err != nil {
		log.Fatal(err)
	}
	os.Exit(exitCode)
}

func TestGetGrantedRolesID(t *testing.T) {
	tests := map[string]struct {
		grantedRoles []iam.RoleGrantedRole
		expectedIDs  []int
	}{
		"granted roles": {
			grantedRoles: []iam.RoleGrantedRole{
				{
					RoleID: 123,
				},
				{
					RoleID: 321,
				},
				{
					RoleID: 456,
				},
			},
			expectedIDs: []int{123, 321, 456},
		},
		"empty granted roles": {
			grantedRoles: []iam.RoleGrantedRole{},
			expectedIDs:  []int{},
		},
		"nil granted roles": {
			grantedRoles: nil,
			expectedIDs:  []int{},
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			grantedRolesIDs := getGrantedRolesID(test.grantedRoles)
			assert.Equal(t, test.expectedIDs, grantedRolesIDs)
		})
	}
}
