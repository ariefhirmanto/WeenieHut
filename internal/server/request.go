package server

type EmailLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type PhoneLoginRequest struct {
	Phone    string `json:"phone" validate:"required,e164"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type EmailRegisterRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type PhoneRegisterRequest struct {
	Phone    string `json:"phone" validate:"required,e164"`
	Password string `json:"password" validate:"required,min=8,max=32"`
}

type UpdateUserProfileRequest struct {
	FileID            string `json:"fileId"`
	BankAccountName   string `json:"bankAccountName" validate:"omitempty,min=4,max=32"`
	BankAccountHolder string `json:"bankAccountHolder" validate:"omitempty,min=4,max=32"`
	BankAccountNumber string `json:"bankAccountNumber" validate:"omitempty,min=4,max=32"`
}

type UpdateUserContactRequest struct {
	Email string `json:"email" validate:"omitempty"`
	Phone string `json:"phone" validate:"omitempty"`
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

type PurchasePaymentRequest struct {
	FileIDs []string `json:"fileIds" validate:"required"`
}
