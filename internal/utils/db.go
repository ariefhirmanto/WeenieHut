package utils

import "strings"

func IsErrDBConstraint(err error) bool {
	return strings.Contains(err.Error(), "unique constraint")
}
