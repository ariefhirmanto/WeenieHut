package service

import (
	"WeenieHut/internal/constants"
	"WeenieHut/internal/model"
	"WeenieHut/internal/repository"
	"WeenieHut/internal/utils"
	"context"
	"strings"
)

type UpdateUserParams struct {
	UserID            int64
	FileID            int64
	BankAccountName   string
	BankAccountHolder string
	BankAccountNumber string
	Phone             string
	Email             string
}

func (s *Service) EmailLogin(ctx context.Context, email string, password string) (token string, phone string, err error) {
	user, err := s.repository.SelectUserCredentialsByEmail(ctx, email)
	if err != nil {
		if strings.Contains(err.Error(), "sql: no rows in result set") {
			return "", "", constants.ErrUserNotFound
		}
		return "", "", err
	}
	if user.ID == 0 { // user not found
		return "", "", constants.ErrUserNotFound
	}

	if !utils.VerifyPassword(password, user.PasswordHash) {
		return "", "", constants.ErrUserWrongPassword
	}

	token, err = utils.GenerateToken(user.ID)
	if err != nil {
		return "", "", err
	}

	return token, user.Phone.String, nil
}

func (s *Service) PhoneLogin(ctx context.Context, phone string, password string) (token string, email string, err error) {
	user, err := s.repository.SelectUserCredentialsByPhone(ctx, phone)
	if err != nil {
		if strings.Contains(err.Error(), "sql: no rows in result set") {
			return "", "", constants.ErrUserNotFound
		}
		return "", "", err
	}
	if user.ID == 0 { // user not found
		return "", "", constants.ErrUserNotFound
	}

	if !utils.VerifyPassword(password, user.PasswordHash) {
		return "", "", constants.ErrUserWrongPassword
	}

	token, err = utils.GenerateToken(user.ID)
	if err != nil {
		return "", "", err
	}

	return token, user.Email.String, nil
}

func (s *Service) Register(ctx context.Context, userReq model.User, password string) (string, error) {
	passwordHash, err := utils.HashPassword(password)
	if err != nil {
		return "", err
	}

	user, err := s.repository.InsertUser(ctx, userReq, passwordHash)
	if err != nil {
		if utils.IsErrDBConstraint(err) {
			return "", constants.ErrDuplicate
		}
		return "", err
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (s *Service) GetUserProfile(ctx context.Context, userId int64) (model.User, error) {
	if userId == 0 {
		return model.User{}, constants.ErrUserNotFound
	}

	user, err := s.repository.GetUserProfile(ctx, userId)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (s *Service) UpdateUserProfile(ctx context.Context, params UpdateUserParams) (model.User, error) {
	_, err := s.IsUserExist(ctx, params.UserID)
	if err != nil {
		return model.User{}, err
	}

	repoParams := repository.UpdateUserProfileParams{
		UserID:            params.UserID,
		FileID:            params.FileID,
		BankAccountName:   params.BankAccountName,
		BankAccountHolder: params.BankAccountHolder,
		BankAccountNumber: params.BankAccountNumber,
	}

	if params.FileID != 0 {
		file, err := s.repository.GetFileUpload(ctx, params.FileID)
		if err != nil {
			return model.User{}, constants.ErrFileNotFound
		}

		repoParams.FileThumbnailURI = file.ThumbnailUri
		repoParams.FileURI = file.Uri
	}

	user, err := s.repository.UpdateUserProfile(ctx, repoParams)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (s *Service) UpdateUserContact(ctx context.Context, params UpdateUserParams) (model.User, error) {
	userID := params.UserID
	if userID == 0 {
		return model.User{}, constants.ErrUserNotFound
	}

	_, err := s.IsUserExist(ctx, userID)
	if err != nil {
		return model.User{}, err
	}

	if params.Email != "" {
		emailExists, err := s.repository.IsEmailExist(ctx, params.Email, userID)
		if err != nil {
			return model.User{}, constants.ErrInternalServer
		}
		if emailExists {
			return model.User{}, constants.ErrDuplicateEmail
		}
	}

	if params.Phone != "" {
		phoneExists, err := s.repository.IsPhoneExist(ctx, params.Phone, userID)
		if err != nil {
			return model.User{}, constants.ErrInternalServer
		}
		if phoneExists {
			return model.User{}, constants.ErrDuplicatePhoneNum
		}
	}

	repoParams := repository.UpdateUserProfileParams{
		UserID: userID,
		Email:  params.Email,
		Phone:  params.Phone,
	}

	user, err := s.repository.UpdateUserProfile(ctx, repoParams)
	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (s *Service) IsUserExist(ctx context.Context, userID int64) (bool, error) {
	isUserExist, err := s.repository.IsUserExist(ctx, userID)
	if err != nil {
		return false, constants.ErrInternalServer
	}
	if !isUserExist {
		return false, constants.ErrUserNotFound
	}

	return true, nil
}
