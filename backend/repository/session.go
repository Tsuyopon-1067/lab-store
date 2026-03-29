package repository

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"purchase-system/model"
	"time"
)

type SessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) *SessionRepository {
	return &SessionRepository{db: db}
}

// GenerateToken はランダムなセッショントークンを生成する
func (r *SessionRepository) GenerateToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// Create はセッションを作成する
func (r *SessionRepository) Create(timeoutMinutes int) (*model.AdminSession, error) {
	token, err := r.GenerateToken()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	expiresAt := now.Add(time.Duration(timeoutMinutes) * time.Minute)

	result, err := r.db.Exec(
		"INSERT INTO admin_sessions (token, created_at, expires_at) VALUES (?, ?, ?)",
		token, now, expiresAt,
	)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &model.AdminSession{
		ID:        int(id),
		Token:     token,
		CreatedAt: now,
		ExpiresAt: expiresAt,
	}, nil
}

// GetByToken はトークンからセッションを取得する
func (r *SessionRepository) GetByToken(token string) (*model.AdminSession, error) {
	session := &model.AdminSession{}
	err := r.db.QueryRow(
		"SELECT id, token, created_at, expires_at FROM admin_sessions WHERE token = ?",
		token,
	).Scan(&session.ID, &session.Token, &session.CreatedAt, &session.ExpiresAt)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	// セッション期限切れ確認
	if time.Now().After(session.ExpiresAt) {
		return nil, nil
	}

	return session, nil
}

// Delete はトークンを削除する
func (r *SessionRepository) Delete(token string) error {
	_, err := r.db.Exec("DELETE FROM admin_sessions WHERE token = ?", token)
	return err
}

// DeleteExpired は期限切れセッションを削除する
func (r *SessionRepository) DeleteExpired() error {
	_, err := r.db.Exec(
		"DELETE FROM admin_sessions WHERE expires_at < ?",
		time.Now(),
	)
	return err
}
