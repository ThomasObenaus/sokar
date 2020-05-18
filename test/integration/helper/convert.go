package helper

import "time"

func StrToPtr(v string) *string {
	return &v
}

func IntToPtr(v int) *int {
	return &v
}

func DurToPtr(v time.Duration) *time.Duration {
	return &v
}
