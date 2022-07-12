// Copyright 2020. Akamai Technologies, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
		schar := string(char)
		if i == 0 {
			fc := regexp.MustCompile("^[a-zA-Z_]*$")
			if !fc.MatchString(schar) {
				outKey += "_"
			}
		}
		if re.MatchString(schar) {
			outKey += schar
		} else {
			outKey += "_"
		}
	}
	return outKey

}
