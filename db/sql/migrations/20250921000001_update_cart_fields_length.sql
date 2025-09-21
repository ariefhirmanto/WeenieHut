-- +goose Up
-- +goose StatementBegin
ALTER TABLE carts
ALTER COLUMN sender_name TYPE VARCHAR(100),
ALTER COLUMN sender_contact_type TYPE VARCHAR(50),
ALTER COLUMN sender_contact_detail TYPE VARCHAR(100);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE carts
ALTER COLUMN sender_name TYPE VARCHAR(20),
ALTER COLUMN sender_contact_type TYPE VARCHAR(20),
ALTER COLUMN sender_contact_detail TYPE VARCHAR(20);
-- +goose StatementEnd