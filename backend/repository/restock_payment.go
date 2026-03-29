package repository

import (
	"database/sql"
	"encoding/json"
	"purchase-system/model"
	"time"
)

type RestockPaymentRepository struct {
	db *sql.DB
}

func NewRestockPaymentRepository(db *sql.DB) *RestockPaymentRepository {
	return &RestockPaymentRepository{db: db}
}

func (r *RestockPaymentRepository) Create(userID int, amount int, settledAt time.Time, note string, createdBy string) (*model.RestockPayment, error) {
	result, err := r.db.Exec(
		"INSERT INTO restock_payments (user_id, amount, settled_at, note, created_by) VALUES (?, ?, ?, ?, ?)",
		userID, amount, settledAt, note, createdBy,
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

func (r *RestockPaymentRepository) GetByID(id int) (*model.RestockPayment, error) {
	restockPayment := &model.RestockPayment{}
	err := r.db.QueryRow(`
		SELECT id, user_id, amount, settled_at, note, created_by, updated_by, deleted_at, created_at, updated_at
		FROM restock_payments
		WHERE id = ?
	`, id).Scan(
		&restockPayment.ID, &restockPayment.UserID, &restockPayment.Amount, &restockPayment.SettledAt, &restockPayment.Note,
		&restockPayment.CreatedBy, &restockPayment.UpdatedBy, &restockPayment.DeletedAt, &restockPayment.CreatedAt, &restockPayment.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return restockPayment, nil
}

func (r *RestockPaymentRepository) ListAll() ([]*model.RestockPayment, error) {
	rows, err := r.db.Query(`
		SELECT id, user_id, amount, settled_at, note, created_by, updated_by, deleted_at, created_at, updated_at
		FROM restock_payments
		WHERE deleted_at IS NULL
		ORDER BY settled_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []*model.RestockPayment
	for rows.Next() {
		payment := &model.RestockPayment{}
		if err := rows.Scan(
			&payment.ID, &payment.UserID, &payment.Amount, &payment.SettledAt, &payment.Note,
			&payment.CreatedBy, &payment.UpdatedBy, &payment.DeletedAt, &payment.CreatedAt, &payment.UpdatedAt,
		); err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}
	return payments, rows.Err()
}

func (r *RestockPaymentRepository) ListByUserID(userID int) ([]*model.RestockPayment, error) {
	rows, err := r.db.Query(`
		SELECT id, user_id, amount, settled_at, note, created_by, updated_by, deleted_at, created_at, updated_at
		FROM restock_payments
		WHERE user_id = ? AND deleted_at IS NULL
		ORDER BY settled_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []*model.RestockPayment
	for rows.Next() {
		payment := &model.RestockPayment{}
		if err := rows.Scan(
			&payment.ID, &payment.UserID, &payment.Amount, &payment.SettledAt, &payment.Note,
			&payment.CreatedBy, &payment.UpdatedBy, &payment.DeletedAt, &payment.CreatedAt, &payment.UpdatedAt,
		); err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}
	return payments, rows.Err()
}

func (r *RestockPaymentRepository) Update(id int, amount int, settledAt time.Time, note string, updatedBy string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 元のレコード取得（監査ログ用）
	oldPayment := &model.RestockPayment{}
	err = tx.QueryRow(`
		SELECT id, user_id, amount, settled_at, note, created_by, updated_by, deleted_at, created_at, updated_at
		FROM restock_payments
		WHERE id = ?
	`, id).Scan(
		&oldPayment.ID, &oldPayment.UserID, &oldPayment.Amount, &oldPayment.SettledAt, &oldPayment.Note,
		&oldPayment.CreatedBy, &oldPayment.UpdatedBy, &oldPayment.DeletedAt, &oldPayment.CreatedAt, &oldPayment.UpdatedAt,
	)
	if err != nil {
		return err
	}

	oldJSON, _ := json.Marshal(oldPayment)

	// 更新
	now := time.Now()
	_, err = tx.Exec(
		"UPDATE restock_payments SET amount = ?, settled_at = ?, note = ?, updated_by = ?, updated_at = ? WHERE id = ?",
		amount, settledAt, note, updatedBy, now, id,
	)
	if err != nil {
		return err
	}

	// 監査ログ記録
	newPayment := &model.RestockPayment{
		ID:        id,
		UserID:    oldPayment.UserID,
		Amount:    amount,
		SettledAt: settledAt,
		Note:      &note,
		CreatedBy: oldPayment.CreatedBy,
		UpdatedBy: &updatedBy,
		CreatedAt: oldPayment.CreatedAt,
		UpdatedAt: &now,
	}
	newJSON, _ := json.Marshal(newPayment)

	_, err = tx.Exec(
		"INSERT INTO audit_logs (table_name, record_id, action, before_json, after_json) VALUES (?, ?, ?, ?, ?)",
		"restock_payments", id, "update", string(oldJSON), string(newJSON),
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *RestockPaymentRepository) Delete(id int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 元のレコード取得
	oldPayment := &model.RestockPayment{}
	err = tx.QueryRow(`
		SELECT id, user_id, amount, settled_at, note, created_by, updated_by, deleted_at, created_at, updated_at
		FROM restock_payments
		WHERE id = ?
	`, id).Scan(
		&oldPayment.ID, &oldPayment.UserID, &oldPayment.Amount, &oldPayment.SettledAt, &oldPayment.Note,
		&oldPayment.CreatedBy, &oldPayment.UpdatedBy, &oldPayment.DeletedAt, &oldPayment.CreatedAt, &oldPayment.UpdatedAt,
	)
	if err != nil {
		return err
	}

	oldJSON, _ := json.Marshal(oldPayment)

	// 論理削除
	now := time.Now()
	_, err = tx.Exec(
		"UPDATE restock_payments SET deleted_at = ?, updated_at = ? WHERE id = ?",
		now, now, id,
	)
	if err != nil {
		return err
	}

	// 監査ログ記録（削除時は after_json は NULL）
	_, err = tx.Exec(
		"INSERT INTO audit_logs (table_name, record_id, action, before_json, after_json) VALUES (?, ?, ?, ?, NULL)",
		"restock_payments", id, "delete", string(oldJSON),
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}
