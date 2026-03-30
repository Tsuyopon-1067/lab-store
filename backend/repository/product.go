package repository

import (
	"database/sql"
	"purchase-system/model"
	"time"
)

type ProductRepository struct {
	db *sql.DB
}

func NewProductRepository(db *sql.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) GetByID(id int) (*model.Product, error) {
	product := &model.Product{}
	err := r.db.QueryRow(
		"SELECT id, name, barcode, is_active, note, created_at FROM products WHERE id = ?",
		id,
	).Scan(&product.ID, &product.Name, &product.Barcode, &product.IsActive, &product.Note, &product.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return product, nil
}

func (r *ProductRepository) GetByBarcode(barcode string) (*model.ProductWithPrice, error) {
	product := &model.ProductWithPrice{}
	err := r.db.QueryRow(`
		SELECT p.id, p.name, p.barcode, p.is_active, p.note, pp.price, p.created_at
		FROM products p
		JOIN product_prices pp ON pp.product_id = p.id
		WHERE p.barcode = ? AND p.is_active = 1 AND pp.valid_to IS NULL
	`, barcode).Scan(&product.ID, &product.Name, &product.Barcode, &product.IsActive, &product.Note, &product.Price, &product.CreatedAt)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return product, nil
}

func (r *ProductRepository) ListAll() ([]*model.ProductWithPrice, error) {
	rows, err := r.db.Query(`
		SELECT p.id, p.name, p.barcode, p.is_active, p.note, COALESCE(pp.price, 0), p.created_at
		FROM products p
		LEFT JOIN product_prices pp ON pp.product_id = p.id AND pp.valid_to IS NULL
		ORDER BY p.created_at DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*model.ProductWithPrice
	for rows.Next() {
		product := &model.ProductWithPrice{}
		if err := rows.Scan(&product.ID, &product.Name, &product.Barcode, &product.IsActive, &product.Note, &product.Price, &product.CreatedAt); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, rows.Err()
}

func (r *ProductRepository) Search(query string) ([]*model.ProductWithPrice, error) {
	rows, err := r.db.Query(`
		SELECT p.id, p.name, p.barcode, p.is_active, p.note, COALESCE(pp.price, 0), p.created_at
		FROM products p
		LEFT JOIN product_prices pp ON pp.product_id = p.id AND pp.valid_to IS NULL
		WHERE p.name LIKE ? OR p.barcode LIKE ?
		ORDER BY p.created_at DESC
	`, "%"+query+"%", "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var products []*model.ProductWithPrice
	for rows.Next() {
		product := &model.ProductWithPrice{}
		if err := rows.Scan(&product.ID, &product.Name, &product.Barcode, &product.IsActive, &product.Note, &product.Price, &product.CreatedAt); err != nil {
			return nil, err
		}
		products = append(products, product)
	}
	return products, rows.Err()
}

func (r *ProductRepository) Create(name, barcode, note string, price int) (*model.ProductWithPrice, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	result, err := tx.Exec(
		"INSERT INTO products (name, barcode, is_active, note) VALUES (?, ?, 1, ?)",
		name, barcode, note,
	)
	if err != nil {
		return nil, err
	}

	productID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	_, err = tx.Exec(
		"INSERT INTO product_prices (product_id, price) VALUES (?, ?)",
		productID, price,
	)
	if err != nil {
		return nil, err
	}

	if err := tx.Commit(); err != nil {
		return nil, err
	}

	product := &model.ProductWithPrice{
		ID:        int(productID),
		Name:      name,
		Barcode:   barcode,
		IsActive:  1,
		Note:      &note,
		Price:     price,
		CreatedAt: time.Now(),
	}
	return product, nil
}

func (r *ProductRepository) Update(id int, name, barcode, note string) (*model.ProductWithPrice, error) {
	_, err := r.db.Exec(
		"UPDATE products SET name = ?, barcode = ?, note = ? WHERE id = ?",
		name, barcode, note, id,
	)
	if err != nil {
		return nil, err
	}

	product := &model.Product{}
	err = r.db.QueryRow(
		"SELECT id, name, barcode, is_active, note, created_at FROM products WHERE id = ?",
		id,
	).Scan(&product.ID, &product.Name, &product.Barcode, &product.IsActive, &product.Note, &product.CreatedAt)
	if err != nil {
		return nil, err
	}

	var price int
	err = r.db.QueryRow(
		"SELECT price FROM product_prices WHERE product_id = ? AND valid_to IS NULL",
		id,
	).Scan(&price)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	return &model.ProductWithPrice{
		ID:        product.ID,
		Name:      product.Name,
		Barcode:   product.Barcode,
		IsActive:  product.IsActive,
		Note:      product.Note,
		Price:     price,
		CreatedAt: product.CreatedAt,
	}, nil
}

func (r *ProductRepository) SetActive(id int, isActive int) error {
	_, err := r.db.Exec(
		"UPDATE products SET is_active = ? WHERE id = ?",
		isActive, id,
	)
	return err
}

func (r *ProductRepository) ChangePrice(id int, newPrice int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	now := time.Now()
	_, err = tx.Exec(
		"UPDATE product_prices SET valid_to = ? WHERE product_id = ? AND valid_to IS NULL",
		now, id,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		"INSERT INTO product_prices (product_id, price, valid_from) VALUES (?, ?, ?)",
		id, newPrice, now,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *ProductRepository) GetPriceHistory(id int) ([]*model.ProductPrice, error) {
	rows, err := r.db.Query(
		"SELECT id, product_id, price, valid_from, valid_to FROM product_prices WHERE product_id = ? ORDER BY valid_from DESC",
		id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var prices []*model.ProductPrice
	for rows.Next() {
		price := &model.ProductPrice{}
		if err := rows.Scan(&price.ID, &price.ProductID, &price.Price, &price.ValidFrom, &price.ValidTo); err != nil {
			return nil, err
		}
		prices = append(prices, price)
	}
	return prices, rows.Err()
}
