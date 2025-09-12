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

	seller := make(map[int64]int64) //{sellerID:price}
	var products []model.ProductCart
	var paymentdetails []model.CartPaymentDetail
	var totalPrices int64
	now := time.Now()
	for _, item := range req.PurchasedItems {
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

		productInCart.Qty = productInCart.Qty - 1
		productInCart.UpdatedAt = now
		// ganti ke response type
		products = append(products, productInCart)
		totalPrices += productInCart.Price

		_, exists := seller[sellerID]
		if !exists {
			seller[sellerID] = productInCart.Price
		} else {
			seller[sellerID] += productInCart.Price
		}

	}

	for sellerID, price := range seller {
		sellerPaymentDetail, err := s.service.GetSellerPaymentDetailBySellerId(ctx, sellerID)
		if err != nil {
			log.Printf("invalid seller: %s\n", err.Error())
			sendErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		sellerPaymentDetail.TotalPrice = price
		// ganti ke response type
		paymentdetails = append(paymentdetails, sellerPaymentDetail)
	}

	// kirim data ke repo purchase: array of product dari loop diatas dan seller id
	// purchase(products, totalPrices, paymentdetails)
	// resp := PurchaseResponse{}

}
