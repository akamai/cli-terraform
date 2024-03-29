package tools

// Float64Ptr returns the address of the float64
func Float64Ptr(f float64) *float64 {
	return &f
}

// StringPtr returns the address of the string
func StringPtr(s string) *string {
	return &s
}

// IntPtr returns the address of the int
func IntPtr(i int) *int {
	return &i
}

// Int64Ptr returns the address of the int64
func Int64Ptr(i int64) *int64 {
	return &i
}
