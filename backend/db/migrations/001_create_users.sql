-- +goose Up
CREATE TABLE users (
    id         INTEGER  PRIMARY KEY AUTOINCREMENT,
    name       TEXT     NOT NULL,
    barcode    TEXT     NOT NULL UNIQUE,
    is_active  INTEGER  NOT NULL DEFAULT 1,
    created_at DATETIME NOT NULL DEFAULT (datetime('now', 'localtime'))
);

-- +goose Down
DROP TABLE users;
