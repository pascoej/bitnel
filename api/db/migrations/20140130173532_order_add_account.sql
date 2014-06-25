-- +goose Up
ALTER TABLE orders ADD COLUMN account_uuid uuid NOT NULL;

-- +goose Down
ALTER TABLE orders DROP COLUMN account_uuid;