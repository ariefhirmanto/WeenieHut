package server

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type ErrorResponse struct {
	Error string `json:"error,omitempty"`
}

func sendResponse(w http.ResponseWriter, statusCode int, body interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if body == nil {
		return
	}

	if err := json.NewEncoder(w).Encode(body); err != nil {
		log.Printf("failed to encode response: %v", err)
	}
}

func sendErrorResponse(w http.ResponseWriter, statusCode int, error string) {
	resp := ErrorResponse{
		Error: error,
	}

	sendResponse(w, statusCode, resp)
}

type LoginResponse struct {
	Email string `json:"email"`
	Token string `json:"token"`
}

type RegisterResponse struct {
	Email string `json:"email"`
	Token string `json:"token"`
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

type PurchasedItemResponse struct {
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

type PaymentDetailResponse struct {
	BankAccountName   string `json:"bankAccountName"`
	BankAccountHolder string `json:"bankAccountHolder"`
	BankAccountNumber string `json:"bankAccountNumber"`
	TotalPrice        int64  `json:"totalPrice"` // Total for this seller
}

type PurchaseResponse struct {
	PurchaseID     string                  `json:"purchaseId"`     // Any ID
	PurchasedItems []PurchasedItemResponse `json:"purchasedItems"` // Must be at least 1
	TotalPrice     int64                   `json:"totalPrice"`     // Sum of all item prices
	PaymentDetails []PaymentDetailResponse `json:"paymentDetails"` // One per seller involved in transaction
}
