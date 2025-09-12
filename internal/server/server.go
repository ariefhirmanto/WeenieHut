package server

import (
	"WeenieHut/internal/model"
	"WeenieHut/internal/service"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
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
