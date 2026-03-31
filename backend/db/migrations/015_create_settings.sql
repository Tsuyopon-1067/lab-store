-- +goose Up
CREATE TABLE settings (
    id                      INTEGER NOT NULL DEFAULT 1 CHECK (id = 1),
    session_timeout_minutes INTEGER NOT NULL DEFAULT 30,
    backup_path             TEXT    NOT NULL DEFAULT '../backup',
    backup_interval_minutes INTEGER NOT NULL DEFAULT 60,
    updated_at              DATETIME NOT NULL DEFAULT (datetime('now', 'localtime')),
    PRIMARY KEY (id)
);
INSERT INTO settings (id, session_timeout_minutes, backup_path, backup_interval_minutes)
VALUES (1, 30, '../backup', 60);

-- +goose Down
DROP TABLE settings;
