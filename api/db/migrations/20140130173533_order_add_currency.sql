-- +goose Up
ALTER TABLE orders ADD COLUMN currency int NOT NULL;

-- +goose Down
ALTER TABLE orders DROP COLUMN currency;