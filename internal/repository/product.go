package repository

import (
	"WeenieHut/internal/model"
	"context"
	"fmt"
)

func (q *Queries) InsertProduct(ctx context.Context, data model.Product) (model.Product, error) {
	query := `
		INSERT INTO product (
			name, category, qty, price, sku, file_id, file_uri, file_thumbnail_uri, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW()
		)
		RETURNING id, created_at, updated_at
	`

	inserted := model.Product{}
	inserted = data

	err := q.db.QueryRowContext(ctx, query,
		data.Name,
		data.Category,
		data.Qty,
		data.Price,
		data.SKU,
		data.FileID,
		data.FileURI,
		data.FileThumbnailURI,
	).Scan(&inserted.ID, &inserted.CreatedAt, &inserted.UpdatedAt)

	if err != nil {
		return model.Product{}, fmt.Errorf("insert product: %w", err)
	}

	return inserted, nil
}
