package utils

import (
	"WeenieHut/internal/constants"
	"path/filepath"
	"slices"
)

func ValidateFileExtensions(filename string, allowedExtensions []string) error {
	ext := filepath.Ext(filename)
	if slices.Contains(allowedExtensions, ext) {
		return nil
	}
	return constants.ErrInvalidFileType
}
