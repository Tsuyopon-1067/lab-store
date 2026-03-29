-- +goose Up
CREATE TABLE restock_items (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    restock_id  INTEGER NOT NULL REFERENCES restocks(id),
    product_id  INTEGER NOT NULL REFERENCES products(id),
    quantity    INTEGER NOT NULL DEFAULT 1,
    unit_price  INTEGER NOT NULL
);

CREATE INDEX idx_restock_items_restock_id ON restock_items(restock_id);

-- +goose Down
DROP INDEX idx_restock_items_restock_id;
DROP TABLE restock_items;
