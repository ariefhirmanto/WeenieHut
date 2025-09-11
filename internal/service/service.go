package service

import (
	"WeenieHut/internal/model"
	"context"
)

type Service struct {
	repository Repository
}

// note: not ideal, might need adapter layer because return type is defined in the repository package
type Repository interface {
	// SelectUserByEmail(ctx context.Context, email string) (repository.SelectUserByEmailRow, error)
	// CreateUser(ctx context.Context, arg repository.CreateUserParams) (int64, error)

	// Product Repository
	InsertProduct(ctx context.Context, data model.Product) (res model.Product, err error)
	GetProducts(ctx context.Context, filter ProductFilter) (res []model.Product, err error)
	UpdateProduct(ctx context.Context, data model.Product) (res model.Product, err error)
	DeleteProductByID(ctx context.Context, id int64) (err error)
}

func New(repository Repository) *Service {
	return &Service{
		repository: repository,
	}
}
