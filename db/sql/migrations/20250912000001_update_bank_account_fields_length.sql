-- +goose Up
-- +goose StatementBegin
ALTER TABLE users 
ALTER COLUMN bank_account_name TYPE VARCHAR(32),
ALTER COLUMN bank_account_holder TYPE VARCHAR(32),
ALTER COLUMN bank_account_number TYPE VARCHAR(32);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE users 
ALTER COLUMN bank_account_name TYPE VARCHAR(20),
ALTER COLUMN bank_account_holder TYPE VARCHAR(20),
ALTER COLUMN bank_account_number NUMERIC(6,2);
-- +goose StatementEnd