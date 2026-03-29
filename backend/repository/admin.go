package repository

import (
	"database/sql"
	"purchase-system/model"
	"time"
)

type AdminRepository struct {
	db *sql.DB
}

func NewAdminRepository(db *sql.DB) *AdminRepository {
	return &AdminRepository{db: db}
}

// GetOrCreate は管理者レコードを取得するか、存在しない場合は空のパスワードで作成する
func (r *AdminRepository) GetOrCreate() (*model.Admin, error) {
	admin := &model.Admin{}
	err := r.db.QueryRow("SELECT id, password_hash, updated_at FROM admin LIMIT 1").
		Scan(&admin.ID, &admin.PasswordHash, &admin.UpdatedAt)

	if err == sql.ErrNoRows {
		// レコード作成
		now := time.Now()
		result, err := r.db.Exec(
			"INSERT INTO admin (password_hash, updated_at) VALUES (?, ?)",
			"", now,
		)
		if err != nil {
			return nil, err
		}

		id, err := result.LastInsertId()
		if err != nil {
			return nil, err
		}

		return &model.Admin{
			ID:        int(id),
			PasswordHash: "",
			UpdatedAt: now,
		}, nil
	}
	if err != nil {
		return nil, err
	}

	return admin, nil
}

// Update はパスワードハッシュを更新する
func (r *AdminRepository) Update(passwordHash string) error {
	now := time.Now()
	_, err := r.db.Exec(
		"UPDATE admin SET password_hash = ?, updated_at = ? WHERE id = 1",
		passwordHash, now,
	)
	return err
}
