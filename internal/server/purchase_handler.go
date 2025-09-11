package server

import (
	"WeenieHut/internal/constants"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
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

	for _, item := range PurchaseCartRequest.PurchasedItems {
		// query item.ProductID dari repo product, cek valid kalo bukan return 400
		// cek quantity dari repo product, kalo kurang return 400
		// kalo valid store ke array
		fmt.Println(item.ProductID)
	}

	// cek repo seller, convert sender name, type and detail ke seller id

	// kirim data ke repo purchase: array of product dari loop diatas dan seller id
	// resp := PurchaseResponse{}

}
