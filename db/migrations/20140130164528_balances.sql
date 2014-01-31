-- +goose Up
CREATE TABLE balances (
    uuid uuid NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_uuid uuid NOT NULL,
    currency int NOT NULL,
    available_balance BIGINT NOT NULL,
    reserved_balance BIGINT NOT NULL
);


-- +goose Down
DROP TABLE balances;