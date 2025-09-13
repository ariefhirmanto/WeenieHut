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
	// User Repository
	InsertUser(ctx context.Context, user model.User, passwordHash string) (model.User, error)
	SelectUserCredentialsByEmail(ctx context.Context, phone string) (model.User, error)
	SelectUserCredentialsByPhone(ctx context.Context, phone string) (model.User, error)
	GetUserProfile(ctx context.Context, userId int64) (model.User, error)
	UpdateUserProfile(ctx context.Context, param repository.UpdateUserProfileParams) (model.User, error)
	IsEmailExist(ctx context.Context, email string, excludeUserID int64) (bool, error)
	IsPhoneExist(ctx context.Context, phone string, excludeUserID int64) (bool, error)
	IsUserExist(ctx context.Context, userID int64) (bool, error)

	// File
	GetFileUpload(ctx context.Context, id int64) (model.File, error)

	// File Repository
	InsertFile(ctx context.Context, file model.File) (model.File, error)
	GetFileByFileID(ctx context.Context, fileID string) (res model.File, err error)

	// Product Repository
	InsertProduct(ctx context.Context, data model.Product) (res model.Product, err error)
	GetProducts(ctx context.Context, filter ProductFilter) (res []model.Product, err error)
	UpdateProduct(ctx context.Context, data model.Product) (res model.Product, err error)
	DeleteProductByID(ctx context.Context, id int64) (err error)
	SelectProductByProductId(ctx context.Context, productIdInput int64) (repository.SelectProductByProductIdRow, error)
	SelectPaymentDetailByUserId(ctx context.Context, userId int64) (repository.SelectPaymentDetailByUserIdRow, error)
	InsertCart(ctx context.Context, arg repository.InsertCartRow) (int64, error)
	InsertCartItem(ctx context.Context, arg repository.InsertCartItemRow) (int64, error)
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
