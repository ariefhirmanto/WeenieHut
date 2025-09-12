package repository

import (
	"context"
	"time"
)

// SELECT id, name, category, qty, price, sku, file_id, file_uri, file_thumbnail_uri, created_at, updated_at
// FROM product
// WHERE user_id = $1

//     id BIGSERIAL PRIMARY KEY,
//     user_id BIGINT NOT NULL,
//     name VARCHAR(255) NOT NULL,
//     category VARCHAR(100),
//     qty INTEGER NOT NULL CHECK (qty >= 0),
//     price NUMERIC(12,2) NOT NULL CHECK (price >= 0),
//     sku VARCHAR(100) NOT NULL,
//     file_id VARCHAR(255),
//     file_uri TEXT,
//     file_thumbnail_uri TEXT,
//     created_at TIMESTAMP DEFAULT NOW(),
//     updated_at TIMESTAMP DEFAULT NOW(),

const selectProductByProductId = `-- SelectProductByProductId
SELECT user_id as userID, name, category, qty, price, sku, file_id as fileID, file_uri as fileUri, file_thumbnail_uri as fileThumbnailUri, created_at as createdAt, updated_at as updatedAt
FROM product
WHERE id = $1`

type SelectProductByProductIdRow struct {
	UserID           int64
	Name             string
	Category         *string
	Qty              int
	Price            int64
	SKU              string
	FileID           *string
	FileURI          *string
	FileThumbnailURI *string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (q *Queries) SelectProductByProductId(ctx context.Context, productIdInput int64) (SelectProductByProductIdRow, error) {
	row := q.db.QueryRowContext(ctx, selectProductByProductId, productIdInput)
	var i SelectProductByProductIdRow
	err := row.Scan(
		&i.UserID,
		&i.Name,
		&i.Category,
		&i.Qty,
		&i.Price,
		&i.SKU,
		&i.FileID,
		&i.FileURI,
		&i.FileThumbnailURI,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const selectPaymentDetailByUserId = `-- SelectPaymentDetailByUserId
SELECT bank_account_name, bank_account_holder, bank_account_number 
FROM users 
WHERE id = $1`

type SelectPaymentDetailByUserIdRow struct {
	BankAccountName   string
	BankAccountHolder string
	BankAccountNumber int64
}

func (q *Queries) SelectPaymentDetailByUserId(ctx context.Context, userId int64) (SelectPaymentDetailByUserIdRow, error) {
	row := q.db.QueryRowContext(ctx, selectPaymentDetailByUserId, userId)
	var i SelectPaymentDetailByUserIdRow
	err := row.Scan(
		&i.BankAccountName,
		&i.BankAccountHolder,
		&i.BankAccountNumber,
	)
	return i, err
}

const insertCart = `-- name: InsertCart :one
INSERT INTO carts (
    total_price,
    sender_name,
    sender_contact_type,
    sender_contact_detail
) VALUES ($1, $2, $3, $4)
RETURNING id`

type InsertCartRow struct {
	TotalPrice          int64
	SenderName          string
	SenderContactType   string
	SenderContactDetail string
}

func (q *Queries) InsertCart(ctx context.Context, arg InsertCartRow) (int64, error) {
	var id int64
	err := q.db.QueryRowContext(ctx, insertCart,
		arg.TotalPrice,
		arg.SenderName,
		arg.SenderContactType,
		arg.SenderContactDetail,
	).Scan(&id)
	return id, err
}

const insertCartItem = `-- name: InsertCartItem :one
INSERT INTO cart_items (
    cart_id,
    seller_id,
    product_id,
    qty,
    price
) VALUES ($1, $2, $3, $4, $5)`

type InsertCartItemRow struct {
	CartID    int64
	SellerID  int64
	ProductID int64
	Qty       int
	Price     int64
}

func (q *Queries) InsertCartItem(ctx context.Context, arg InsertCartItemRow) (int64, error) {
	var id int64
	err := q.db.QueryRowContext(ctx, insertCartItem,
		arg.CartID,
		arg.SellerID,
		arg.ProductID,
		arg.Qty,
		arg.Price,
	).Scan(&id)
	return id, err
}
