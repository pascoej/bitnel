-- +goose Up
CREATE TABLE oauth_tokens (
    uuid uuid NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_uuid uuid NOT NULL,
    access_token TEXT,
    expires_at timestamptz NOT NULL,
    created_at timestamptz NOT NULL
);

ALTER TABLE oauth_tokens ADD CONSTRAINT oauth_tokens_user_uuid_fkey FOREIGN KEY (user_uuid) REFERENCES users (uuid);

-- +goose Down
DROP TABLE oauth_tokens;