package dns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_createNamedModulePath(t *testing.T) {
	tests := map[string]struct {
		tfWorkPath   string
		modName      string
		expectedPath string
	}{
		"tfWorkPath = ./": {
			tfWorkPath:   "./",
			modName:      "_0007770b-08a8-4b5f-a46b-081b772ba605-test_com",
			expectedPath: moduleFolder + "/_0007770b-08a8-4b5f-a46b-081b772ba605-test_com",
		},
		"tfWorkPath = test_path": {
			tfWorkPath:   "test_path",
			modName:      "_0007770b-08a8-4b5f-a46b-081b772ba605-test_com",
			expectedPath: "test_path/" + moduleFolder + "/_0007770b-08a8-4b5f-a46b-081b772ba605-test_com",
		},
		"tfWorkPath = ../": {
			tfWorkPath:   "../",
			modName:      "_0007770b-08a8-4b5f-a46b-081b772ba605-test_com",
			expectedPath: "../" + moduleFolder + "/_0007770b-08a8-4b5f-a46b-081b772ba605-test_com",
		},
		"blank tfWorkPath": {
			tfWorkPath:   "",
			modName:      "_0007770b-08a8-4b5f-a46b-081b772ba605-test_com",
			expectedPath: moduleFolder + "/_0007770b-08a8-4b5f-a46b-081b772ba605-test_com",
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equalf(t, test.expectedPath, createNamedModulePath(test.modName, test.tfWorkPath), "createNamedModulePath(%v, %v)", test.modName, test.tfWorkPath)
		})
	}
}
