-- +goose Up
CREATE TABLE admin_sessions (
    id         INTEGER  PRIMARY KEY AUTOINCREMENT,
    token      TEXT     NOT NULL UNIQUE,
    created_at DATETIME NOT NULL DEFAULT (datetime('now', 'localtime')),
    expires_at DATETIME NOT NULL
);

CREATE INDEX idx_admin_sessions_token ON admin_sessions(token);

-- +goose Down
DROP INDEX idx_admin_sessions_token;
DROP TABLE admin_sessions;
