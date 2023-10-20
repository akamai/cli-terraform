package tools

import (
	"strings"
	"text/template"
)

// DecorateWithMultilineHandlingFunctions adds necessary functions to escape multiline fields in exported terraform files.
func DecorateWithMultilineHandlingFunctions(additionalFuncs map[string]any) template.FuncMap {
	multiLineFuncs := map[string]any{
		"Escape":            Escape,
		"IsMultiline":       IsMultiline,
		"NoNewlineAtTheEnd": NoNewlineAtTheEnd,
		"RemoveLastNewline": RemoveLastNewline,
		"GetEOT":            GetEOT,
	}
	for k, v := range additionalFuncs {
		multiLineFuncs[k] = v
	}
	return multiLineFuncs
}

// IsMultiline returns true if the input string contains at least one new line character
func IsMultiline(str string) bool {
	return strings.LastIndex(str, "\n") >= 0
}

// NoNewlineAtTheEnd returns true if there is no new line character at the end of the string
func NoNewlineAtTheEnd(str string) bool {
	if str == "" {
		return true
	}
	return str[len(str)-1:] != "\n"
}

// RemoveLastNewline removes the new line character if this is the last character in the string
func RemoveLastNewline(str string) string {
	return strings.TrimSuffix(str, "\n")
}

// GetEOT generates unique delimiter word for heredoc, by default it is EOT
func GetEOT(str string) string {
	eot := "EOT"
	for strings.LastIndex(str, eot) >= 0 {
		eot += "A"
	}
	return eot
}
