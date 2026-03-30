package repository

import (
	"database/sql"
	"purchase-system/model"
)

type BackupRepository struct {
	db *sql.DB
}

func NewBackupRepository(db *sql.DB) *BackupRepository {
	return &BackupRepository{db: db}
}

func (r *BackupRepository) Create(filename string) (*model.Backup, error) {
	result, err := r.db.Exec(
		"INSERT INTO backups (filename) VALUES (?)",
		filename,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return r.GetByID(int(id))
}

func (r *BackupRepository) GetByID(id int) (*model.Backup, error) {
	backup := &model.Backup{}
	err := r.db.QueryRow(`
		SELECT id, filename, created_at
		FROM backups
		WHERE id = ?
	`, id).Scan(&backup.ID, &backup.Filename, &backup.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return backup, nil
}

func (r *BackupRepository) ListAll() ([]*model.Backup, error) {
	rows, err := r.db.Query(`
		SELECT id, filename, created_at
		FROM backups
		ORDER BY created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var backups []*model.Backup
	for rows.Next() {
		backup := &model.Backup{}
		if err := rows.Scan(&backup.ID, &backup.Filename, &backup.CreatedAt); err != nil {
			return nil, err
		}
		backups = append(backups, backup)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return backups, nil
}
