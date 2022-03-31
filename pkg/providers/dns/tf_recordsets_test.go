package dns

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProcessStringNoQuotes(t *testing.T) {

	sourceString := "no quotes"
	expectedString := "no quotes"

	returnedString := processString(sourceString)

	assert.Equal(t, returnedString, expectedString)
}

func TestProcessStringEscapedQuotes(t *testing.T) {

	sourceString := "test \"four\" test"
	expectedString := `test \"four\" test`

	returnedString := processString(sourceString)

	assert.Equal(t, returnedString, expectedString)
}

func TestProcessStringEmbeddedQuotes(t *testing.T) {

	sourceString := `first string" "secondString`
	expectedString := `first string\" \"secondString`

	returnedString := processString(sourceString)

	assert.Equal(t, returnedString, expectedString)
}
