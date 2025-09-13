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
	ErrFileNotFound    = errors.New("file not found")
	ErrInternalServer  = errors.New("internal server error")

	ErrFileIDNotValid                 = errors.New("fileId is not valid / exists")
	ErrDuplicateSKU                   = errors.New("duplicate sku")
	ErrProductNotFound                = errors.New("productId is not found")
	ErrNotEqualAvailableSellersInCart = errors.New("not equal to the available sellers in the cart")
	ErrPurchaseNotFound               = errors.New("purchase not found")
)
