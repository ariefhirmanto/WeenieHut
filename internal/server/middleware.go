package server

import (
	"WeenieHut/internal/constants"
	"WeenieHut/internal/utils"
	"context"
	"net/http"
	"strings"
)

func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/health" || path == "/v1/login/email" || path == "/v1/login/phone" || path == "/v1/register/email" || path == "/v1/register/phone" {
			next.ServeHTTP(w, r)
			return
		}

		authorizationHeader := r.Header.Get("Authorization")
		if !strings.Contains(authorizationHeader, "Bearer") {
			sendErrorResponse(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		tokenString := strings.Replace(authorizationHeader, "Bearer ", "", -1)

		userID, err := utils.ParseUserIDFromToken(tokenString)
		if err != nil {
			sendErrorResponse(w, http.StatusUnauthorized, "unauthorized")
			return
		}

		ctx := context.WithValue(r.Context(), constants.UserIDCtxKey, userID)

		// Proceed with the next handler
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
