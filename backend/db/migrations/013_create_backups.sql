-- +goose Up
CREATE TABLE backups (
    id         INTEGER  PRIMARY KEY AUTOINCREMENT,
    filename   TEXT     NOT NULL,
    created_at DATETIME NOT NULL DEFAULT (datetime('now', 'localtime'))
);

-- +goose Down
DROP TABLE backups;
