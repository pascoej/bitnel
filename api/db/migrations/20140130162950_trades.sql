-- +goose Up
CREATE TABLE trades (
    uuid uuid NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    amount BIGINT NOT NULL,
    price BIGINT NOT NULL,
    created_at timestamptz NOT NULL DEFAULT NOW()  
);


-- +goose Down
DROP TABLE trades;