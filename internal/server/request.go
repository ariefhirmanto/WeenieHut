package server

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type RegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
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
