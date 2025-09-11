package model

import "time"

type ProductCart struct {
	ProductID        string    `json:"productId"` // Any ID
	Name             string    `json:"name"`
	Category         string    `json:"category"`
	Qty              int       `json:"qty"`   // Quantity before bought
	Price            int64     `json:"price"` // Price per item
	SKU              string    `json:"sku"`
	FileID           string    `json:"fileId"`
	FileURI          string    `json:"fileUri"`          // Related file URI
	FileThumbnailURI string    `json:"fileThumbnailUri"` // Related thumbnail URI
	CreatedAt        time.Time `json:"createdAt"`
	UpdatedAt        time.Time `json:"updatedAt"`
}

type CartPaymentDetail struct {
	BankAccountName   string `json:"bankAccountName"`
	BankAccountHolder string `json:"bankAccountHolder"`
	BankAccountNumber string `json:"bankAccountNumber"`
	TotalPrice        int64  `json:"totalPrice"` // Total for this seller
}

type PurchaseCartStore struct {
	PurchasedItems []PurchaseCartItem `json:"purchasedItems"`
	TotalPrice     int64              `json:"totalPrice`
	PaymentDetails []PaymentDetail    `json:"paymentDetails"`
}

type PurchaseCartReturn struct {
	PurchaseID     string             `json:"purchaseID"`
	PurchasedItems []PurchaseCartItem `json:"purchasedItems"`
	TotalPrice     int64              `json:"totalPrice`
	PaymentDetails []PaymentDetail    `json:"paymentDetails"`
}

// can be adjusted, start with a simple item details
type PurchaseCartItem struct {
	ProductID string `json:"productId"`
	Qty       int    `json:"qty"`
}
