package server

import (
	"SaltySpitoon/internal/constants"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
)

func (s *Server) loginHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req LoginRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println("invalid login request")
		sendErrorResponse(w, http.StatusBadRequest, "invalid request")
		return
	}

	err = s.validator.Struct(req)
	if err != nil {
		log.Printf("invalid login request: %s\n", err.Error())
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	token, err := s.service.Login(ctx, req.Email, req.Password)
	if err != nil {
		log.Printf("failed to login: %s\n", err.Error())
		if errors.Is(err, constants.ErrUserNotFound) {
			sendErrorResponse(w, http.StatusNotFound, fmt.Sprintf("email %s not found", req.Email))
			return
		}
		if errors.Is(err, constants.ErrUserWrongPassword) {
			sendErrorResponse(w, http.StatusBadRequest, "wrong password")
			return
		}
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := LoginResponse{
		Email: req.Email,
		Token: token,
	}
	sendResponse(w, http.StatusOK, resp)
}

func (s *Server) registerHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	var req RegisterRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		log.Println("invalid login request")
		sendErrorResponse(w, http.StatusBadRequest, "invalid request")
		return
	}

	err = s.validator.Struct(req)
	if err != nil {
		log.Printf("invalid register request: %s\n", err.Error())
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	token, err := s.service.Register(ctx, req.Email, req.Password)
	if err != nil {
		log.Printf("failed to register: %s\n", err.Error())
		if errors.Is(err, constants.ErrEmailAlreadyExists) {
			sendErrorResponse(w, http.StatusConflict, fmt.Sprintf("email %s already exists", req.Email))
			return
		}
		sendErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	resp := RegisterResponse{
		Email: req.Email,
		Token: token,
	}
	sendResponse(w, http.StatusOK, resp)
}
