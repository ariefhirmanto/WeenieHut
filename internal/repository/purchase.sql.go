package repository

import (
	"WeenieHut/internal/model"
	"context"
	"database/sql"
	"errors"
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
	Price            float64
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
	TotalPrice          float64
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
	Price     float64
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

func (q *Queries) InsertCartTransaction(ctx context.Context, arg model.StoreCart, items map[int64]model.StoreCartItems) (int64, error) {
	// Start the transaction
	tx, err := q.db.BeginTx(ctx, nil)
	if err != nil {
		return 0, err
	}

	// Insert cart and get the inserted id
	var cartID int64
	err = tx.QueryRowContext(ctx, insertCart,
		arg.TotalPrice,
		arg.SenderName,
		arg.SenderContactType,
		arg.SenderContactDetail,
	).Scan(&cartID)

	if err != nil {
		tx.Rollback()
		return 0, err
	}

	// Insert cart items dynamically
	stmt, err := tx.PrepareContext(ctx, insertCartItem)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	defer stmt.Close()

	for _, item := range items {
		_, err := stmt.ExecContext(ctx, cartID, item.SellerID, item.ProductID, item.Qty, item.Price)
		if err != nil {
			tx.Rollback()
			return 0, err
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return 0, err
	}

	return cartID, nil
}

const selectProductsByCartId = `-- SelectProductsByCartId
SELECT
    p.id as productID,
    p.name,
    p.category,
    p.price as originalPrice,
    p.sku,
    p.file_id as fileID,
    p.file_uri as fileUri,
    p.file_thumbnail_uri as fileThumbnailUri,
    ci.qty as cartQty,
    ci.price as cartPrice,
    ci.seller_id as sellerID,
    p.created_at as createdAt,
    p.updated_at as updatedAt
FROM cart_items ci
JOIN product p ON ci.product_id = p.id
WHERE ci.cart_id = $1`

type SelectProductsByCartIdRow struct {
	ProductID        int64
	Name             string
	Category         *string
	OriginalPrice    float64
	SKU              string
	FileID           *string
	FileURI          *string
	FileThumbnailURI *string
	CartQty          int
	CartPrice        float64
	SellerID         int64
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

func (q *Queries) SelectProductsByCartId(ctx context.Context, cartId int64) ([]SelectProductsByCartIdRow, error) {
	rows, err := q.db.QueryContext(ctx, selectProductsByCartId, cartId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []SelectProductsByCartIdRow
	for rows.Next() {
		var i SelectProductsByCartIdRow
		if err := rows.Scan(
			&i.ProductID,
			&i.Name,
			&i.Category,
			&i.OriginalPrice,
			&i.SKU,
			&i.FileID,
			&i.FileURI,
			&i.FileThumbnailURI,
			&i.CartQty,
			&i.CartPrice,
			&i.SellerID,
			&i.CreatedAt,
			&i.UpdatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const selectCartById = `-- SelectCartById
SELECT id, total_price, sender_name, sender_contact_type, sender_contact_detail
FROM carts
WHERE id = $1`

type SelectCartByIdRow struct {
	ID                  int64
	TotalPrice          int64
	SenderName          *string
	SenderContactType   *string
	SenderContactDetail *string
}

func (q *Queries) SelectCartById(ctx context.Context, cartId int64) (SelectCartByIdRow, error) {
	row := q.db.QueryRowContext(ctx, selectCartById, cartId)
	var i SelectCartByIdRow
	err := row.Scan(
		&i.ID,
		&i.TotalPrice,
		&i.SenderName,
		&i.SenderContactType,
		&i.SenderContactDetail,
	)
	return i, err
}

func (q *Queries) CartExists(ctx context.Context, cartID int64) (bool, error) {
	query := `
		SELECT 1 FROM carts
		WHERE id = $1
	`

	var exists int
	err := q.db.QueryRowContext(ctx, query, cartID).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, err
	}

	return true, nil
}
