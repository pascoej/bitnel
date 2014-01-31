-- +goose Up
CREATE TABLE accounts (
    uuid uuid NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_uuid uuid NOT NULL
);


-- +goose Down
DROP TABLE accounts;