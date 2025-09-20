package server

import (
	"WeenieHut/internal/constants"
	"WeenieHut/internal/model"
	"WeenieHut/internal/utils"
	"strconv"

	"github.com/go-playground/validator/v10"
)

type UserValidator struct {
	validator *validator.Validate
}

func NewUserValidator(v *validator.Validate) *UserValidator {
	return &UserValidator{
		validator: v,
	}
}

func (uv *UserValidator) ValidateUpdateProfileRequest(req UpdateUserProfileRequest) error {
	if (req == UpdateUserProfileRequest{}) {
		return constants.ErrInvalidRequest
	}

	if req.Email != "" {
		return constants.ErrInvalidRequest
	}

	if req.Phone != "" {
		return constants.ErrInvalidRequest
	}

	if req.Name != "" {
		return constants.ErrInvalidRequest
	}

	if req.FileURI != "" {
		return constants.ErrInvalidRequest
	}

	if req.FileThumbnailURI != "" {
		return constants.ErrInvalidRequest
	}

	if req.FileID == "" {
		return constants.ErrInvalidRequest
	}

	if _, err := strconv.ParseInt(req.FileID, 10, 64); err != nil {
		return constants.ErrInvalidFileID
	}

	return uv.validator.Struct(req)
}

func (uv *UserValidator) ValidateUpdateContactRequest(req UpdateUserContactRequest) error {
	emailProvided := req.Email != ""
	phoneProvided := req.Phone != ""

	// exactly one field must be provided
	if !emailProvided && !phoneProvided {
		return constants.ErrEmailOrPhoneMustBeProvided
	}

	if emailProvided && phoneProvided {
		return constants.ErrCannotUpdateEmailAndPhone
	}

	if emailProvided && !utils.IsEmail(req.Email) {
		return constants.ErrInvalidEmailFormat
	}

	if phoneProvided && !utils.IsValidPhoneNumber(req.Phone) {
		return constants.ErrInvalidPhoneNumberFormat
	}

	return uv.validator.Struct(req)
}

func (uv *UserValidator) ParseFileID(fileIDStr string) (int64, error) {
	if fileIDStr == "" {
		return 0, nil
	}
	return strconv.ParseInt(fileIDStr, 10, 64)
}

type UserResponseBuilder struct{}

func NewUserResponseBuilder() *UserResponseBuilder {
	return &UserResponseBuilder{}
}

func (urb *UserResponseBuilder) BuildUserResponse(user model.User) (UserResponse, error) {
	var fileIDStr string
	if user.FileID.Valid {
		fileIDStr = strconv.FormatInt(user.FileID.Int64, 10)
	} else {
		fileIDStr = ""
	}

	return UserResponse{
		Email:             utils.ToString(user.Email),
		Phone:             utils.ToString(user.Phone),
		Name:              utils.ToString(user.Name),
		FileID:            fileIDStr,
		FileURI:           utils.ToString(user.FileURI),
		FileThumbnailURI:  utils.ToString(user.FileThumbnailURI),
		BankAccountName:   utils.ToString(user.BankAccountName),
		BankAccountHolder: utils.ToString(user.BankAccountHolder),
		BankAccountNumber: utils.ToString(user.BankAccountNumber),
	}, nil
}
