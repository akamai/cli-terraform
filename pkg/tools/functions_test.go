package tools

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestToJSON(t *testing.T) {
	test := `{"context":"DOSATCK","id":2670509,"name":"DoS Attackers (High Threat)","sharedIpHandling":"NON_SHARED","threshold":9}`
	var i interface{}
	err := json.Unmarshal([]byte(test), &i)
	assert.NoError(t, err)

	j, err := ToJSON(i)
	assert.NoError(t, err)

	assert.Equal(t, "{\n    \"context\": \"DOSATCK\",\n    \"id\": 2670509,\n    \"name\": \"DoS Attackers (High Threat)\",\n    \"sharedIpHandling\": \"NON_SHARED\",\n    \"threshold\": 9\n}", j)
}

func TestToList(t *testing.T) {
	tests := []string{"this", "is", "a", "list", "of", "strings"}
	assert.Equal(t, "\"this\", \"is\", \"a\", \"list\", \"of\", \"strings\"", ToList(tests))
}

func TestEscapeTFName(t *testing.T) {
	tests := []string{"This is a test", "This Is A Test", "123 This is a test", "!This is a test!", "This is a test!", "This_is-a$test!"}
	expected := []string{"this_is_a_test", "this_is_a_test", "ak_123_this_is_a_test", "this_is_a_test", "this_is_a_test", "this_isatest"}
	for i, s := range tests {
		e, err := EscapeName(s)
		require.NoError(t, err)
		assert.Equal(t, expected[i], e)
	}
}
func TestEscapeQuotedStringLit(t *testing.T) {
	tests := map[string]struct {
		data   string
		expect string
	}{
		"string": {
			data:   "foo",
			expect: "foo",
		},
		"string with quotes": {
			data:   `"foo"`,
			expect: `\"foo\"`,
		},
		"new line character": {
			data:   "hello\nworld\n",
			expect: `hello\nworld\n`,
		},
		"new line and carriage return": {
			data:   "hello\r\nworld\r\n",
			expect: `hello\r\nworld\r\n`,
		},
		"backslash": {
			data:   `what\what`,
			expect: `what\\what`,
		},
		"unicode char": {
			data:   "𝄞",
			expect: "𝄞",
		},
		"non backslash escape sequence": {
			data:   "%{",
			expect: "%%{",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			got := EscapeQuotedStringLit(test.data)
			assert.Equal(t, got, test.expect)
		})
	}
}
