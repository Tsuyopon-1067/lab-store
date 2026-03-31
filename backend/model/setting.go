package model

import "time"

// Setting represents the single-row settings table (id always = 1).
type Setting struct {
	ID                    int       `json:"id"`
	SessionTimeoutMinutes int       `json:"session_timeout_minutes"`
	BackupPath            string    `json:"backup_path"`
	BackupIntervalMinutes int       `json:"backup_interval_minutes"`
	UpdatedAt             time.Time `json:"updated_at"`
}

// UpdateSettingRequest is the JSON body for PUT /api/settings.
// All fields are required; the client always sends the full object.
type UpdateSettingRequest struct {
	SessionTimeoutMinutes int    `json:"session_timeout_minutes" binding:"required,min=1"`
	BackupPath            string `json:"backup_path"             binding:"required"`
	BackupIntervalMinutes int    `json:"backup_interval_minutes" binding:"required,min=10"`
}
