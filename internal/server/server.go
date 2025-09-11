package server

import (
	"WeenieHut/internal/constants"
	"context"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"

	_ "github.com/joho/godotenv/autoload"
)

type Service interface {
	Login(ctx context.Context, email string, password string) (string, error)
	Register(ctx context.Context, email string, password string) (string, error)
}

type Server struct {
	port      int
	service   Service
	validator *validator.Validate
}

func NewServer(service Service) *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	NewServer := &Server{
		port:      port,
		service:   service,
		validator: validator.New(),
	}

	// Custom validator for product type
	NewServer.validator.RegisterValidation("productType", func(fl validator.FieldLevel) bool {
		productType := fl.Field().String()
		for _, pt := range constants.ProductTypes {
			if pt == productType {
				return true
			}
		}
		return false
	})

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}
