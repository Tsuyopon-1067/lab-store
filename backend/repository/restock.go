package repository

import (
	"database/sql"
	"encoding/json"
	"purchase-system/model"
	"time"
)

type RestockRepository struct {
	db *sql.DB
}

func NewRestockRepository(db *sql.DB) *RestockRepository {
	return &RestockRepository{db: db}
}

func (r *RestockRepository) Create(userID int, totalAmount int, note string, items []struct {
	ProductID int
	Quantity  int
	UnitPrice int
}) (*model.Restock, []*model.RestockItem, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, nil, err
	}
	defer tx.Rollback()

	now := time.Now()
	result, err := tx.Exec(
		"INSERT INTO restocks (user_id, total_amount, restocked_at, note) VALUES (?, ?, ?, ?)",
		userID, totalAmount, now, note,
	)
	if err != nil {
		return nil, nil, err
	}

	restockID, err := result.LastInsertId()
	if err != nil {
		return nil, nil, err
	}

	var restockItems []*model.RestockItem
	for _, item := range items {
		res, err := tx.Exec(
			"INSERT INTO restock_items (restock_id, product_id, quantity, unit_price) VALUES (?, ?, ?, ?)",
			restockID, item.ProductID, item.Quantity, item.UnitPrice,
		)
		if err != nil {
			return nil, nil, err
		}

		itemID, err := res.LastInsertId()
		if err != nil {
			return nil, nil, err
		}

		restockItems = append(restockItems, &model.RestockItem{
			ID:        int(itemID),
			RestockID: int(restockID),
			ProductID: item.ProductID,
			Quantity:  item.Quantity,
			UnitPrice: item.UnitPrice,
		})
	}

	if err := tx.Commit(); err != nil {
		return nil, nil, err
	}

	return &model.Restock{
		ID:          int(restockID),
		UserID:      userID,
		TotalAmount: totalAmount,
		RestockedAt: now,
		Note:        &note,
		CreatedAt:   now,
	}, restockItems, nil
}

func (r *RestockRepository) GetByID(id int) (*model.Restock, error) {
	restock := &model.Restock{}
	err := r.db.QueryRow(
		"SELECT id, user_id, total_amount, restocked_at, note, deleted_at, created_at, updated_at FROM restocks WHERE id = ?",
		id,
	).Scan(&restock.ID, &restock.UserID, &restock.TotalAmount, &restock.RestockedAt, &restock.Note, &restock.DeletedAt, &restock.CreatedAt, &restock.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return restock, nil
}

func (r *RestockRepository) GetItemsByRestockID(restockID int) ([]*model.RestockItem, error) {
	rows, err := r.db.Query(
		"SELECT id, restock_id, product_id, quantity, unit_price FROM restock_items WHERE restock_id = ?",
		restockID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*model.RestockItem
	for rows.Next() {
		item := &model.RestockItem{}
		if err := rows.Scan(&item.ID, &item.RestockID, &item.ProductID, &item.Quantity, &item.UnitPrice); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *RestockRepository) GetByUserID(userID int) ([]*model.Restock, error) {
	rows, err := r.db.Query(
		"SELECT id, user_id, total_amount, restocked_at, note, deleted_at, created_at, updated_at FROM restocks WHERE user_id = ? AND deleted_at IS NULL ORDER BY restocked_at DESC",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var restocks []*model.Restock
	for rows.Next() {
		restock := &model.Restock{}
		if err := rows.Scan(&restock.ID, &restock.UserID, &restock.TotalAmount, &restock.RestockedAt, &restock.Note, &restock.DeletedAt, &restock.CreatedAt, &restock.UpdatedAt); err != nil {
			return nil, err
		}
		restocks = append(restocks, restock)
	}
	return restocks, rows.Err()
}

func (r *RestockRepository) ListAll() ([]*model.Restock, error) {
	rows, err := r.db.Query(
		"SELECT id, user_id, total_amount, restocked_at, note, deleted_at, created_at, updated_at FROM restocks WHERE deleted_at IS NULL ORDER BY restocked_at DESC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var restocks []*model.Restock
	for rows.Next() {
		restock := &model.Restock{}
		if err := rows.Scan(&restock.ID, &restock.UserID, &restock.TotalAmount, &restock.RestockedAt, &restock.Note, &restock.DeletedAt, &restock.CreatedAt, &restock.UpdatedAt); err != nil {
			return nil, err
		}
		restocks = append(restocks, restock)
	}
	return restocks, rows.Err()
}

func (r *RestockRepository) Update(id int, totalAmount int, note string, items []struct {
	ProductID int
	Quantity  int
	UnitPrice int
}) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 元のレコード取得
	oldRestock := &model.Restock{}
	err = tx.QueryRow(
		"SELECT id, user_id, total_amount, restocked_at, note, deleted_at, created_at, updated_at FROM restocks WHERE id = ?",
		id,
	).Scan(&oldRestock.ID, &oldRestock.UserID, &oldRestock.TotalAmount, &oldRestock.RestockedAt, &oldRestock.Note, &oldRestock.DeletedAt, &oldRestock.CreatedAt, &oldRestock.UpdatedAt)
	if err != nil {
		return err
	}

	// 監査ログ記録用のJSON
	oldJSON, _ := json.Marshal(oldRestock)

	now := time.Now()
	_, err = tx.Exec(
		"UPDATE restocks SET total_amount = ?, note = ?, updated_at = ? WHERE id = ?",
		totalAmount, note, now, id,
	)
	if err != nil {
		return err
	}

	// 既存アイテムを削除
	_, err = tx.Exec("DELETE FROM restock_items WHERE restock_id = ?", id)
	if err != nil {
		return err
	}

	// 新しいアイテムを作成
	for _, item := range items {
		_, err := tx.Exec(
			"INSERT INTO restock_items (restock_id, product_id, quantity, unit_price) VALUES (?, ?, ?, ?)",
			id, item.ProductID, item.Quantity, item.UnitPrice,
		)
		if err != nil {
			return err
		}
	}

	// 監査ログ記録
	newRestock := &model.Restock{
		ID:          id,
		UserID:      oldRestock.UserID,
		TotalAmount: totalAmount,
		RestockedAt: oldRestock.RestockedAt,
		Note:        &note,
		CreatedAt:   oldRestock.CreatedAt,
		UpdatedAt:   &now,
	}
	newJSON, _ := json.Marshal(newRestock)

	_, err = tx.Exec(
		"INSERT INTO audit_logs (table_name, record_id, action, before_json, after_json) VALUES (?, ?, ?, ?, ?)",
		"restocks", id, "update", string(oldJSON), string(newJSON),
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *RestockRepository) Delete(id int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// 元のレコード取得
	oldRestock := &model.Restock{}
	err = tx.QueryRow(
		"SELECT id, user_id, total_amount, restocked_at, note, deleted_at, created_at, updated_at FROM restocks WHERE id = ?",
		id,
	).Scan(&oldRestock.ID, &oldRestock.UserID, &oldRestock.TotalAmount, &oldRestock.RestockedAt, &oldRestock.Note, &oldRestock.DeletedAt, &oldRestock.CreatedAt, &oldRestock.UpdatedAt)
	if err != nil {
		return err
	}

	oldJSON, _ := json.Marshal(oldRestock)

	now := time.Now()
	_, err = tx.Exec(
		"UPDATE restocks SET deleted_at = ?, updated_at = ? WHERE id = ?",
		now, now, id,
	)
	if err != nil {
		return err
	}

	// 監査ログ記録（削除時は after_json は NULL）
	_, err = tx.Exec(
		"INSERT INTO audit_logs (table_name, record_id, action, before_json, after_json) VALUES (?, ?, ?, ?, NULL)",
		"restocks", id, "delete", string(oldJSON),
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}
