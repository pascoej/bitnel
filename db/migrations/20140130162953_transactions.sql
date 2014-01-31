-- +goose Up
CREATE TABLE transactions (
    uuid uuid NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    balance_uuid uuid NOT NULL,
    type int NOT NULL,
    amount BIGINT NOT NULL,
    fee_amount BIGINT NOT NULL,
    created_at timestamptz NOT NULL DEFAULT NOW(),  
    trade uuid
);


-- +goose Down
DROP TABLE transactions;