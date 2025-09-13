-- +goose Up
-- +goose StatementBegin
CREATE TABLE carts (
    id BIGSERIAL PRIMARY KEY,
    total_price NUMERIC(12,2) NOT NULL CHECK (total_price >= 0),
    sender_name VARCHAR(20),
    sender_contact_type VARCHAR(20),
    sender_contact_detail VARCHAR(20)
);

CREATE TABLE cart_items (
  id BIGSERIAL PRIMARY KEY,
  cart_id BIGINT NOT NULL,
  seller_id BIGINT NOT NULL,
  product_id BIGINT NOT NULL,
  qty INTEGER NOT NULL CHECK (qty >= 0),
  price NUMERIC(12,2) NOT NULL CHECK (price >= 0),

  CONSTRAINT fk_cart
    FOREIGN KEY (cart_id)
    REFERENCES carts(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE,

  CONSTRAINT fk_seller
    FOREIGN KEY (seller_id)
    REFERENCES users(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE,

  CONSTRAINT fk_product
    FOREIGN KEY (product_id)
    REFERENCES product(id)
    ON DELETE CASCADE
    ON UPDATE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS carts;
DROP TABLE IF EXISTS cart_items;
-- +goose StatementEnd
