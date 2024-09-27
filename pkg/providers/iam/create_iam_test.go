package iam

import (
	"log"
	"os"
	"testing"

	"github.com/akamai/AkamaiOPEN-edgegrid-golang/v8/pkg/iam"
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

func TestCIDRName(t *testing.T) {
	tests := map[string]struct {
		original string
		expected string
	}{
		"ipv4": {
			original: "1.1.1.1/24",
			expected: "cidr_1_1_1_1-24",
		},
		"ipv4 - no netmask": {
			original: "1.1.1.1",
			expected: "cidr_1_1_1_1",
		},
		"ipv6": {
			original: "2002::1234:abcd:ffff:c0a8:101/64",
			expected: "cidr_2002__1234_abcd_ffff_c0a8_101-64",
		},
		"ipv6 - no netmask": {
			original: "2002::1234:abcd:ffff:c0a8:101",
			expected: "cidr_2002__1234_abcd_ffff_c0a8_101",
		},
	}

	for _, test := range tests {
		actual := cidrName(test.original)
		assert.Equal(t, test.expected, actual)
	}
}
