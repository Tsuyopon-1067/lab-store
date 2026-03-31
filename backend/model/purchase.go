package model

import "time"

type Purchase struct {
	ID          int        `json:"id"`
	UserID      int        `json:"user_id"`
	PurchasedAt time.Time  `json:"purchased_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at,omitempty"`
}

type PurchaseItem struct {
	ID         int        `json:"id"`
	PurchaseID int        `json:"purchase_id"`
	ProductID  int        `json:"product_id"`
	Quantity   int        `json:"quantity"`
	UnitPrice  int        `json:"unit_price"`
	DeletedAt  *time.Time `json:"deleted_at,omitempty"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  *time.Time `json:"updated_at,omitempty"`
}

type PurchaseDetail struct {
	Purchase      *Purchase
	Items         []*PurchaseItemDetail
	TotalAmount   int
}

type PurchaseHistory struct {
	ID          int                    `json:"id"`
	PurchasedAt time.Time              `json:"purchased_at"`
	TotalAmount int                    `json:"total_amount"`
	Items       []*PurchaseItemDetail  `json:"items"`
}

type PurchaseItemDetail struct {
	ID          int    `json:"id"`
	ProductID   int    `json:"product_id"`
	Quantity    int    `json:"quantity"`
	UnitPrice   int    `json:"unit_price"`
	Subtotal    int    `json:"subtotal"`
	ProductName string `json:"product_name"`
}

type BalanceSummary struct {
	UserID               int    `json:"user_id"`
	UserName             string `json:"user_name"`
	PurchaseUnpaid       int    `json:"purchase_unpaid"`
	RestockUnclaimed     int    `json:"restock_unclaimed"`
	NetBalance           int    `json:"net_balance"`
}

type CreatePurchaseRequest struct {
	UserBarcode string `json:"user_barcode" binding:"required"`
	Items       []struct {
		ProductID int `json:"product_id" binding:"required"`
		Quantity  int `json:"quantity" binding:"required,min=1"`
	} `json:"items" binding:"required,min=1"`
}

type PurchaseResponse struct {
	PurchaseID     int                   `json:"purchase_id"`
	UserID         int                   `json:"user_id"`
	UserName       string                `json:"user_name"`
	Items          []*PurchaseItemDetail `json:"items"`
	TotalAmount    int                   `json:"total_amount"`
	PurchasedAt    time.Time             `json:"purchased_at"`
	UpdatedBalance *BalanceSummary       `json:"updated_balance"`
}
