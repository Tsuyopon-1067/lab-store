package model

import "time"

type Admin struct {
	ID            int       `json:"id"`
	PasswordHash   string    `json:"-"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type AdminSession struct {
	ID        int       `json:"id"`
	Token     string    `json:"token"`
	CreatedAt time.Time `json:"created_at"`
	ExpiresAt time.Time `json:"expires_at"`
}

type LoginRequest struct {
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token   string    `json:"token"`
	Message string    `json:"message"`
}
