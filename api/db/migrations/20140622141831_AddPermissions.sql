
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
ALTER TABLE oauth_tokens ADD COLUMN scope TEXT;

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back

ALTER TABLE oauth_tokens DROP COLUMN scope;