-- +goose Up
CREATE TABLE admin (
    id            INTEGER  PRIMARY KEY AUTOINCREMENT,
    password_hash TEXT     NOT NULL,
    updated_at    DATETIME NOT NULL DEFAULT (datetime('now', 'localtime'))
);

-- +goose Down
DROP TABLE admin;
