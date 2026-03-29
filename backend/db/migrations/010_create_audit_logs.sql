-- +goose Up
CREATE TABLE audit_logs (
    id          INTEGER  PRIMARY KEY AUTOINCREMENT,
    table_name  TEXT     NOT NULL,
    record_id   INTEGER  NOT NULL,
    action      TEXT     NOT NULL,
    before_json TEXT     NOT NULL,
    after_json  TEXT,
    operated_at DATETIME NOT NULL DEFAULT (datetime('now', 'localtime'))
);

CREATE INDEX idx_audit_logs_table_record ON audit_logs(table_name, record_id);

-- +goose Down
DROP INDEX idx_audit_logs_table_record;
DROP TABLE audit_logs;
