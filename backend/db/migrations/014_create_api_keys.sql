-- +goose Up
CREATE TABLE api_keys (
    id           INTEGER  PRIMARY KEY AUTOINCREMENT,
    name         TEXT     NOT NULL,
    key_hash     TEXT     NOT NULL UNIQUE,
    key_prefix   TEXT     NOT NULL,
    is_active    INTEGER  NOT NULL DEFAULT 1,
    created_at   DATETIME NOT NULL DEFAULT (datetime('now', 'localtime')),
    last_used_at DATETIME
);
-- +goose Down
DROP TABLE api_keys;
