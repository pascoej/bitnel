-- +goose Up
CREATE EXTENSION "uuid-ossp";

CREATE TABLE users (
    uuid uuid NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    email text NOT NULL,
    password_hash text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT NOW()
);

CREATE TYPE order_side AS ENUM('bid', 'ask');

CREATE TABLE orders (
    uuid uuid NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
    size bigint NOT NULL,
    initial_size bigint NOT NULL,
    price bigint,
    side smallint NOT NULL,
    status smallint NOT NULL,
    created_at timestamptz NOT NULL DEFAULT NOW()
);

-- +goose Down
DROP TABLE users;
DROP TABLE orders;
DROP EXTENSION "uuid-ossp";