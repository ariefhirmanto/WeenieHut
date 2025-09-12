package service

import (
	"WeenieHut/internal/model"
	"WeenieHut/internal/repository"
	"context"
)

type Service struct {
	repository Repository
}

// note: not ideal, might need adapter layer because return type is defined in the repository package
type Repository interface {
	// SelectUserByEmail(ctx context.Context, email string) (repository.SelectUserByEmailRow, error)
	// CreateUser(ctx context.Context, arg repository.CreateUserParams) (int64, error)
	InsertProduct(ctx context.Context, data model.Product) (model.Product, error)
	SelectProductByProductId(ctx context.Context, productIdInput int64) (repository.SelectProductByProductIdRow, error)
	SelectPaymentDetailByUserId(ctx context.Context, userId int64) (repository.SelectPaymentDetailByUserIdRow, error)
	InsertCart(ctx context.Context, arg repository.InsertCartRow) (int64, error)
	InsertCartItem(ctx context.Context, arg repository.InsertCartItemRow) (int64, error)
}

func New(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}
