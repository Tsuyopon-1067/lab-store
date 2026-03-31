package model

import "time"

type Product struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	Barcode       string    `json:"barcode"`
	IsActive      int       `json:"is_active"`
	Note          *string   `json:"note"`
	StockQuantity int       `json:"stock_quantity"`
	CreatedAt     time.Time `json:"created_at"`
}

type ProductPrice struct {
	ID        int       `json:"id"`
	ProductID int       `json:"product_id"`
	Price     int       `json:"price"`
	ValidFrom time.Time `json:"valid_from"`
	ValidTo   *time.Time `json:"valid_to"`
}

type ProductWithPrice struct {
	ID            int       `json:"id"`
	Name          string    `json:"name"`
	Barcode       string    `json:"barcode"`
	IsActive      int       `json:"is_active"`
	Note          *string   `json:"note"`
	Price         int       `json:"current_price"`
	StockQuantity int       `json:"stock_quantity"`
	CreatedAt     time.Time `json:"created_at"`
}
