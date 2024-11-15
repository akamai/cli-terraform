package dns

import (
	"regexp"
)

// Utility function to normalize resource names. A name must start with a letter or
// underscore and may contain only letters, digits, underscores, and dashes.
func normalizeResourceName(inKey string) string {
	outKey := ""
	re := regexp.MustCompile("^[a-zA-Z0-9_-]*$")
	for i, char := range inKey {
		sChar := string(char)
		if i == 0 {
			fc := regexp.MustCompile("^[a-zA-Z_]*$")
			if !fc.MatchString(sChar) {
				outKey += "_"
			}
		}
		if re.MatchString(sChar) {
			outKey += sChar
		} else {
			outKey += "_"
		}
	}

	return outKey
}
