package model

import "time"

type RestockPayment struct {
	ID        int        `json:"id"`
	UserID    int        `json:"user_id"`
	Amount    int        `json:"amount"`
	SettledAt time.Time  `json:"settled_at"`
	Note      *string    `json:"note"`
	CreatedBy string     `json:"created_by"`
	UpdatedBy *string    `json:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type RestockPaymentWithUserName struct {
	RestockPayment *RestockPayment
	UserName       string
}

type CreateRestockPaymentRequest struct {
	UserID    int       `json:"user_id" binding:"required"`
	Amount    int       `json:"amount" binding:"required,gt=0"`
	SettledAt time.Time `json:"settled_at" binding:"required"`
	Note      string    `json:"note"`
}

type UpdateRestockPaymentRequest struct {
	Amount    int       `json:"amount" binding:"required,gt=0"`
	SettledAt time.Time `json:"settled_at" binding:"required"`
	Note      string    `json:"note"`
}

type RestockPaymentResponse struct {
	ID        int        `json:"id"`
	UserID    int        `json:"user_id"`
	UserName  string     `json:"user_name"`
	Amount    int        `json:"amount"`
	SettledAt time.Time  `json:"settled_at"`
	Note      *string    `json:"note"`
	CreatedBy string     `json:"created_by"`
	UpdatedBy *string    `json:"updated_by"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
