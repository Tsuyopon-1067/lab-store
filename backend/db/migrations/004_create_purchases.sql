-- +goose Up
CREATE TABLE purchases (
    id           INTEGER  PRIMARY KEY AUTOINCREMENT,
    user_id      INTEGER  NOT NULL REFERENCES users(id),
    purchased_at DATETIME NOT NULL DEFAULT (datetime('now', 'localtime'))
);

CREATE INDEX idx_purchases_user_id      ON purchases(user_id);
CREATE INDEX idx_purchases_purchased_at ON purchases(purchased_at);

-- +goose Down
DROP INDEX idx_purchases_purchased_at;
DROP INDEX idx_purchases_user_id;
DROP TABLE purchases;
