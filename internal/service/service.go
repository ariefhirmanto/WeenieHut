package service

import "context"

type Service struct {
	repository Repository
	storage    Storage
}

// note: not ideal, might need adapter layer because return type is defined in the repository package
type Repository interface {
	// SelectUserByEmail(ctx context.Context, email string) (repository.SelectUserByEmailRow, error)
	// CreateUser(ctx context.Context, arg repository.CreateUserParams) (int64, error)
}

type Storage interface {
	UploadFile(ctx context.Context, bucket, localPath, remotePath string) (string, error)
}

type ImageProcessor interface {
	Compress(ctx context.Context, src, dest string) error
}

func New(repository Repository, storage Storage) *Service {
	return &Service{
		repository: repository,
		storage:    storage,
	}
}
