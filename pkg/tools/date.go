package tools

import "time"

// ParseRFC3339 returns RFC3339 string representation for provided time. Intended for tests only.
func ParseRFC3339(s string) time.Time {
	parsedTime, err := time.Parse(time.RFC3339, s)
	if err != nil {
		panic(err)
	}
	return parsedTime
}
