package server

import (
	"WeenieHut/internal/constants"
	"WeenieHut/internal/model"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
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

	productSellerDummy := map[string]int{
		"Arief Hirmanto": 1,
		"Arief":          2,
		"Hirmanto":       3,
	}

	seller := make(map[string]string) //{sellerID:price}
	var products []model.ProductCart
	var paymentdetails []model.CartPaymentDetail
	var totalPrices int64
	now := time.Now()
	for _, item := range PurchaseCartRequest.PurchasedItems {
		// query item.ProductID dari repo product, cek valid kalo bukan return 400
		// cek quantity dari repo product, kalo kurang return 400
		fmt.Println(item.ProductID)

		// dapet seller id dari product, harusnya dari db
		// asumsi dapet data dari nama seller
		sellerID := productSellerDummy["Arief Hirmanto"]

		// dari productID harusnya bisa dapet seller id (atau disini user id)
		// kalo valid store

		pc := model.ProductCart{
			ProductID:        "random",
			Name:             "celana dalam",
			Category:         "Clothes",
			Qty:              999 - 1,
			Price:            1000,
			SKU:              "randomSKU",
			FileID:           "randomID",
			FileURI:          "randomURI",
			FileThumbnailURI: "randomThumbnailUri",
			CreatedAt:        now,
			UpdatedAt:        now, //only keep this field, all the items beside this is come from db
		}

		_, exists := seller[sellerID]
		if !exists {
			seller[sellerID] = pc.price
		} else {
			seller[sellerID] += pc.price
		}
		products = append(products, pc)
		totalPrices += pc.Price
	}

	// query repo seller. sender name, type and detail untuk dapet bankaccName, bankAccHolder, bankAccNum

	for _, price := range seller { //should be id and price, or any distinguishable things
		cartPayment := model.CartPaymentDetail{
			BankAccountName:   "Arief Hirmanto",
			BankAccountHolder: "Arief Hirmanto",
			BankAccountNumber: "0210220234",
			TotalPrice:        price,
		}
		paymentdetails = append(paymentdetails, cartPayment)
	}
	// kirim data ke repo purchase: array of product dari loop diatas dan seller id
	// purchase(products, totalPrices, paymentdetails)
	// resp := PurchaseResponse{}

}
