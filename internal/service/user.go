package service

import (
	"SaltySpitoon/internal/constants"
	"SaltySpitoon/internal/repository"
	"SaltySpitoon/internal/utils"
	"context"
	"strings"
)

func (s *Service) Login(ctx context.Context, email string, password string) (string, error) {
	user, err := s.repository.SelectUserByEmail(ctx, email)
	if err != nil {
		if strings.Contains(err.Error(), "sql: no rows in result set") {
			return "", constants.ErrUserNotFound
		}
		return "", err
	}
	if user.ID == 0 { // user not found
		return "", constants.ErrUserNotFound
	}

	if !utils.VerifyPassword(password, user.PasswordHash) {
		return "", constants.ErrUserWrongPassword
	}

	token, err := utils.GenerateToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *Service) Register(ctx context.Context, email string, password string) (string, error) {
	passwordHash, err := utils.HashPassword(password)
	if err != nil {
		return "", err
	}
	params := repository.CreateUserParams{
		Email:        email,
		PasswordHash: passwordHash,
	}
	userID, err := s.repository.CreateUser(ctx, params)
	if err != nil {
		if utils.IsErrDBConstraint(err) {
			return "", constants.ErrEmailAlreadyExists
		}
		return "", err
	}

	token, err := utils.GenerateToken(userID)
	if err != nil {
		return "", err
	}
	return token, nil
}
