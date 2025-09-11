package repository

import (
	"WeenieHut/internal/model"
	"WeenieHut/internal/service"
	"context"
	"fmt"
	"strings"
)

func (q *Queries) InsertProduct(ctx context.Context, data model.Product) (res model.Product, err error) {
	query := `
		INSERT INTO product (
			name, category, qty, price, sku, file_id, file_uri, file_thumbnail_uri, created_at, updated_at
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW()
		)
		RETURNING id, created_at, updated_at
	`
	fmt.Println("Executing Query:", query, "with data:", data)
	inserted := data

	err = q.db.QueryRowContext(ctx, query,
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

	res = inserted

	return
}

func (q *Queries) GetProducts(ctx context.Context, filter service.ProductFilter) (res []model.Product, err error) {
	query := `
		SELECT id, name, category, qty, price, sku, file_id, file_uri, file_thumbnail_uri, created_at, updated_at
		FROM product
	`
	conds := []string{}
	args := []interface{}{}
	argIdx := 1

	// filter by productId
	if filter.ProductID != nil {
		conds = append(conds, fmt.Sprintf("id = $%d", argIdx))
		args = append(args, *filter.ProductID)
		argIdx++
	}

	// filter by sku (exact match)
	if filter.Sku != "" {
		conds = append(conds, fmt.Sprintf("sku = $%d", argIdx))
		args = append(args, filter.Sku)
		argIdx++
	}

	// filter by category (case-insensitive exact match)
	if filter.Category != "" {
		conds = append(conds, fmt.Sprintf("LOWER(category) = LOWER($%d)", argIdx))
		args = append(args, filter.Category)
		argIdx++
	}

	// conditions
	if len(conds) > 0 {
		query += " WHERE " + strings.Join(conds, " AND ")
	}

	// sorting
	switch filter.SortBy {
	case "newest":
		query += " ORDER BY updated_at DESC, created_at DESC"
	case "oldest":
		query += " ORDER BY updated_at ASC, created_at ASC"
	case "cheapest":
		query += " ORDER BY price ASC"
	case "expensive":
		query += " ORDER BY price DESC"
	default:
		query += " ORDER BY updated_at DESC, created_at DESC"
	}

	limit := filter.Limit
	if limit <= 0 {
		limit = 5
	}
	offset := filter.Offset
	if offset < 0 {
		offset = 0
	}

	query += fmt.Sprintf(" LIMIT %d OFFSET %d", limit, offset)

	// debug log
	fmt.Println("Executing Query:", query, "Args:", args)

	rows, err := q.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []model.Product
	for rows.Next() {
		var p model.Product
		if err := rows.Scan(
			&p.ID,
			&p.Name,
			&p.Category,
			&p.Qty,
			&p.Price,
			&p.SKU,
			&p.FileID,
			&p.FileURI,
			&p.FileThumbnailURI,
			&p.CreatedAt,
			&p.UpdatedAt,
		); err != nil {
			return nil, err
		}
		products = append(products, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	res = products

	return
}

func (q *Queries) DeleteProductByID(ctx context.Context, id int64) (err error) {
	query := `DELETE FROM product WHERE id = $1`

	_, err = q.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return
}
