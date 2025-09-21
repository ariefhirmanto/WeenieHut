package server

import (
	"WeenieHut/internal/constants"
	"WeenieHut/internal/model"
	"WeenieHut/internal/service"
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
	PhoneLogin(ctx context.Context, phone string, password string) (string, string, error)
	Register(ctx context.Context, user model.User, password string) (string, error)

	GetUserProfile(ctx context.Context, userId int64) (model.User, error)
	UpdateUserProfile(ctx context.Context, params service.UpdateUserParams) (model.User, error)
	UpdateUserContact(ctx context.Context, params service.UpdateUserParams) (model.User, error)
	IsUserExist(ctx context.Context, userID int64) (bool, error)

	UploadFile(ctx context.Context, file io.Reader, filename string, sizeInBytes int64) (model.File, error)
	PostProduct(ctx context.Context, req model.PostProductRequest) (res model.PostProductResponse, err error)
	GetProducts(ctx context.Context, req model.GetProductsRequest) (res []model.GetProductResponse, err error)
	UpdateProduct(ctx context.Context, req model.PutProductRequest) (res model.PutProductResponse, err error)
	DeleteProduct(ctx context.Context, req model.DeleteProductRequest) (err error)
	GetProductByProductId(ctx context.Context, productIdInput int64) (model.ProductCart, int64, error)
	GetSellerPaymentDetailBySellerId(ctx context.Context, sellerID int64) (model.CartPaymentDetail, error)
	PushCart(ctx context.Context, cart model.StoreCart) (int64, error)
	PushCartItem(ctx context.Context, cartItem model.StoreCartItems) error
	PushCartAndItems(ctx context.Context, cart model.StoreCart, items map[int64]model.StoreCartItems) (int64, error)

	PurchasePayment(ctx context.Context, purchaseId string, fileIds []string) error
}

type Server struct {
	port            int
	service         Service
	validator       *validator.Validate
	userValidator   *UserValidator
	responseBuilder *UserResponseBuilder
}

func NewServer(service Service) *http.Server {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	v := validator.New()
	NewServer := &Server{
		port:            port,
		service:         service,
		validator:       v,
		userValidator:   NewUserValidator(v),
		responseBuilder: NewUserResponseBuilder(),
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
