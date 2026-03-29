-- +goose Up
CREATE TABLE purchase_items (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    purchase_id INTEGER NOT NULL REFERENCES purchases(id),
    product_id  INTEGER NOT NULL REFERENCES products(id),
    quantity    INTEGER NOT NULL DEFAULT 1,
    unit_price  INTEGER NOT NULL
);

CREATE INDEX idx_purchase_items_purchase_id ON purchase_items(purchase_id);

-- +goose Down
DROP INDEX idx_purchase_items_purchase_id;
DROP TABLE purchase_items;
