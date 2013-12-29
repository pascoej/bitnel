-- +goose Up
CREATE TABLE markets (
    uuid uuid NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    base_currency smallint NOT NULL,
    quote_currency smallint NOT NULL,
    currency_pair text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT NOW()
);

ALTER TABLE orders ADD COLUMN market_uuid uuid NOT NULL;
ALTER TABLE orders ADD CONSTRAINT orders_market_uuid_fkey FOREIGN KEY (market_uuid) REFERENCES markets (uuid);

-- +goose Down
ALTER TABLE orders DROP COLUMN market_uuid;
DROP TABLE markets;