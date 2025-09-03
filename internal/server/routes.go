package server

import (
	"SaltySpitoon/internal/constants"
	"SaltySpitoon/internal/utils"
	"context"
	"net/http"
	"strings"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", s.HelloWorldHandler)
	mux.HandleFunc("/health", s.healthHandler)
	mux.HandleFunc("POST /v1/register", s.registerHandler)
	mux.HandleFunc("POST /v1/login", s.loginHandler)

	return s.authMiddleware(mux)
}

func (s *Server) authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		if path == "/health" || path == "/v1/login" || path == "/v1/register" {
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
