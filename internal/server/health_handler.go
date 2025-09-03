package server

import (
	"SaltySpitoon/internal/utils"
	"encoding/json"
	"log"
	"net/http"
)

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := map[string]string{"message": "Hello World"}
	jsonResp, err := json.Marshal(resp)
	log.Println(utils.GetUserIDFromCtx(r.Context()))
	if err != nil {
		http.Error(w, "Failed to marshal response", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(jsonResp); err != nil {
		log.Printf("Failed to write response: %v", err)
	}
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	sendResponse(w, http.StatusOK, "OK")
}
