-- +goose Up
ALTER TABLE purchases ADD COLUMN deleted_at DATETIME;
ALTER TABLE purchases ADD COLUMN created_at DATETIME NOT NULL DEFAULT (datetime('now', 'localtime'));
ALTER TABLE purchases ADD COLUMN updated_at DATETIME;

ALTER TABLE purchase_items ADD COLUMN deleted_at DATETIME;
ALTER TABLE purchase_items ADD COLUMN created_at DATETIME NOT NULL DEFAULT (datetime('now', 'localtime'));
ALTER TABLE purchase_items ADD COLUMN updated_at DATETIME;

-- +goose Down
-- SQLite does not support DROP COLUMN easily; migration is one-way
