package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsMultiline(t *testing.T) {
	tests := map[string]struct {
		given    string
		expected bool
	}{
		"has new lines": {
			given:    "this\nis test\n",
			expected: true,
		},
		"empty": {
			given:    "",
			expected: false,
		},
		"no new lines": {
			given:    "no new lines",
			expected: false,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.expected, IsMultiline(test.given))
		})
	}
}

func TestNoNewlineAtTheEnd(t *testing.T) {
	tests := map[string]struct {
		given    string
		expected bool
	}{
		"has new line at the end": {
			given:    "this\nis test\n",
			expected: false,
		},
		"has new line in the middle": {
			given:    "this\nis test",
			expected: true,
		},
		"empty": {
			given:    "",
			expected: true,
		},
		"no new lines": {
			given:    "no new lines",
			expected: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.expected, NoNewlineAtTheEnd(test.given))
		})
	}
}

func TestRemoveLastNewline(t *testing.T) {
	tests := map[string]struct {
		given    string
		expected string
	}{
		"has new line at the end": {
			given:    "this\nis test\n",
			expected: "this\nis test",
		},
		"has new line in the middle": {
			given:    "this\nis test",
			expected: "this\nis test",
		},
		"empty": {
			given:    "",
			expected: "",
		},
		"no new lines": {
			given:    "no new lines",
			expected: "no new lines",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.expected, RemoveLastNewline(test.given))
		})
	}
}

func TestGetEOT(t *testing.T) {
	tests := map[string]struct {
		given    string
		expected string
	}{
		"has new line": {
			given:    "this\nis test\n",
			expected: "EOT",
		},
		"has EOT inside": {
			given:    "this\nEOT",
			expected: "EOTA",
		},
		"empty": {
			given:    "",
			expected: "EOT",
		},
		"has two EOTs": {
			given:    "some\nEOT\nEOTA\ntext",
			expected: "EOTAA",
		},
		"has EOT": {
			given:    "comment\nnewline\nand\nEOT\ninside\n",
			expected: "EOTA",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, test.expected, GetEOT(test.given))
		})
	}
}
