package service

import (
	"WeenieHut/internal/model"
	"context"
)

type Service struct {
	repository      Repository
	storage         Storage
	imageCompressor ImageCompressor
}

// note: not ideal, might need adapter layer because return type is defined in the repository package
type Repository interface {
	// SelectUserByEmail(ctx context.Context, email string) (repository.SelectUserByEmailRow, error)
	// CreateUser(ctx context.Context, arg repository.CreateUserParams) (int64, error)
	InsertFile(ctx context.Context, file model.File) (model.File, error)
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
