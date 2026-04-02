package model

import "time"

type APIKey struct {
	ID        int        `json:"id"`
	Name      string     `json:"name"`
	KeyPrefix string     `json:"key_prefix"`
	IsActive  int        `json:"is_active"`
	CreatedAt time.Time  `json:"created_at"`
	LastUsedAt *time.Time `json:"last_used_at"`
}

type CreateAPIKeyRequest struct {
	Name string `json:"name" binding:"required"`
}

type CreateAPIKeyResponse struct {
	ID         int        `json:"id"`
	Name       string     `json:"name"`
	KeyPrefix  string     `json:"key_prefix"`
	RawKey     string     `json:"key"`
	IsActive   int        `json:"is_active"`
	CreatedAt  time.Time  `json:"created_at"`
	LastUsedAt *time.Time `json:"last_used_at"`
}

type ListAPIKeysResponse struct {
	Keys []APIKey `json:"keys"`
}

type UserBalanceResponse struct {
	Barcode          string `json:"barcode"`
	UserName         string `json:"user_name"`
	PurchaseTotal    int    `json:"purchase_total"`
	PurchasePaid     int    `json:"purchase_paid"`
	PurchaseUnpaid   int    `json:"purchase_unpaid"`
	RestockTotal     int    `json:"restock_total"`
	RestockSettled   int    `json:"restock_settled"`
	RestockUnclaimed int    `json:"restock_unclaimed"`
	NetBalance       int    `json:"net_balance"`
}
