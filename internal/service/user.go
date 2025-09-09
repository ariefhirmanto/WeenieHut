package service

import (
	"context"
)

func (s *Service) Login(ctx context.Context, email string, password string) (string, error) {
	// TODO: implement login
	// user, err := s.repository.SelectUserByEmail(ctx, email)
	// if err != nil {
	// 	if strings.Contains(err.Error(), "sql: no rows in result set") {
	// 		return "", constants.ErrUserNotFound
	// 	}
	// 	return "", err
	// }
	// if user.ID == 0 { // user not found
	// 	return "", constants.ErrUserNotFound
	// }

	// if !utils.VerifyPassword(password, user.PasswordHash) {
	// 	return "", constants.ErrUserWrongPassword
	// }

	// token, err := utils.GenerateToken(user.ID)
	// if err != nil {
	// 	return "", err
	// }

	// return token, nil
	return "", nil
}

func (s *Service) Register(ctx context.Context, email string, password string) (string, error) {
	// TODO: implement register
	// passwordHash, err := utils.HashPassword(password)
	// if err != nil {
	// 	return "", err
	// }
	// params := repository.CreateUserParams{
	// 	Email:        email,
	// 	PasswordHash: passwordHash,
	// }
	// userID, err := s.repository.CreateUser(ctx, params)
	// if err != nil {
	// 	if utils.IsErrDBConstraint(err) {
	// 		return "", constants.ErrEmailAlreadyExists
	// 	}
	// 	return "", err
	// }

	// token, err := utils.GenerateToken(userID)
	// if err != nil {
	// 	return "", err
	// }
	// return token, nil
	return "", nil
}
