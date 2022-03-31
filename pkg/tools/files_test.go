package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckFiles(t *testing.T) {
	tests := map[string]struct {
		given     []string
		withError bool
	}{
		"files do not exist": {
			given:     []string{"testdata/f3.txt", "testdata/f4.txt"},
			withError: false,
		},
		"some files exist": {
			given:     []string{"testdata/f1.txt", "testdata/f3.txt"},
			withError: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			err := CheckFiles(test.given...)
			if test.withError {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
		})
	}
}
