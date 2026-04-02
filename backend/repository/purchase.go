package repository

import (
	"database/sql"
	"encoding/json"
	"purchase-system/model"
	"time"
)

type PurchaseRepository struct {
	db *sql.DB
}

func NewPurchaseRepository(db *sql.DB) *PurchaseRepository {
	return &PurchaseRepository{db: db}
}

func (r *PurchaseRepository) Create(userID int, items []struct {
	ProductID int
	Price     int
	Quantity  int
}) (*model.Purchase, []*model.PurchaseItem, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, nil, err
	}
	defer tx.Rollback()

	now := time.Now()
	result, err := tx.Exec(
		"INSERT INTO purchases (user_id, purchased_at) VALUES (?, ?)",
		userID, now,
	)
	if err != nil {
		return nil, nil, err
	}

	purchaseID, err := result.LastInsertId()
	if err != nil {
		return nil, nil, err
	}

	var purchaseItems []*model.PurchaseItem
	for _, item := range items {
		res, err := tx.Exec(
			"INSERT INTO purchase_items (purchase_id, product_id, quantity, unit_price) VALUES (?, ?, ?, ?)",
			purchaseID, item.ProductID, item.Quantity, item.Price,
		)
		if err != nil {
			return nil, nil, err
		}

		itemID, err := res.LastInsertId()
		if err != nil {
			return nil, nil, err
		}

		// Decrement product stock quantity
		_, err = tx.Exec(
			"UPDATE products SET stock_quantity = stock_quantity - ? WHERE id = ?",
			item.Quantity, item.ProductID,
		)
		if err != nil {
			return nil, nil, err
		}

		purchaseItems = append(purchaseItems, &model.PurchaseItem{
			ID:         int(itemID),
			PurchaseID: int(purchaseID),
			ProductID:  item.ProductID,
			Quantity:   item.Quantity,
			UnitPrice:  item.Price,
		})
	}

	if err := tx.Commit(); err != nil {
		return nil, nil, err
	}

	return &model.Purchase{
		ID:          int(purchaseID),
		UserID:      userID,
		PurchasedAt: now,
	}, purchaseItems, nil
}

func (r *PurchaseRepository) GetByID(id int) (*model.Purchase, error) {
	purchase := &model.Purchase{}
	err := r.db.QueryRow(
		"SELECT id, user_id, purchased_at, deleted_at, created_at, updated_at FROM purchases WHERE id = ? AND deleted_at IS NULL",
		id,
	).Scan(&purchase.ID, &purchase.UserID, &purchase.PurchasedAt, &purchase.DeletedAt, &purchase.CreatedAt, &purchase.UpdatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return purchase, nil
}

func (r *PurchaseRepository) GetItemsByPurchaseID(purchaseID int) ([]*model.PurchaseItem, error) {
	rows, err := r.db.Query(
		"SELECT id, purchase_id, product_id, quantity, unit_price FROM purchase_items WHERE purchase_id = ?",
		purchaseID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []*model.PurchaseItem
	for rows.Next() {
		item := &model.PurchaseItem{}
		if err := rows.Scan(&item.ID, &item.PurchaseID, &item.ProductID, &item.Quantity, &item.UnitPrice); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *PurchaseRepository) GetByUserID(userID int) ([]*model.Purchase, error) {
	rows, err := r.db.Query(
		"SELECT id, user_id, purchased_at, deleted_at, created_at, updated_at FROM purchases WHERE user_id = ? AND deleted_at IS NULL ORDER BY purchased_at DESC",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var purchases []*model.Purchase
	for rows.Next() {
		purchase := &model.Purchase{}
		if err := rows.Scan(&purchase.ID, &purchase.UserID, &purchase.PurchasedAt, &purchase.DeletedAt, &purchase.CreatedAt, &purchase.UpdatedAt); err != nil {
			return nil, err
		}
		purchases = append(purchases, purchase)
	}
	return purchases, rows.Err()
}

func (r *PurchaseRepository) ListAll() ([]*model.Purchase, error) {
	rows, err := r.db.Query(
		"SELECT id, user_id, purchased_at, deleted_at, created_at, updated_at FROM purchases WHERE deleted_at IS NULL ORDER BY purchased_at DESC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var purchases []*model.Purchase
	for rows.Next() {
		purchase := &model.Purchase{}
		if err := rows.Scan(&purchase.ID, &purchase.UserID, &purchase.PurchasedAt, &purchase.DeletedAt, &purchase.CreatedAt, &purchase.UpdatedAt); err != nil {
			return nil, err
		}
		purchases = append(purchases, purchase)
	}
	return purchases, rows.Err()
}

func (r *PurchaseRepository) Delete(id int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Get the old record
	oldPurchase := &model.Purchase{}
	err = tx.QueryRow(
		"SELECT id, user_id, purchased_at, deleted_at, created_at, updated_at FROM purchases WHERE id = ?",
		id,
	).Scan(&oldPurchase.ID, &oldPurchase.UserID, &oldPurchase.PurchasedAt, &oldPurchase.DeletedAt, &oldPurchase.CreatedAt, &oldPurchase.UpdatedAt)
	if err != nil {
		return err
	}

	// Restore stock quantities for all items in this purchase
	rows, err := tx.Query(
		"SELECT product_id, quantity FROM purchase_items WHERE purchase_id = ?",
		id,
	)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var productID, quantity int
		if err := rows.Scan(&productID, &quantity); err != nil {
			return err
		}
		_, err = tx.Exec(
			"UPDATE products SET stock_quantity = stock_quantity + ? WHERE id = ?",
			quantity, productID,
		)
		if err != nil {
			return err
		}
	}

	oldJSON, _ := json.Marshal(oldPurchase)

	now := time.Now()
	_, err = tx.Exec(
		"UPDATE purchases SET deleted_at = ?, updated_at = ? WHERE id = ?",
		now, now, id,
	)
	if err != nil {
		return err
	}

	// Record audit log (after_json is NULL for deletions)
	_, err = tx.Exec(
		"INSERT INTO audit_logs (table_name, record_id, action, before_json, after_json) VALUES (?, ?, ?, ?, NULL)",
		"purchases", id, "delete", string(oldJSON),
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}
