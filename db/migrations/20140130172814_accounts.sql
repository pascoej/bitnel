-- +goose Up
CREATE TABLE accounts (
    uuid uuid NOT NULL PRIMARY KEY DEFAULT uuid_generate_v4(),
<<<<<<< HEAD
    user_uuid uuid NOT NULL,
    type int 
=======
    user_uuid uuid NOT NULL
>>>>>>> 8822ead2d45d8caa6d290ab78fc0e24a8ef488d4
);


-- +goose Down
DROP TABLE accounts;