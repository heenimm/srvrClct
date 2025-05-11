-- +goose Up
CREATE TABLE IF NOT EXISTS users (
                                     id INTEGER PRIMARY KEY AUTOINCREMENT,
                                     login TEXT NOT NULL UNIQUE,
                                     password TEXT NOT NULL
);

-- +goose Down
DROP TABLE IF EXISTS users;
