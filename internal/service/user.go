package service

import (
	"WeenieHut/internal/constants"
	"WeenieHut/internal/model"
	"WeenieHut/internal/repository"
	"WeenieHut/internal/utils"
	"context"
	"strings"
	"sync"
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
	type validationResult struct {
		validationType string
		err            error
		fileURI        string
		thumbnailURI   string
	}

	var wg sync.WaitGroup
	resultChan := make(chan validationResult, 2)

	wg.Add(1)
	go func() {
		defer wg.Done()
		_, err := s.IsUserExist(ctx, params.UserID)
		resultChan <- validationResult{
			validationType: "user",
			err:            err,
		}
	}()

	if params.FileID != 0 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			file, err := s.repository.GetFileUpload(ctx, params.FileID)
			if err != nil {
				resultChan <- validationResult{
					validationType: "file",
					err:            constants.ErrFileNotFound,
				}
				return
			}
			resultChan <- validationResult{
				validationType: "file",
				err:            nil,
				fileURI:        file.Uri,
				thumbnailURI:   file.ThumbnailUri,
			}
		}()
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	repoParams := repository.UpdateUserProfileParams{
		UserID:            params.UserID,
		FileID:            params.FileID,
		BankAccountName:   params.BankAccountName,
		BankAccountHolder: params.BankAccountHolder,
		BankAccountNumber: params.BankAccountNumber,
	}

	for result := range resultChan {
		if result.err != nil {
			return model.User{}, result.err
		}
		if result.validationType == "file" {
			repoParams.FileThumbnailURI = result.thumbnailURI
			repoParams.FileURI = result.fileURI
		}
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

	type validationResult struct {
		fieldType string
		exists    bool
		err       error
	}

	var wg sync.WaitGroup
	resultChan := make(chan validationResult, 2)

	if params.Email != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			emailExists, err := s.repository.IsEmailExist(ctx, params.Email, userID)
			resultChan <- validationResult{
				fieldType: "email",
				exists:    emailExists,
				err:       err,
			}
		}()
	}

	if params.Phone != "" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			phoneExists, err := s.repository.IsPhoneExist(ctx, params.Phone, userID)
			resultChan <- validationResult{
				fieldType: "phone",
				exists:    phoneExists,
				err:       err,
			}
		}()
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	for result := range resultChan {
		if result.err != nil {
			return model.User{}, constants.ErrInternalServer
		}
		if result.exists {
			if result.fieldType == "email" {
				return model.User{}, constants.ErrDuplicateEmail
			}
			if result.fieldType == "phone" {
				return model.User{}, constants.ErrDuplicatePhoneNum
			}
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
