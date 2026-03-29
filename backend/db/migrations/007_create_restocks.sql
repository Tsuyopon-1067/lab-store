-- +goose Up
CREATE TABLE restocks (
    id            INTEGER  PRIMARY KEY AUTOINCREMENT,
    user_id       INTEGER  NOT NULL REFERENCES users(id),
    total_amount  INTEGER  NOT NULL,
    restocked_at  DATETIME NOT NULL,
    note          TEXT,
    deleted_at    DATETIME,
    created_at    DATETIME NOT NULL DEFAULT (datetime('now', 'localtime')),
    updated_at    DATETIME
);

CREATE INDEX idx_restocks_user_id      ON restocks(user_id);
CREATE INDEX idx_restocks_restocked_at ON restocks(restocked_at);

-- +goose Down
DROP INDEX idx_restocks_restocked_at;
DROP INDEX idx_restocks_user_id;
DROP TABLE restocks;
