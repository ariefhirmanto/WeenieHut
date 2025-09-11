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
	Name     string `json:"name" validate:"string,required,min=4,max=32"`
	Category string `json:"category" validate:"string,required,productType"`
	Qty      int    `json:"qty" validate:"numeric,required,min=1"`
	Price    int    `json:"price" validate:"numeric,required,min=100"`
	Sku      string `json:"sku" validate:"string,required,min=0,max=32"`
	FileID   string `json:"fileId" validate:"string,required"`
}
