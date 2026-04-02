package repository

import (
	"database/sql"
	"encoding/json"
	"purchase-system/model"
	"time"
)

type PaymentRepository struct {
	db *sql.DB
}

func NewPaymentRepository(db *sql.DB) *PaymentRepository {
	return &PaymentRepository{db: db}
}

func (r *PaymentRepository) Create(userID int, amount int, paidAt time.Time, note string, createdBy string) (*model.Payment, error) {
	result, err := r.db.Exec(
		"INSERT INTO payments (user_id, amount, paid_at, note, created_by) VALUES (?, ?, ?, ?, ?)",
		userID, amount, paidAt, note, createdBy,
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

func (r *PaymentRepository) GetByID(id int) (*model.Payment, error) {
	payment := &model.Payment{}
	err := r.db.QueryRow(`
		SELECT id, user_id, amount, paid_at, note, created_by, updated_by, deleted_at, created_at, updated_at
		FROM payments
		WHERE id = ?
	`, id).Scan(
		&payment.ID, &payment.UserID, &payment.Amount, &payment.PaidAt, &payment.Note,
		&payment.CreatedBy, &payment.UpdatedBy, &payment.DeletedAt, &payment.CreatedAt, &payment.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return payment, nil
}

func (r *PaymentRepository) ListAll() ([]*model.Payment, error) {
	rows, err := r.db.Query(`
		SELECT id, user_id, amount, paid_at, note, created_by, updated_by, deleted_at, created_at, updated_at
		FROM payments
		WHERE deleted_at IS NULL
		ORDER BY paid_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []*model.Payment
	for rows.Next() {
		payment := &model.Payment{}
		if err := rows.Scan(
			&payment.ID, &payment.UserID, &payment.Amount, &payment.PaidAt, &payment.Note,
			&payment.CreatedBy, &payment.UpdatedBy, &payment.DeletedAt, &payment.CreatedAt, &payment.UpdatedAt,
		); err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}
	return payments, rows.Err()
}

func (r *PaymentRepository) ListByUserID(userID int) ([]*model.Payment, error) {
	rows, err := r.db.Query(`
		SELECT id, user_id, amount, paid_at, note, created_by, updated_by, deleted_at, created_at, updated_at
		FROM payments
		WHERE user_id = ? AND deleted_at IS NULL
		ORDER BY paid_at DESC
	`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []*model.Payment
	for rows.Next() {
		payment := &model.Payment{}
		if err := rows.Scan(
			&payment.ID, &payment.UserID, &payment.Amount, &payment.PaidAt, &payment.Note,
			&payment.CreatedBy, &payment.UpdatedBy, &payment.DeletedAt, &payment.CreatedAt, &payment.UpdatedAt,
		); err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}
	return payments, rows.Err()
}

func (r *PaymentRepository) Update(id int, amount int, paidAt time.Time, note string, updatedBy string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 元のレコード取得（監査ログ用）
	oldPayment := &model.Payment{}
	err = tx.QueryRow(`
		SELECT id, user_id, amount, paid_at, note, created_by, updated_by, deleted_at, created_at, updated_at
		FROM payments
		WHERE id = ?
	`, id).Scan(
		&oldPayment.ID, &oldPayment.UserID, &oldPayment.Amount, &oldPayment.PaidAt, &oldPayment.Note,
		&oldPayment.CreatedBy, &oldPayment.UpdatedBy, &oldPayment.DeletedAt, &oldPayment.CreatedAt, &oldPayment.UpdatedAt,
	)
	if err != nil {
		return err
	}

	oldJSON, _ := json.Marshal(oldPayment)

	// 更新
	now := time.Now()
	_, err = tx.Exec(
		"UPDATE payments SET amount = ?, paid_at = ?, note = ?, updated_by = ?, updated_at = ? WHERE id = ?",
		amount, paidAt, note, updatedBy, now, id,
	)
	if err != nil {
		return err
	}

	// 監査ログ記録
	newPayment := &model.Payment{
		ID:        id,
		UserID:    oldPayment.UserID,
		Amount:    amount,
		PaidAt:    paidAt,
		Note:      &note,
		CreatedBy: oldPayment.CreatedBy,
		UpdatedBy: &updatedBy,
		CreatedAt: oldPayment.CreatedAt,
		UpdatedAt: &now,
	}
	newJSON, _ := json.Marshal(newPayment)

	_, err = tx.Exec(
		"INSERT INTO audit_logs (table_name, record_id, action, before_json, after_json) VALUES (?, ?, ?, ?, ?)",
		"payments", id, "update", string(oldJSON), string(newJSON),
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *PaymentRepository) Delete(id int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 元のレコード取得
	oldPayment := &model.Payment{}
	err = tx.QueryRow(`
		SELECT id, user_id, amount, paid_at, note, created_by, updated_by, deleted_at, created_at, updated_at
		FROM payments
		WHERE id = ?
	`, id).Scan(
		&oldPayment.ID, &oldPayment.UserID, &oldPayment.Amount, &oldPayment.PaidAt, &oldPayment.Note,
		&oldPayment.CreatedBy, &oldPayment.UpdatedBy, &oldPayment.DeletedAt, &oldPayment.CreatedAt, &oldPayment.UpdatedAt,
	)
	if err != nil {
		return err
	}

	oldJSON, _ := json.Marshal(oldPayment)

	// 論理削除
	now := time.Now()
	_, err = tx.Exec(
		"UPDATE payments SET deleted_at = ?, updated_at = ? WHERE id = ?",
		now, now, id,
	)
	if err != nil {
		return err
	}

	// 監査ログ記録（削除時は after_json は NULL）
	_, err = tx.Exec(
		"INSERT INTO audit_logs (table_name, record_id, action, before_json, after_json) VALUES (?, ?, ?, ?, NULL)",
		"payments", id, "delete", string(oldJSON),
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}
