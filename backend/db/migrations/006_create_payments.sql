-- +goose Up
CREATE TABLE payments (
    id          INTEGER  PRIMARY KEY AUTOINCREMENT,
    user_id     INTEGER  NOT NULL REFERENCES users(id),
    amount      INTEGER  NOT NULL,
    paid_at     DATETIME NOT NULL,
    note        TEXT,
    created_by  TEXT     NOT NULL DEFAULT 'admin',
    updated_by  TEXT,
    deleted_at  DATETIME,
    created_at  DATETIME NOT NULL DEFAULT (datetime('now', 'localtime')),
    updated_at  DATETIME
);

CREATE INDEX idx_payments_user_id ON payments(user_id);

-- +goose Down
DROP INDEX idx_payments_user_id;
DROP TABLE payments;
