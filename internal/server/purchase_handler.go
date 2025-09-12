package server

import (
	"WeenieHut/internal/model"

	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"
)

func (s *Server) purchaseCartHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req PurchaseCartRequest

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Printf("invalid payload: %s\n", err.Error())
		sendErrorResponse(w, http.StatusBadRequest, "invalid request")
		return
	}

	err = s.validator.Struct(req)
	if err != nil {
		log.Printf("invalid validator: %s\n", err.Error())
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	seller := make(map[int64]model.StoreCartItems) //{sellerID:price}
	var products []model.ProductCart
	var paymentdetails []model.CartPaymentDetail
	var totalPrices int64
	for _, item := range req.PurchasedItems {
		now := time.Now()
		productIDInt, err := strconv.ParseInt(item.ProductID, 10, 64)
		if err != nil {
			log.Printf("invalid productID: %s\n", err.Error())
			sendErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		productInCart, sellerID, err := s.service.GetProductByProductId(ctx, productIDInt)
		if err != nil {
			log.Println("error get product: %s\n", err.Error())
			sendErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}

		if item.Qty <= 0 {
			log.Printf("product qty is out")
			sendErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		if productInCart.Qty < item.Qty {
			log.Printf("bought more than the available qty")
			sendErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		productInCart.UpdatedAt = now
		// ganti ke response type
		products = append(products, productInCart)
		totalPrices += productInCart.Price

		_, exists := seller[sellerID]
		if !exists {
			par := model.StoreCartItems{
				CartID:    0,
				SellerID:  sellerID,
				ProductID: productIDInt,
				Qty:       productInCart.Qty,
				Price:     productInCart.Price,
			}
			seller[sellerID] = par
		} else {
			existing := seller[sellerID]
			existing.Price += productInCart.Price
			seller[sellerID] = existing
		}

	}

	paramsCart := model.StoreCart{
		TotalPrice:          totalPrices,
		SenderName:          req.SenderName,
		SenderContactType:   req.SenderContactType,
		SenderContactDetail: req.SenderContactDetail,
	}

	cartIDfromDB, err := s.service.PushCart(ctx, paramsCart)
	if err != nil {
		log.Println("error storing purchase: %s\n", err.Error())
		sendErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	for sellerID, cartItems := range seller {
		cartItems.CartID = cartIDfromDB
		err := s.service.PushCartItem(ctx, cartItems)
		if err != nil {
			log.Println("error storing purchase: %s\n", err.Error())
			sendErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
			return
		}
		sellerPaymentDetail, err := s.service.GetSellerPaymentDetailBySellerId(ctx, sellerID)
		if err != nil {
			log.Printf("invalid seller: %s\n", err.Error())
			sendErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		sellerPaymentDetail.TotalPrice = cartItems.Price
		// ganti ke response type
		paymentdetails = append(paymentdetails, sellerPaymentDetail)
	}

	// dapet cart id

	// kirim data ke repo purchase: array of product dari loop diatas dan seller id
	// purchase(products, totalPrices, paymentdetails)
	// resp := PurchaseResponse{}
	resp := formatOutput(products, paymentdetails, cartIDfromDB, totalPrices)
	sendResponse(w, http.StatusOK, resp)
}

func formatOutput(products []model.ProductCart, paymentdetails []model.CartPaymentDetail, purchaseID int64, totalP int64) PurchaseResponse {
	var purchasedItemR []PurchasedItemResponse
	var paymentDetailR []PaymentDetailResponse
	for _, product := range products {
		params := PurchasedItemResponse{
			ProductID:        strconv.Itoa(int(product.ProductID)),
			Name:             product.Name,
			Category:         product.Category,
			Qty:              product.Qty,
			Price:            product.Price,
			SKU:              product.SKU,
			FileID:           product.FileID,
			FileURI:          product.FileURI,
			FileThumbnailURI: product.FileThumbnailURI,
			CreatedAt:        product.CreatedAt,
			UpdatedAt:        product.UpdatedAt,
		}
		purchasedItemR = append(purchasedItemR, params)
	}

	for _, paymentD := range paymentdetails {
		paramsP := PaymentDetailResponse{
			BankAccountName:   paymentD.BankAccountName,
			BankAccountHolder: paymentD.BankAccountHolder,
			BankAccountNumber: strconv.Itoa(int(paymentD.BankAccountNumber)),
			TotalPrice:        paymentD.TotalPrice,
		}
		paymentDetailR = append(paymentDetailR, paramsP)
	}

	result := PurchaseResponse{
		PurchaseID:     strconv.Itoa(int(purchaseID)),
		PurchasedItems: purchasedItemR,
		TotalPrice:     totalP,
		PaymentDetails: paymentDetailR,
	}

	return result
}
