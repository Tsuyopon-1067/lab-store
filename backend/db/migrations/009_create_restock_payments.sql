-- +goose Up
CREATE TABLE restock_payments (
    id          INTEGER  PRIMARY KEY AUTOINCREMENT,
    user_id     INTEGER  NOT NULL REFERENCES users(id),
    amount      INTEGER  NOT NULL,
    settled_at  DATETIME NOT NULL,
    note        TEXT,
    created_by  TEXT     NOT NULL DEFAULT 'admin',
    updated_by  TEXT,
    deleted_at  DATETIME,
    created_at  DATETIME NOT NULL DEFAULT (datetime('now', 'localtime')),
    updated_at  DATETIME
);

CREATE INDEX idx_restock_payments_user_id ON restock_payments(user_id);

-- +goose Down
DROP INDEX idx_restock_payments_user_id;
DROP TABLE restock_payments;
