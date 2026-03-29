package model

import "time"

type Payment struct {
	ID        int        `json:"id"`
	UserID    int        `json:"user_id"`
	Amount    int        `json:"amount"`
	PaidAt    time.Time  `json:"paid_at"`
	Note      *string    `json:"note"`
	CreatedBy string     `json:"created_by"`
	UpdatedBy *string    `json:"updated_by"`
	DeletedAt *time.Time `json:"deleted_at"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}

type PaymentWithUserName struct {
	Payment  *Payment
	UserName string
}

type CreatePaymentRequest struct {
	UserID int       `json:"user_id" binding:"required"`
	Amount int       `json:"amount" binding:"required,gt=0"`
	PaidAt time.Time `json:"paid_at" binding:"required"`
	Note   string    `json:"note"`
}

type UpdatePaymentRequest struct {
	Amount int       `json:"amount" binding:"required,gt=0"`
	PaidAt time.Time `json:"paid_at" binding:"required"`
	Note   string    `json:"note"`
}

type PaymentResponse struct {
	ID        int        `json:"id"`
	UserID    int        `json:"user_id"`
	UserName  string     `json:"user_name"`
	Amount    int        `json:"amount"`
	PaidAt    time.Time  `json:"paid_at"`
	Note      *string    `json:"note"`
	CreatedBy string     `json:"created_by"`
	UpdatedBy *string    `json:"updated_by"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
