package service

import (
	"WeenieHut/internal/constants"
	"WeenieHut/internal/model"
	"WeenieHut/internal/utils"
	"context"
	"database/sql"
	"errors"
	"strconv"
)

func (s *Service) PostProduct(ctx context.Context, req model.PostProductRequest) (res model.PostProductResponse, err error) {
	//
	// Get File Uri & Thumbnail Uri
	//
	var fileUri, fileThumbnailUri string
	if req.FileID != "" {
		file, err := s.repository.GetFileByFileID(ctx, req.FileID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return res, constants.ErrFileIDNotValid
			}
			return res, err
		}
		fileUri = file.Uri
		fileThumbnailUri = file.ThumbnailUri
	}

	//
	// Set Product Value
	//
	newProduct := model.Product{
		Name:             req.Name,
		Category:         utils.ToPointer(req.Category),
		Qty:              req.Qty,
		Price:            req.Price,
		SKU:              req.Sku,
		FileID:           utils.ToPointer(req.FileID),
		FileURI:          utils.ToPointer(fileUri),
		FileThumbnailURI: utils.ToPointer(fileThumbnailUri),
	}

	//
	// Insert Product
	//
	insertedProduct, err := s.repository.InsertProduct(ctx, newProduct)
	if err != nil {
		if errors.Is(err, constants.ErrDuplicateSKU) {
			return res, constants.ErrDuplicateSKU
		}
		return res, err
	}

	res = model.PostProductResponse{
		ProductID:        strconv.Itoa(int(insertedProduct.ID)),
		Name:             insertedProduct.Name,
		Category:         utils.PointerValue(insertedProduct.Category, ""),
		Qty:              insertedProduct.Qty,
		Price:            insertedProduct.Price,
		Sku:              insertedProduct.SKU,
		FileID:           utils.PointerValue(insertedProduct.FileID, ""),
		FileUri:          utils.PointerValue(insertedProduct.FileURI, ""),
		FileThumbnailUri: utils.PointerValue(insertedProduct.FileThumbnailURI, ""),
		CreatedAt:        insertedProduct.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:        insertedProduct.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return res, nil
}

func (s *Service) GetProducts(ctx context.Context, req model.GetProductsRequest) (res []model.GetProductResponse, err error) {
	var productIDPtr *int64
	if req.ProductID != "" {
		if productIDInt, err := strconv.Atoi(req.ProductID); err == nil {
			productIDPtr = utils.ToPointer(int64(productIDInt))
		}
	}

	limitInt, _ := strconv.Atoi(req.Limit)
	offsetInt, _ := strconv.Atoi(req.Offset)
	filter := model.ProductFilter{
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

	res = make([]model.GetProductResponse, 0, len(products))
	for _, p := range products {
		res = append(res, model.GetProductResponse{
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

func (s *Service) DeleteProduct(ctx context.Context, req model.DeleteProductRequest) (err error) {
	productIDInt := 0
	if req.ProductID != "" {
		productIDInt, err = strconv.Atoi(req.ProductID)
		if err != nil {
			return err
		}
	}

	err = s.repository.DeleteProductByID(ctx, int64(productIDInt))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return constants.ErrProductNotFound
		}
		return err
	}

	return
}

func (s *Service) UpdateProduct(ctx context.Context, req model.PutProductRequest) (res model.PutProductResponse, err error) {
	//
	// Get File Uri & Thumbnail Uri
	//
	var fileUri, fileThumbnailUri string
	if req.FileID != "" {
		file, err := s.repository.GetFileByFileID(ctx, req.FileID)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return res, constants.ErrFileIDNotValid
			}
			return res, err
		}
		fileUri = file.Uri
		fileThumbnailUri = file.ThumbnailUri
	}

	//
	// Set Value for Update Product
	//
	productIDInt := 0
	if req.ProductID != "" {
		productIDInt, err = strconv.Atoi(req.ProductID)
		if err != nil {
			return res, err
		}
	}

	updateData := model.Product{
		ID:               int64(productIDInt),
		Name:             req.Name,
		Category:         utils.ToPointer(req.Category),
		Qty:              req.Qty,
		Price:            req.Price,
		SKU:              req.Sku,
		FileID:           utils.ToPointer(req.FileID),
		FileURI:          utils.ToPointer(fileUri),
		FileThumbnailURI: utils.ToPointer(fileThumbnailUri),
	}

	//
	// Update Product
	//
	val, err := s.repository.UpdateProduct(ctx, updateData)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return res, constants.ErrProductNotFound
		}
		if errors.Is(err, constants.ErrDuplicateSKU) {
			return res, errors.New("sku already exists (per account basis)")
		}
		return res, err
	}

	res = model.PutProductResponse{
		ProductID:        strconv.Itoa(int(val.ID)),
		Name:             val.Name,
		Category:         utils.PointerValue(val.Category, ""),
		Qty:              val.Qty,
		Price:            val.Price,
		Sku:              val.SKU,
		FileID:           utils.PointerValue(val.FileID, ""),
		FileUri:          utils.PointerValue(val.FileURI, ""),
		FileThumbnailUri: utils.PointerValue(val.FileThumbnailURI, ""),
		CreatedAt:        val.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:        val.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	return res, nil
}
