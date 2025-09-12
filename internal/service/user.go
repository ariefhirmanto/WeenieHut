package service

import (
	"WeenieHut/internal/constants"
	"WeenieHut/internal/model"
	"WeenieHut/internal/utils"
	"context"
	"strings"
)

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
