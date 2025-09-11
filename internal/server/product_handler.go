package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
)

func (s *Server) productHandler(w http.ResponseWriter, r *http.Request) {
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

	res, err := s.service.PostProduct(ctx, req)
	if err != nil {
		log.Println("failed to create product:", err)
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	sendResponse(w, http.StatusCreated, res)
}
