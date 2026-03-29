package repository

import (
	"database/sql"
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
		"SELECT id, user_id, purchased_at FROM purchases WHERE id = ?",
		id,
	).Scan(&purchase.ID, &purchase.UserID, &purchase.PurchasedAt)
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
		"SELECT id, user_id, purchased_at FROM purchases WHERE user_id = ? ORDER BY purchased_at DESC",
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var purchases []*model.Purchase
	for rows.Next() {
		purchase := &model.Purchase{}
		if err := rows.Scan(&purchase.ID, &purchase.UserID, &purchase.PurchasedAt); err != nil {
			return nil, err
		}
		purchases = append(purchases, purchase)
	}
	return purchases, rows.Err()
}

func (r *PurchaseRepository) ListAll() ([]*model.Purchase, error) {
	rows, err := r.db.Query(
		"SELECT id, user_id, purchased_at FROM purchases ORDER BY purchased_at DESC",
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var purchases []*model.Purchase
	for rows.Next() {
		purchase := &model.Purchase{}
		if err := rows.Scan(&purchase.ID, &purchase.UserID, &purchase.PurchasedAt); err != nil {
			return nil, err
		}
		purchases = append(purchases, purchase)
	}
	return purchases, rows.Err()
}
