package repository

import (
	"WeenieHut/internal/model"
	"context"
	"fmt"
)

func (q *Queries) InsertFile(ctx context.Context, data model.File) (res model.File, err error) {
	query := `
		INSERT INTO files (
			uri, thumbnail_uri, created_at, updated_at
		) VALUES (
			$1, $2, NOW(), NOW()
		)
		RETURNING id, created_at, updated_at
	`

	err = q.db.QueryRowContext(ctx, query,
		data.Uri,
		data.ThumbnailUri,
	).Scan(&data.ID, &data.CreatedAt, &data.UpdatedAt)

	if err != nil {
		return model.File{}, fmt.Errorf("error inserting file: %w", err)
	}

	return data, nil
}

func (q *Queries) GetFileByFileID(ctx context.Context, fileID string) (res model.File, err error) {
	query := `SELECT * FROM files WHERE id = $1`
	err = q.db.QueryRowContext(ctx, query, fileID).Scan(&res.ID, &res.Uri, &res.ThumbnailUri, &res.SizeInBytes, &res.CreatedAt, &res.UpdatedAt)
	if err != nil {
		return model.File{}, err
	}
	return res, nil
}
