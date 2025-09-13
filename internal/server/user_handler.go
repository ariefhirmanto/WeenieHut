package server

import (
	"WeenieHut/internal/constants"
	"WeenieHut/internal/service"
	"WeenieHut/internal/utils"
	"encoding/json"
	"net/http"
)

func (s *Server) getUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _ := utils.GetUserIDFromCtx(ctx)
	user, err := s.service.GetUserProfile(ctx, userID)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	resp, err := s.responseBuilder.BuildUserResponse(user)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	sendResponse(w, http.StatusOK, resp)
}

func (s *Server) updateUserProfileHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _ := utils.GetUserIDFromCtx(ctx)

	var req UpdateUserProfileRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "invalid request")
		return
	}

	if err := s.userValidator.ValidateUpdateProfileRequest(req); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	fileId, err := s.userValidator.ParseFileID(req.FileID)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "invalid file ID")
		return
	}

	params := service.UpdateUserParams{
		UserID:            userID,
		FileID:            fileId,
		BankAccountName:   req.BankAccountName,
		BankAccountHolder: req.BankAccountHolder,
		BankAccountNumber: req.BankAccountNumber,
	}

	user, err := s.service.UpdateUserProfile(ctx, params)
	if err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := s.responseBuilder.BuildUserResponse(user)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	sendResponse(w, http.StatusOK, resp)
}

func (s *Server) updateUserContactHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _ := utils.GetUserIDFromCtx(ctx)

	var req UpdateUserContactRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, "invalid request")
		return
	}

	if err := s.userValidator.ValidateUpdateContactRequest(req); err != nil {
		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	params := service.UpdateUserParams{
		UserID: userID,
		Phone:  req.Phone,
		Email:  req.Email,
	}

	user, err := s.service.UpdateUserContact(ctx, params)
	if err != nil {
		if err == constants.ErrDuplicatePhoneNum {
			sendErrorResponse(w, http.StatusConflict, err.Error())
			return
		}

		if err == constants.ErrDuplicateEmail {
			sendErrorResponse(w, http.StatusConflict, err.Error())
			return
		}

		sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := s.responseBuilder.BuildUserResponse(user)
	if err != nil {
		sendErrorResponse(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	sendResponse(w, http.StatusOK, resp)
}
