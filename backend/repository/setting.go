package repository

import (
	"database/sql"
	"purchase-system/model"
	"time"
)

type SettingRepository struct {
	db *sql.DB
}

func NewSettingRepository(db *sql.DB) *SettingRepository {
	return &SettingRepository{db: db}
}

// Get retrieves the single settings row (id = 1).
func (r *SettingRepository) Get() (*model.Setting, error) {
	s := &model.Setting{}
	err := r.db.QueryRow(`
		SELECT id, session_timeout_minutes, backup_path, backup_interval_minutes, updated_at
		FROM settings
		WHERE id = 1
	`).Scan(&s.ID, &s.SessionTimeoutMinutes, &s.BackupPath, &s.BackupIntervalMinutes, &s.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return s, nil
}

// Update overwrites all mutable columns on the single row.
func (r *SettingRepository) Update(req *model.UpdateSettingRequest) (*model.Setting, error) {
	now := time.Now()
	_, err := r.db.Exec(`
		UPDATE settings
		SET session_timeout_minutes = ?,
		    backup_path             = ?,
		    backup_interval_minutes = ?,
		    updated_at              = ?
		WHERE id = 1
	`, req.SessionTimeoutMinutes, req.BackupPath, req.BackupIntervalMinutes, now)
	if err != nil {
		return nil, err
	}
	return r.Get()
}
