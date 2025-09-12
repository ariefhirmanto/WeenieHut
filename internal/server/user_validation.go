package server

import (
	"WeenieHut/internal/constants"
	"WeenieHut/internal/model"
	"WeenieHut/internal/utils"
	"errors"
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
	return uv.validator.Struct(req)
}

func (uv *UserValidator) ValidateUpdateContactRequest(req UpdateUserContactRequest) error {
	emailProvided := req.Email != ""
	phoneProvided := req.Phone != ""

	// exactly one field must be provided
	if !emailProvided && !phoneProvided {
		return errors.New("either email or phone must be provided")
	}

	if emailProvided && phoneProvided {
		return errors.New("cannot update both email and phone simultaneously")
	}

	if emailProvided && !utils.IsEmail(req.Email) {
		return errors.New("invalid email format")
	}

	if phoneProvided && !utils.IsValidPhoneNumber(req.Phone) {
		return errors.New("invalid phone number format")
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
	fileIdResp := int(user.FileID)
	if int64(fileIdResp) != user.FileID {
		return UserResponse{}, constants.ErrInternalServer
	}

	return UserResponse{
		Email:             utils.ToString(user.Email),
		Phone:             utils.ToString(user.Phone),
		Name:              user.Name,
		FileID:            strconv.Itoa(fileIdResp),
		FileURI:           user.FileURI,
		FileThumbnailURI:  user.FileThumbnailURI,
		BankAccountName:   user.BankAccountName,
		BankAccountHolder: user.BankAccountHolder,
		BankAccountNumber: user.BankAccountNumber,
	}, nil
}
