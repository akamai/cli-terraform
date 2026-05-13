package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSanitizeResourceName(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected string
	}{
		"plain name unchanged":          {input: "mygroup", expected: "mygroup"},
		"spaces replaced":               {input: "my group", expected: "my_group"},
		"dots replaced":                 {input: "my.group", expected: "my_group"},
		"leading digit gets underscore": {input: "42group", expected: "_42group"},
		"leading dot gets underscore":   {input: ".group", expected: "_group"},
		"hyphens preserved":             {input: "my-group", expected: "my-group"},
		"mixed":                         {input: "my group.name", expected: "my_group_name"},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expected, SanitizeResourceName(tc.input))
		})
	}
}

func TestNotations(t *testing.T) {
	tests := map[string]struct {
		given    string
		expected string
	}{
		"basic": {
			given:    "test name ",
			expected: "test_name",
		},
		"multiple words with numbers": {
			given:    "word test 1 hello2",
			expected: "word_test_1_hello2",
		},
		"capitalized words": {
			given:    "TestNameHello1 world",
			expected: "test_name_hello1_world",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			actual := TerraformName(test.given)
			assert.Equal(t, test.expected, actual)
		})
	}
}
