-- +goose Up
CREATE TABLE sessions (
    uuid uuid NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    token uuid NOT NULL DEFAULT uuid_generate_v4(),
    user_uuid uuid NOT NULL,
    created_at timestamptz NOT NULL DEFAULT NOW(),
    expires_at timestamptz NOT NULL
);


-- +goose Down
DROP TABLE sessions;