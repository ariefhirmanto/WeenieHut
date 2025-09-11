package server

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type PostProductRequest struct {
	Name     string  `json:"name" validate:"required,min=4,max=32"`
	Category string  `json:"category" validate:"required,productType"`
	Qty      int     `json:"qty" validate:"required,min=1"`
	Price    float64 `json:"price" validate:"required,gt=0"`
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

type DeleteProductRequest struct {
	ProductID string `query:"productId"`
}
