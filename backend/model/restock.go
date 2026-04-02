package model

import "time"

type Restock struct {
	ID          int       `json:"id"`
	UserID      int       `json:"user_id"`
	TotalAmount int       `json:"total_amount"`
	RestockedAt time.Time `json:"restocked_at"`
	Note        *string   `json:"note"`
	DeletedAt   *time.Time `json:"deleted_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   *time.Time `json:"updated_at"`
}

type RestockItem struct {
	ID        int `json:"id"`
	RestockID int `json:"restock_id"`
	ProductID int `json:"product_id"`
	Quantity  int `json:"quantity"`
	UnitPrice int `json:"unit_price"`
}

type RestockDetail struct {
	Restock     *Restock
	Items       []*RestockItemDetail
	TotalAmount int
}

type RestockHistory struct {
	ID          int                   `json:"id"`
	TotalAmount int                   `json:"total_amount"`
	RestockedAt time.Time             `json:"restocked_at"`
	Note        *string               `json:"note"`
	Items       []*RestockItemDetail  `json:"items"`
}

type RestockItemDetail struct {
	ID        int    `json:"id"`
	ProductID int    `json:"product_id"`
	Quantity  int    `json:"quantity"`
	UnitPrice int    `json:"unit_price"`
	Subtotal  int    `json:"subtotal"`
	ProductName string `json:"product_name"`
}

type CreateRestockRequest struct {
	UserBarcode string `json:"user_barcode" binding:"required"`
	Items       []struct {
		ProductID int `json:"product_id" binding:"required"`
		Quantity  int `json:"quantity" binding:"required,min=1"`
		UnitPrice int `json:"unit_price" binding:"required"`
		Subtotal  int `json:"subtotal"`
	} `json:"items" binding:"required,min=1"`
	TotalAmount int    `json:"total_amount" binding:"required"`
	Note        string `json:"note"`
}

type UpdateRestockRequest struct {
	Items       []struct {
		ProductID int `json:"product_id" binding:"required"`
		Quantity  int `json:"quantity" binding:"required,min=1"`
		UnitPrice int `json:"unit_price" binding:"required"`
		Subtotal  int `json:"subtotal"`
	} `json:"items"`
	TotalAmount int    `json:"total_amount" binding:"required"`
	Note        string `json:"note"`
}

type RestockResponse struct {
	RestockID   int                   `json:"restock_id"`
	UserID      int                   `json:"user_id"`
	UserName    string                `json:"user_name"`
	Items       []*RestockItemDetail  `json:"items"`
	TotalAmount int                   `json:"total_amount"`
	RestockedAt time.Time             `json:"restocked_at"`
	Note        *string               `json:"note"`
}
