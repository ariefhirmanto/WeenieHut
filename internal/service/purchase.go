package service

import (
	"WeenieHut/internal/model"
	"context"
)

// push data to purchase db, return purchase id. consist of purchased items {id and qty -> product table}, total price, payment details {user id -> user table}

// func (s *Service) Purchase(products []model.ProductCart, totalPrices int64, paymentDetail []model.CartPaymentDetail) (model.PurchaseCartReturn, error) {
// yang di store
// products.productID
// products.qty
// totalPrices
// paymentDetail.semuanya
// }

func ptrtostring(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

func (s *Service) GetProductByProductId(ctx context.Context, productIdInput int64) (model.ProductCart, int64, error) {
	product, err := s.repository.SelectProductByProductId(ctx, productIdInput)
	// wrong id
	if err != nil {
		return model.ProductCart{}, 0, err
	}

	// validator null (asumsi udah beres buat category, fileid uri thumbnailuri)

	pc := model.ProductCart{
		ProductID:        productIdInput,
		Name:             product.Name,
		Category:         ptrtostring(product.Category),
		Qty:              product.Qty,
		Price:            product.Price,
		SKU:              product.SKU,
		FileID:           ptrtostring(product.FileID),
		FileURI:          ptrtostring(product.FileURI),
		FileThumbnailURI: ptrtostring(product.FileThumbnailURI),
		CreatedAt:        product.CreatedAt,
		UpdatedAt:        product.UpdatedAt,
	}

	return pc, product.UserID, nil
}

// type CartPaymentDetail struct {
// 	BankAccountName   string `json:"bankAccountName"`
// 	BankAccountHolder string `json:"bankAccountHolder"`
// 	BankAccountNumber string `json:"bankAccountNumber"`
// 	TotalPrice        int64  `json:"totalPrice"` // Total for this seller
// }

func (s *Service) GetSellerPaymentDetailBySellerId(ctx context.Context, sellerID int64) (model.CartPaymentDetail, error) {
	sellerDetail, err := s.repository.SelectPaymentDetailByUserId(ctx, sellerID)
	if err != nil {
		return model.CartPaymentDetail{}, err
	}

	spd := model.CartPaymentDetail{
		BankAccountName:   sellerDetail.BankAccountName,
		BankAccountHolder: sellerDetail.BankAccountHolder,
		BankAccountNumber: sellerDetail.BankAccountNumber,
		TotalPrice:        0,
	}

	return spd, nil
}
