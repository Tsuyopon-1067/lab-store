package repository

import (
	"database/sql"
	"purchase-system/model"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByID(id int) (*model.User, error) {
	user := &model.User{}
	err := r.db.QueryRow(
		"SELECT id, name, barcode, is_active, created_at FROM users WHERE id = ?",
		id,
	).Scan(&user.ID, &user.Name, &user.Barcode, &user.IsActive, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetByBarcode(barcode string) (*model.User, error) {
	user := &model.User{}
	err := r.db.QueryRow(
		"SELECT id, name, barcode, is_active, created_at FROM users WHERE barcode = ?",
		barcode,
	).Scan(&user.ID, &user.Name, &user.Barcode, &user.IsActive, &user.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) ListAll() ([]*model.User, error) {
	rows, err := r.db.Query("SELECT id, name, barcode, is_active, created_at FROM users ORDER BY created_at DESC")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		user := &model.User{}
		if err := rows.Scan(&user.ID, &user.Name, &user.Barcode, &user.IsActive, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, rows.Err()
}

func (r *UserRepository) Create(name, barcode string) (*model.User, error) {
	result, err := r.db.Exec(
		"INSERT INTO users (name, barcode, is_active) VALUES (?, ?, 1)",
		name, barcode,
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

func (r *UserRepository) Update(id int, name, barcode string) (*model.User, error) {
	_, err := r.db.Exec(
		"UPDATE users SET name = ?, barcode = ? WHERE id = ?",
		name, barcode, id,
	)
	if err != nil {
		return nil, err
	}

	return r.GetByID(id)
}

func (r *UserRepository) SetActive(id int, isActive int) error {
	_, err := r.db.Exec(
		"UPDATE users SET is_active = ? WHERE id = ?",
		isActive, id,
	)
	return err
}
