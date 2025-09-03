package constants

import "errors"

var (
	ErrUserNotFound       = errors.New("user not found")
	ErrUserWrongPassword  = errors.New("wrong password")
	ErrEmailAlreadyExists = errors.New("email already exists")
)
