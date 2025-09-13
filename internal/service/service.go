package service

import (
	"WeenieHut/internal/model"
	"WeenieHut/internal/repository"
	"context"
)

type Service struct {
	repository      Repository
	storage         Storage
	imageCompressor ImageCompressor
}

// note: not ideal, might need adapter layer because return type is defined in the repository package
type Repository interface {
	// User
	InsertUser(ctx context.Context, user model.User, passwordHash string) (model.User, error)
	SelectUserCredentialsByEmail(ctx context.Context, phone string) (model.User, error)
	SelectUserCredentialsByPhone(ctx context.Context, phone string) (model.User, error)
	GetUserProfile(ctx context.Context, userId int64) (model.User, error)
	UpdateUserProfile(ctx context.Context, param repository.UpdateUserProfileParams) (model.User, error)
	IsEmailExist(ctx context.Context, email string, excludeUserID int64) (bool, error)
	IsPhoneExist(ctx context.Context, phone string, excludeUserID int64) (bool, error)
	IsUserExist(ctx context.Context, userID int64) (bool, error)

	// File
	InsertFile(ctx context.Context, file model.File) (model.File, error)
	GetFileUpload(ctx context.Context, id int64) (model.File, error)
}

type Storage interface {
	UploadFile(ctx context.Context, bucket, localPath, remotePath string) (string, error)
}

type ImageCompressor interface {
	Compress(ctx context.Context, src string) (string, error)
}

func New(repository Repository, storage Storage, imageCompressor ImageCompressor) *Service {
	return &Service{
		repository:      repository,
		storage:         storage,
		imageCompressor: imageCompressor,
	}
}
