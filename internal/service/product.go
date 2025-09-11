package service

import (
	"WeenieHut/internal/model"
	"WeenieHut/internal/server"
	"WeenieHut/internal/utils"
	"context"
	"errors"
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

type ProductFilter struct {
	ProductID *int64
	Sku       string
	Category  string
	SortBy    string
	Limit     int
	Offset    int
}

func (s *Service) GetProducts(ctx context.Context, req server.GetProductsRequest) (res []server.GetProductResponse, err error) {
	var productIDPtr *int64
	if req.ProductID != "" {
		if productIDInt, err := strconv.Atoi(req.ProductID); err == nil {
			productIDPtr = utils.ToPointer(int64(productIDInt))
		}
	}

	limitInt, _ := strconv.Atoi(req.Limit)
	offsetInt, _ := strconv.Atoi(req.Offset)
	filter := ProductFilter{
		ProductID: productIDPtr,
		Sku:       req.Sku,
		Category:  req.Category,
		SortBy:    req.SortBy,
		Limit:     limitInt,
		Offset:    offsetInt,
	}

	if filter.Limit <= 0 {
		filter.Limit = 5
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}

	switch filter.SortBy {
	case "newest", "oldest", "cheapest", "expensive":

	case "":
		filter.SortBy = "newest"
	default:
		filter.SortBy = "newest"
	}

	products, err := s.repository.GetProducts(ctx, filter)
	if err != nil {
		return nil, err
	}
	if len(products) == 0 {
		return nil, errors.New("no products found")
	}

	res = make([]server.GetProductResponse, 0, len(products))
	for _, p := range products {
		res = append(res, server.GetProductResponse{
			ProductID:        strconv.FormatInt(p.ID, 10),
			Name:             p.Name,
			Category:         utils.PointerValue(p.Category, ""),
			Qty:              p.Qty,
			Price:            p.Price,
			Sku:              p.SKU,
			FileID:           utils.PointerValue(p.FileID, ""),
			FileUri:          utils.PointerValue(p.FileURI, ""),
			FileThumbnailUri: utils.PointerValue(p.FileThumbnailURI, ""),
			CreatedAt:        p.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt:        p.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	return
}
