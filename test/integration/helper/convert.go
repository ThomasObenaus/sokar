package helper

import "time"

// StrToPtr returns the pointer of the given string
func StrToPtr(v string) *string {
	return &v
}

// IntToPtr returns the pointer of the given int
func IntToPtr(v int) *int {
	return &v
}

// DurToPtr returns the pointer of the given time.Duration
func DurToPtr(v time.Duration) *time.Duration {
	return &v
}
