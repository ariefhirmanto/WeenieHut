-- name: SelectUserByEmail :one
SELECT id, password_hash FROM users where email = $1;

-- name: CreateUser :exec
INSERT INTO users (email, password_hash) VALUES ($1, $2);