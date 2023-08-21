-- +goose Up

-- email ka column nhi h table m pata h tabhi first last h 
-- okay got it,
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    first_name VARCHAR(255) NOT NULL,
    last_name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE UNIQUE INDEX users_email_uindex ON users (email);

CREATE TABLE IF NOT EXISTS spends (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    amount INT NOT NULL,
    description VARCHAR(255) NOT NULL,
    tag VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- created_at m date hogi na, e.g. created_at = 'Wed, 09 Aug 2023 10:53:48 GMT'

CREATE UNIQUE INDEX spends_user_id_uindex ON spends (user_id);

-- +goose Down
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS spends;
