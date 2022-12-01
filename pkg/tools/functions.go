// Package tools contains various functions to help with processing of the templates
package tools

import (
	"encoding/json"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// EscapeName takes a string and makes it suitable for a tf resource instance name
// USAGE EXAMPLE: resource "akamai_appsec_waf_mode" "{{ escapename $policyName }}" {
func EscapeName(s string) (string, error) {

	// Sanitize
	s = EscapeQuotedStringLit(s)
	s = strings.TrimSpace(s)
	s = strings.ReplaceAll(s, " ", "_")

	reg, err := regexp.Compile("[^_a-zA-Z0-9]+")
	if err != nil {
		return "", err
	}
	s = reg.ReplaceAllString(s, "")

	// Terraform names can't start with a number
	i := 0
	n, _ := fmt.Sscanf(s, "%d", &i)
	if n > 0 {
		s = "ak_" + s
	}

	// Convert to lower case
	s = strings.ToLower(s)

	return s, nil
}

// ToList returns a list as a comma delimited string
// USAGE EXAMPLE: security_policy_ids = [ {{ tolist .Siem.FirewallPolicyIds }} ]
func ToList(l []string) string {
	n := []string{}
	for _, v := range l {
		v = strconv.Quote(EscapeQuotedStringLit(v))
		n = append(n, v)

	}
	return strings.Join(n, ", ")
}

// ToJSON returns a JSON representation of the given object
// USAGE EXAMPLE: "bypassNetworkLists": {{ tojson .BypassNetworkLists }},
func ToJSON(o interface{}) (string, error) {

	// Serialize and make pretty
	b, err := json.MarshalIndent(o, "", "    ")
	if err != nil {
		return "", err
	}

	return string(b), nil
}

// EscapeQuotedStringLit returns escaped terraform string literal
// https://www.terraform.io/docs/language/expressions/strings.html#escape-sequences
// This function is based on https://github.com/hashicorp/hcl/blob/c7ee8b78101c33b4dfed2641d78cf5e9651eabb8/hclwrite/generate.go#L207-L246
func EscapeQuotedStringLit(s string) string {
	if len(s) == 0 {
		return ""
	}
	buf := strings.Builder{}
	for i, r := range s {
		switch r {
		case '\n':
			buf.Write([]byte{'\\', 'n'})
		case '\r':
			buf.Write([]byte{'\\', 'r'})
		case '\t':
			buf.Write([]byte{'\\', 't'})
		case '"':
			buf.Write([]byte{'\\', '"'})
		case '\\':
			buf.Write([]byte{'\\', '\\'})
		case '$', '%':
			buf.WriteRune(r)
			remain := s[i+1:]
			if len(remain) > 0 && remain[0] == '{' {
				// Double up our template introducer symbol to escape it.
				buf.WriteRune(r)
			}
		default:
			if !unicode.IsPrint(r) {
				var fmted string
				if r < 65536 {
					fmted = fmt.Sprintf("\\u%04x", r)
				} else {
					fmted = fmt.Sprintf("\\U%08x", r)
				}
				buf.WriteString(fmted)
			} else {
				buf.WriteRune(r)
			}
		}
	}
	return buf.String()
}
