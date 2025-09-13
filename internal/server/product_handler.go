package server

import (
	"WeenieHut/internal/constants"
	"WeenieHut/internal/model"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strings"
)

func (s *Server) postProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		sendErrorResponse(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	ct := r.Header.Get("Content-Type")
	if ct == "" || !strings.HasPrefix(ct, "application/json") {
		sendErrorResponse(w, http.StatusBadRequest, "invalid content type")
		return
	}

	ctx := r.Context()
	var req PostProductRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if err := s.validator.Struct(req); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := s.service.PostProduct(ctx, model.PostProductRequest(req))
	if err != nil {
		log.Println("failed to create product:", err)

		switch {
		case errors.Is(err, constants.ErrFileIDNotValid):
			sendErrorResponse(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, constants.ErrDuplicateSKU):
			sendErrorResponse(w, http.StatusConflict, err.Error())
			sendErrorResponse(w, http.StatusUnauthorized, err.Error())
		default:
			sendErrorResponse(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	sendResponse(w, http.StatusCreated, res)
}

func (s *Server) getProductsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		sendErrorResponse(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	query := r.URL.Query()

	req := GetProductsRequest{
		ProductID: query.Get("productId"),
		Sku:       query.Get("sku"),
		Category:  query.Get("category"),
		SortBy:    query.Get("sortBy"),
		Limit:     query.Get("limit"),
		Offset:    query.Get("offset"),
	}

	products, err := s.service.GetProducts(r.Context(), model.GetProductsRequest(req))
	if err != nil {
		if err.Error() == "no products found" {
			sendErrorResponse(w, http.StatusNotFound, err.Error())
		} else {
			sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	var res []GetProductResponse
	for _, product := range products {
		res = append(res, GetProductResponse(product))
	}
	sendResponse(w, http.StatusOK, res)
}

func (s *Server) deleteProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		sendErrorResponse(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	query := r.URL.Query()
	req := DeleteProductRequest{
		ProductID: query.Get("productId"),
	}

	err := s.service.DeleteProduct(r.Context(), model.DeleteProductRequest(req))
	if err != nil {
		if errors.Is(err, constants.ErrProductNotFound) {
			sendErrorResponse(w, http.StatusNotFound, err.Error())
		} else {
			sendErrorResponse(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	sendResponse(w, http.StatusOK, map[string]string{"message": "product deleted"})
}

func (s *Server) updateProductHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		sendErrorResponse(w, http.StatusMethodNotAllowed, "method not allowed")
		return
	}

	ct := r.Header.Get("Content-Type")
	if ct == "" || !strings.HasPrefix(ct, "application/json") {
		sendErrorResponse(w, http.StatusBadRequest, "invalid content type")
		return
	}

	query := r.URL.Query()
	req := PutProductRequest{
		ProductID: query.Get("productId"),
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	if err := s.validator.Struct(req); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := s.service.UpdateProduct(r.Context(), model.PutProductRequest(req))
	if err != nil {
		switch {
		case errors.Is(err, constants.ErrFileIDNotValid):
			sendErrorResponse(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, constants.ErrProductNotFound):
			sendErrorResponse(w, http.StatusNotFound, err.Error())
		case errors.Is(err, constants.ErrDuplicateSKU):
			sendErrorResponse(w, http.StatusConflict, err.Error())
		default:
			sendErrorResponse(w, http.StatusInternalServerError, "internal server error")
		}
		return
	}

	sendResponse(w, http.StatusOK, PutProductResponse(res))
}
