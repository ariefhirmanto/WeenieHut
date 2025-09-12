package model

import (
	"database/sql"

	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	jwt.RegisteredClaims
	UserID int64 `json:"userID"`
}

type User struct {
	ID                int64
	Email             sql.NullString
	Phone             sql.NullString
	PasswordHash      string
	FileID            string
	FileURI           string
	FileThumbnailURI  string
	BankAccountName   string
	BankAccountHolder string
	BankAccountNumber string
}
