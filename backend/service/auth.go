package service

import (
	"database/sql"
	"purchase-system/model"
	"purchase-system/repository"

	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	db *sql.DB
}

func NewAuthService(db *sql.DB) *AuthService {
	return &AuthService{db: db}
}

// HashPassword はパスワードをbcryptでハッシュ化する
func (s *AuthService) HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

// VerifyPassword はパスワードが正しいか確認する
func (s *AuthService) VerifyPassword(hash, password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// SetPassword は管理者パスワードを設定する
func (s *AuthService) SetPassword(password string) error {
	hash, err := s.HashPassword(password)
	if err != nil {
		return err
	}

	adminRepo := repository.NewAdminRepository(s.db)
	return adminRepo.Update(hash)
}

// Login はパスワードを確認してセッションを作成する
func (s *AuthService) Login(password string, sessionTimeoutMinutes int) (*model.AdminSession, error) {
	adminRepo := repository.NewAdminRepository(s.db)
	admin, err := adminRepo.GetOrCreate()
	if err != nil {
		return nil, err
	}

	// パスワード未設定の場合はエラー
	if admin.PasswordHash == "" {
		return nil, ErrPasswordNotSet
	}

	// パスワード確認
	if !s.VerifyPassword(admin.PasswordHash, password) {
		return nil, ErrInvalidPassword
	}

	// セッション作成
	sessionRepo := repository.NewSessionRepository(s.db)
	session, err := sessionRepo.Create(sessionTimeoutMinutes)
	if err != nil {
		return nil, err
	}

	return session, nil
}

// Logout はセッションを削除する
func (s *AuthService) Logout(token string) error {
	sessionRepo := repository.NewSessionRepository(s.db)
	return sessionRepo.Delete(token)
}

// VerifySession はセッショントークンが有効か確認する
func (s *AuthService) VerifySession(token string) (*model.AdminSession, error) {
	sessionRepo := repository.NewSessionRepository(s.db)
	session, err := sessionRepo.GetByToken(token)
	if err != nil {
		return nil, err
	}
	if session == nil {
		return nil, ErrSessionNotFound
	}
	return session, nil
}
