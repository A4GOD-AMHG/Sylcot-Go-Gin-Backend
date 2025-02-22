package utils

import "strings"

func IsDuplicateError(err error) bool {
	return strings.Index(strings.ToLower(err.Error()), "duplicate") >= 0
}
