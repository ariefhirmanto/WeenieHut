package repository

import (
	"WeenieHut/internal/model"
	"context"
)

const (
	insertUserQuery                   = `INSERT INTO users (email, phone, password_hash) VALUES ($1, $2, $3) RETURNING id`
	selectUserCredentialsByEmailQuery = `SELECT id, phone, password_hash FROM users WHERE email = $1`
	selectUserCredentialsByPhoneQuery = `SELECT id, email, password_hash FROM users WHERE phone = $1`
)

func (q *Queries) InsertUser(ctx context.Context, user model.User, passwordHash string) (res model.User, err error) {
	err = q.db.QueryRowContext(ctx, insertUserQuery,
		user.Email,
		user.Phone,
		passwordHash,
	).Scan(&user.ID)

	if err != nil {
		return model.User{}, err
	}

	return user, nil
}

func (q *Queries) SelectUserCredentialsByEmail(ctx context.Context, email string) (res model.User, err error) {
	err = q.db.QueryRowContext(ctx, selectUserCredentialsByEmailQuery, email).Scan(&res.ID, &res.Phone, &res.PasswordHash)
	if err != nil {
		return model.User{}, err
	}
	return res, nil
}

func (q *Queries) SelectUserCredentialsByPhone(ctx context.Context, email string) (res model.User, err error) {
	err = q.db.QueryRowContext(ctx, selectUserCredentialsByPhoneQuery, email).Scan(&res.ID, &res.Email, &res.PasswordHash)
	if err != nil {
		return model.User{}, err
	}
	return res, nil
}
