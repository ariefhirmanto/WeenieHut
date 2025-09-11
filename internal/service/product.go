package service

import (
	"WeenieHut/internal/model"
	"WeenieHut/internal/server"
	"WeenieHut/internal/utils"
	"context"
	"strconv"
)

func (s *Service) PostProduct(ctx context.Context, req server.PostProductRequest) (res server.PostProductResponse, err error) {
	//
	// Set Product Value
	//
	newProduct := model.Product{
		Name:     req.Name,
		Category: utils.ToPointer(req.Category),
		Qty:      req.Qty,
		Price:    req.Price,
		SKU:      req.Sku,
		FileID:   utils.ToPointer(req.FileID),
	}

	//
	// Insert Product
	//
	insertedProduct, err := s.repository.InsertProduct(ctx, newProduct)
	if err != nil {
		return
	}

	res = server.PostProductResponse{
		ProductID:        strconv.Itoa(int(insertedProduct.ID)),
		Name:             insertedProduct.Name,
		Category:         utils.PointerValue(insertedProduct.Category, ""),
		Qty:              insertedProduct.Qty,
		Price:            insertedProduct.Price,
		Sku:              insertedProduct.SKU,
		FileID:           utils.PointerValue(insertedProduct.FileID, ""),
		FileUri:          "",
		FileThumbnailUri: "",
		CreatedAt:        insertedProduct.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:        insertedProduct.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return
}
