package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNormalizeWholeFile(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected string
	}{
		"lines out of order are sorted alphabetically": {
			input:    "line_c\nline_b\nline_a",
			expected: "line_a\nline_b\nline_c",
		},
		"already sorted input is unchanged": {
			input:    "aaa\nbbb\nccc",
			expected: "aaa\nbbb\nccc",
		},
		"single line is unchanged": {
			input:    "single line",
			expected: "single line",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.expected, NormalizeWholeFile(test.input))
		})
	}
}

func TestNormalizeFieldsInBlock(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected string
	}{
		"fields within a single block are sorted": {
			input: `resource "some_resource" "example" {
    field_c = "value_c"
    field_a = "value_a"
    field_b = "value_b"
}`,
			expected: `resource "some_resource" "example" {
    field_a = "value_a"
    field_b = "value_b"
    field_c = "value_c"
}`,
		},
		"each block is sorted independently": {
			input: `resource "some_resource" "first" {
    field_b = "value_b"
    field_a = "value_a"
}
resource "some_resource" "second" {
    field_d = "value_d"
    field_c = "value_c"
}`,
			expected: `resource "some_resource" "first" {
    field_a = "value_a"
    field_b = "value_b"
}
resource "some_resource" "second" {
    field_c = "value_c"
    field_d = "value_d"
}`,
		},
		"non-block lines are preserved in their original position": {
			input: `locals {
    some_local = "value"
}

resource "some_resource" "example" {
    field_b = "value_b"
    field_a = "value_a"
}`,
			expected: `locals {
    some_local = "value"
}

resource "some_resource" "example" {
    field_a = "value_a"
    field_b = "value_b"
}`,
		},
		"content with no matching blocks is unchanged": {
			input: `locals {
    some_local = "value"
}`,
			expected: `locals {
    some_local = "value"
}`,
		},
		"empty block is preserved": {
			input: `resource "some_resource" "example" {
}`,
			expected: `resource "some_resource" "example" {
}`,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.expected, NormalizeFieldsInBlock(test.input, `resource "some_resource"`))
		})
	}
}

func TestNormalizeBlocksInFile(t *testing.T) {
	tests := map[string]struct {
		input    string
		expected string
	}{
		"blocks out of order are sorted alphabetically": {
			input: `module "module_c" {
    field = "value"
}
module "module_a" {
    field = "value"
}`,
			expected: `
module "module_a" {
    field = "value"
}
module "module_c" {
    field = "value"
}`,
		},
		"non-block lines are placed before sorted blocks": {
			input: `locals {
    some_local = "value"
}

module "module_b" {
    field = "value"
}
module "module_a" {
    field = "value"
}`,
			expected: `locals {
    some_local = "value"
}

module "module_a" {
    field = "value"
}
module "module_b" {
    field = "value"
}`,
		},
		"already sorted blocks are unchanged": {
			input: `module "module_a" {
    field = "value"
}
module "module_b" {
    field = "value"
}`,
			expected: `
module "module_a" {
    field = "value"
}
module "module_b" {
    field = "value"
}`,
		},
		"content with no matching blocks is unchanged": {
			input: `locals {
    some_local = "value"
}`,
			expected: `locals {
    some_local = "value"
}`,
		},
		"single block is preserved": {
			input: `module "module_a" {
    field = "value"
}`,
			expected: `
module "module_a" {
    field = "value"
}`,
		},
		"block with multiple fields is preserved as a unit when sorted": {
			input: `module "module_b" {
    field_1 = "value_1"
    field_2 = "value_2"
    field_3 = "value_3"
}
module "module_a" {
    field_1 = "value_1"
    field_2 = "value_2"
    field_3 = "value_3"
}`,
			expected: `
module "module_a" {
    field_1 = "value_1"
    field_2 = "value_2"
    field_3 = "value_3"
}
module "module_b" {
    field_1 = "value_1"
    field_2 = "value_2"
    field_3 = "value_3"
}`,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.expected, NormalizeBlocksInFile(test.input, `module "`))
		})
	}
}
