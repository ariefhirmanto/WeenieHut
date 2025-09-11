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
type PurchaseCartRequest struct {
	PurchasedItems      []PurchasedItem `json:"purchasedItems" validate:"required,min=1,dive"`
	SenderName          string          `json:"senderName" validate:"required,min=4,max=55"`
	SenderContactType   string          `json:"senderContactType" validate:"required,oneof=email phone"`
	SenderContactDetail string          `json:"senderContactDetail" validate:"required"`
}

type PurchasedItem struct {
	ProductID string `json:"productId" validate:"required"`
	Qty       int    `json:"qty" validate:"required,min=2"`
}
