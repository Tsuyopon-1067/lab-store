package model

import "time"

type User struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Barcode   string    `json:"barcode"`
	IsActive  int       `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}
