-- +goose Up
CREATE TABLE products (
    id              INTEGER  PRIMARY KEY AUTOINCREMENT,
    name            TEXT     NOT NULL,
    barcode         TEXT     NOT NULL UNIQUE,
    is_active       INTEGER  NOT NULL DEFAULT 1,
    note            TEXT,
    stock_quantity  INTEGER  NOT NULL DEFAULT 0,
    created_at      DATETIME NOT NULL DEFAULT (datetime('now', 'localtime'))
);

-- +goose Down
DROP TABLE products;
