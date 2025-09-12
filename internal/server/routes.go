package server

import (
	"net/http"
)

func (s *Server) RegisterRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/", s.HelloWorldHandler)
	mux.HandleFunc("/health", s.healthHandler)
	mux.HandleFunc("POST /v1/register/email", s.emailRegisterHandler)
	mux.HandleFunc("POST /v1/register/phone", s.phoneRegisterHandler)
	mux.HandleFunc("POST /v1/login/email", s.emailLoginHandler)
	mux.HandleFunc("POST /v1/login/phone", s.phoneLoginHandler)

	mux.HandleFunc("POST /v1/file", s.fileUploadHandler)
	mux.HandleFunc("POST /v1/product", s.postProductHandler)
	mux.HandleFunc("GET /v1/product", s.getProductsHandler)
	mux.HandleFunc("PUT /v1/product/", s.updateProductHandler)
	mux.HandleFunc("DELETE /v1/product/", s.deleteProductHandler)

	return s.authMiddleware(mux)
}
