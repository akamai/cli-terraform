package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
