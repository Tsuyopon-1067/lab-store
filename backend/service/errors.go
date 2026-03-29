package service

import "errors"

var (
	ErrInvalidPassword = errors.New("invalid password")
	ErrPasswordNotSet  = errors.New("password not set")
	ErrSessionNotFound = errors.New("session not found or expired")
)
