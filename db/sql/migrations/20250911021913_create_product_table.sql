-- +goose Up
-- +goose StatementBegin
CREATE TABLE product (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT NOT NULL,
    name VARCHAR(255) NOT NULL,
    category VARCHAR(100),                
    qty INTEGER NOT NULL CHECK (qty >= 0),
    price NUMERIC(12,2) NOT NULL CHECK (price >= 0),
    sku VARCHAR(100) NOT NULL,
    file_id VARCHAR(255),
    file_uri TEXT,
    file_thumbnail_uri TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    CONSTRAINT uq_product_user_sku UNIQUE (user_id, sku)
    CONSTRAINT fk_user
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS product;
-- +goose StatementEnd
