package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID                int64
	Email             sql.NullString
	Phone             sql.NullString
	PasswordHash      string
	Name              string
	FileID            int64
	FileURI           string
	FileThumbnailURI  string
	BankAccountName   string
	BankAccountHolder string
	BankAccountNumber string
	CreatedAt         time.Time
	UpdatedAt         time.Time
}
