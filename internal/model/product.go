package model

import "time"

type Product struct {
	ID               int64     `db:"id"`
	Name             string    `db:"name"`
	Category         *string   `db:"category"`
	Qty              int       `db:"qty"`
	Price            float64   `db:"price"`
	SKU              string    `db:"sku"`
	FileID           *string   `db:"file_id"`
	FileURI          *string   `db:"file_uri"`
	FileThumbnailURI *string   `db:"file_thumbnail_uri"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}

type ProductFilter struct {
	ProductID *int64
	Sku       string
	Category  string
	SortBy    string
	Limit     int
	Offset    int
}

type PostProductRequest struct {
	Name     string  `json:"name" validate:"required,min=4,max=32"`
	Category string  `json:"category" validate:"required,productType"`
	Qty      int     `json:"qty" validate:"required,min=1"`
	Price    float64 `json:"price" validate:"required,gte=100"`
	Sku      string  `json:"sku" validate:"required,min=1,max=32"`
	FileID   string  `json:"fileId" validate:"required"`
}

type GetProductsRequest struct {
	ProductID string `query:"productId"`
	Sku       string `query:"sku"`
	Category  string `query:"category"`
	SortBy    string `query:"sortBy"`
	Limit     string `query:"limit"`
	Offset    string `query:"offset"`
}

type PutProductRequest struct {
	ProductID string  `query:"productId"`
	Name      string  `json:"name" validate:"required,min=4,max=32"`
	Category  string  `json:"category" validate:"required,productType"`
	Qty       int     `json:"qty" validate:"required,min=1"`
	Price     float64 `json:"price" validate:"required,gte=100"`
	Sku       string  `json:"sku" validate:"required,min=1,max=32"`
	FileID    string  `json:"fileId" validate:"required"`
}

type DeleteProductRequest struct {
	ProductID string `query:"productId"`
}

type PostProductResponse struct {
	ProductID        string  `json:"productId"`
	Name             string  `json:"name"`
	Category         string  `json:"category"`
	Qty              int     `json:"qty"`
	Price            float64 `json:"price"`
	Sku              string  `json:"sku"`
	FileID           string  `json:"fileId"`
	FileUri          string  `json:"fileUri"`
	FileThumbnailUri string  `json:"fileThumbnailUri"`
	CreatedAt        string  `json:"createdAt"`
	UpdatedAt        string  `json:"updatedAt"`
}

type GetProductResponse struct {
	ProductID        string  `json:"productId"`
	Name             string  `json:"name"`
	Category         string  `json:"category"`
	Qty              int     `json:"qty"`
	Price            float64 `json:"price"`
	Sku              string  `json:"sku"`
	FileID           string  `json:"fileId"`
	FileUri          string  `json:"fileUri"`
	FileThumbnailUri string  `json:"fileThumbnailUri"`
	CreatedAt        string  `json:"createdAt"`
	UpdatedAt        string  `json:"updatedAt"`
}

type PutProductResponse struct {
	ProductID        string  `json:"productId"`
	Name             string  `json:"name"`
	Category         string  `json:"category"`
	Qty              int     `json:"qty"`
	Price            float64 `json:"price"`
	Sku              string  `json:"sku"`
	FileID           string  `json:"fileId"`
	FileUri          string  `json:"fileUri"`
	FileThumbnailUri string  `json:"fileThumbnailUri"`
	CreatedAt        string  `json:"createdAt"`
	UpdatedAt        string  `json:"updatedAt"`
}
