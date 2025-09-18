package service

import (
	"WeenieHut/internal/constants"
	"WeenieHut/internal/model"
	"WeenieHut/internal/repository"
	"context"
	"strconv"
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

func (s *Service) PushCartAndItems(ctx context.Context, cart model.StoreCart, items map[int64]model.StoreCartItems) (int64, error) {
	cartID, err := s.repository.InsertCartTransaction(ctx, cart, items)

	if err != nil {
		return 0, err
	}

	return cartID, nil
}

func (s *Service) PurchasePayment(ctx context.Context, purchaseId string, fileIds []string) error {
	// Parse purchaseId to get cartId (assuming purchaseId is actually cartId as string)
	cartId, err := strconv.ParseInt(purchaseId, 10, 64)
	if err != nil {
		return constants.ErrInternalServer
	}

	// Get Cart
	exists, err := s.repository.CartExists(ctx, cartId)
	if err != nil {
		return err
	}

	if !exists {
		return constants.ErrPurchaseNotFound
	}

	// Get products by cartId
	cartProducts, err := s.repository.SelectProductsByCartId(ctx, cartId)
	if err != nil {
		return err
	}

	if len(cartProducts) == 0 {
		return constants.ErrProductNotFound
	}

	// Group products by seller and validate against uploaded files
	sellerProducts := make(map[int64][]repository.SelectProductsByCartIdRow)
	expectedFileCount := 0

	for _, product := range cartProducts {
		sellerProducts[product.SellerID] = append(sellerProducts[product.SellerID], product)
		expectedFileCount++
	}

	// Validate file count matches cart items count
	if len(fileIds) != expectedFileCount {
		return constants.ErrNotEqualAvailableSellersInCart
	}

	// Check file existence
	for _, fileId := range fileIds {
		exists, err := s.repository.FileExists(ctx, fileId)
		if err != nil {
			return err
		}

		if !exists {
			return constants.ErrFileNotFound
		}
	}

	return nil
}
