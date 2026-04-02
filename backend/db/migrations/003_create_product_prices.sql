-- +goose Up
CREATE TABLE product_prices (
    id         INTEGER  PRIMARY KEY AUTOINCREMENT,
    product_id INTEGER  NOT NULL REFERENCES products(id),
    price      INTEGER  NOT NULL,
    valid_from DATETIME NOT NULL DEFAULT (datetime('now', 'localtime')),
    valid_to   DATETIME
);

CREATE INDEX idx_product_prices_product_id ON product_prices(product_id);

-- +goose Down
DROP INDEX idx_product_prices_product_id;
DROP TABLE product_prices;
