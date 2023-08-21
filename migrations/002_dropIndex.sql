-- +goose Up

DROP INDEX spends_user_id_uindex ON spends;

-- +goose Down

CREATE UNIQUE INDEX spends_user_id_uindex ON spends (user_id);

