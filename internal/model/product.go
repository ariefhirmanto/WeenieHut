package model

import "time"

type Product struct {
	ID               int64     `db:"id"`
	Name             string    `db:"name"`
	Category         *string   `db:"category"`
	Qty              int       `db:"qty"`
	Price            float64   `db:"price"`
	SKU              string    `db:"sku"`
	FileID           *string   `db:"file_id"`
	FileURI          *string   `db:"file_uri"`
	FileThumbnailURI *string   `db:"file_thumbnail_uri"`
	CreatedAt        time.Time `db:"created_at"`
	UpdatedAt        time.Time `db:"updated_at"`
}
