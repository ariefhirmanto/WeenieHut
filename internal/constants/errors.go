package constants

import "errors"

var (
	ErrUserNotFound      = errors.New("user not found")
	ErrUserWrongPassword = errors.New("wrong password")
	ErrDuplicate         = errors.New("duplicate items")
	ErrDuplicatePhoneNum = errors.New("phone number already exists")
	ErrDuplicateEmail    = errors.New("email already exists")

	ErrInvalidFileType = errors.New("invalid file type")
	ErrMaximumFileSize = errors.New("size exceeds the maximum allowed file size")

	ErrInternalServer = errors.New("internal server error")
)
