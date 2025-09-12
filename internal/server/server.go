package server

import (
	"WeenieHut/internal/model"
	"WeenieHut/internal/constants"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"

	_ "github.com/joho/godotenv/autoload"
)

type Service interface {
	EmailLogin(ctx context.Context, user string, password string) (string, string, error)
	Login(ctx context.Context, email string, password string) (string, error)
	PhoneLogin(ctx context.Context, phone string, password string) (string, string, error)
	Register(ctx context.Context, user model.User, password string) (string, error)

	UploadFile(ctx context.Context, file io.Reader, filename string, sizeInBytes int64) (model.File, error)
	PostProduct(ctx context.Context, req PostProductRequest) (res PostProductResponse, err error)
	GetProducts(ctx context.Context, req GetProductsRequest) (res []GetProductResponse, err error)
	UpdateProduct(ctx context.Context, req PutProductRequest) (res PutProductResponse, err error)
	DeleteProduct(ctx context.Context, req DeleteProductRequest) (err error)
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
			if strings.EqualFold(pt, productType) {
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
