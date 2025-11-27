package tools

import (
	"regexp"
	"strings"
)

var matchFirstCap = regexp.MustCompile("([^ _])([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")
var nameRegexp = regexp.MustCompile(`[^\p{L}\p{Nl}\p{Mn}\p{Mc}\p{Nd}\p{Pc}\d\-_ ]`)

// ToSnakeCase returns name using snake case notation - SomeName -> some_name
func ToSnakeCase(str string) string {
	snake := strings.ReplaceAll(str, " ", "_")
	snake, _ = strings.CutSuffix(snake, "_")
	snake = matchFirstCap.ReplaceAllString(snake, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

// TerraformName is used to convert rule name into valid name of the exported data source
// Current implementation is not covering all the cases defined in the terraform specification
// https://github.com/hashicorp/hcl/blob/main/hclsyntax/spec.md#identifiers and http://unicode.org/reports/tr31/ ,
// but only a reasonable subset.
func TerraformName(str string) string {
	str = nameRegexp.ReplaceAllString(str, "-")
	return ToSnakeCase(str)
}
