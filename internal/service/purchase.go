package service

import (
	"WeenieHut/internal/model"
	"WeenieHut/internal/repository"
	"context"
)

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

func (s *Service) PushCart(ctx context.Context, cart model.StoreCart) (int64, error) {
	params := repository.InsertCartRow{
		TotalPrice:          cart.TotalPrice,
		SenderName:          cart.SenderName,
		SenderContactType:   cart.SenderContactType,
		SenderContactDetail: cart.SenderContactDetail,
	}

	cartId, err := s.repository.InsertCart(ctx, params)
	if err != nil {
		return 0, err
	}

	return cartId, nil
}

func (s *Service) PushCartItem(ctx context.Context, cartItem model.StoreCartItems) error {
	params := repository.InsertCartItemRow{
		CartID:    cartItem.CartID,
		SellerID:  cartItem.SellerID,
		ProductID: cartItem.ProductID,
		Qty:       cartItem.Qty,
		Price:     cartItem.Price,
	}

	_, err := s.repository.InsertCartItem(ctx, params)
	if err != nil {
		return err
	}

	return nil
}
