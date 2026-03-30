package repository

import (
	"crypto/rand"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"purchase-system/model"
)

type APIKeyRepository struct {
	db *sql.DB
}

func NewAPIKeyRepository(db *sql.DB) *APIKeyRepository {
	return &APIKeyRepository{db: db}
}

// GenerateKey generates a random 64-char hex string (32 bytes)
func (r *APIKeyRepository) GenerateKey() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("failed to generate random key: %w", err)
	}
	return hex.EncodeToString(b), nil
}

// hashKey returns sha256 hex of the raw key
func (r *APIKeyRepository) hashKey(rawKey string) string {
	hash := sha256.Sum256([]byte(rawKey))
	return hex.EncodeToString(hash[:])
}

// Create generates a new API key, hashes it, and stores it
func (r *APIKeyRepository) Create(name string) (*model.APIKey, string, error) {
	rawKey, err := r.GenerateKey()
	if err != nil {
		return nil, "", err
	}

	keyHash := r.hashKey(rawKey)
	keyPrefix := rawKey[:8]

	result, err := r.db.Exec(
		"INSERT INTO api_keys (name, key_hash, key_prefix, is_active, created_at) VALUES (?, ?, ?, 1, datetime('now', 'localtime'))",
		name, keyHash, keyPrefix,
	)
	if err != nil {
		return nil, "", fmt.Errorf("failed to insert api key: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, "", fmt.Errorf("failed to get inserted id: %w", err)
	}

	apiKey, err := r.GetByID(int(id))
	if err != nil {
		return nil, "", err
	}

	return apiKey, rawKey, nil
}

// GetByID retrieves an API key by ID
func (r *APIKeyRepository) GetByID(id int) (*model.APIKey, error) {
	row := r.db.QueryRow(
		"SELECT id, name, key_prefix, is_active, created_at, last_used_at FROM api_keys WHERE id = ?",
		id,
	)

	apiKey := &model.APIKey{}
	err := row.Scan(&apiKey.ID, &apiKey.Name, &apiKey.KeyPrefix, &apiKey.IsActive, &apiKey.CreatedAt, &apiKey.LastUsedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to scan api key: %w", err)
	}

	return apiKey, nil
}

// GetByRawKey retrieves an API key by raw key, validates it's active, and updates last_used_at
func (r *APIKeyRepository) GetByRawKey(rawKey string) (*model.APIKey, error) {
	keyHash := r.hashKey(rawKey)

	row := r.db.QueryRow(
		"SELECT id, name, key_prefix, is_active, created_at, last_used_at FROM api_keys WHERE key_hash = ?",
		keyHash,
	)

	apiKey := &model.APIKey{}
	err := row.Scan(&apiKey.ID, &apiKey.Name, &apiKey.KeyPrefix, &apiKey.IsActive, &apiKey.CreatedAt, &apiKey.LastUsedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("failed to scan api key: %w", err)
	}

	if apiKey.IsActive != 1 {
		return nil, nil
	}

	// Update last_used_at
	_, err = r.db.Exec(
		"UPDATE api_keys SET last_used_at = datetime('now', 'localtime') WHERE id = ?",
		apiKey.ID,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to update last_used_at: %w", err)
	}

	return apiKey, nil
}

// ListAll retrieves all API keys
func (r *APIKeyRepository) ListAll() ([]model.APIKey, error) {
	rows, err := r.db.Query(
		"SELECT id, name, key_prefix, is_active, created_at, last_used_at FROM api_keys ORDER BY created_at DESC",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to query api keys: %w", err)
	}
	defer rows.Close()

	var apiKeys []model.APIKey
	for rows.Next() {
		apiKey := model.APIKey{}
		err := rows.Scan(&apiKey.ID, &apiKey.Name, &apiKey.KeyPrefix, &apiKey.IsActive, &apiKey.CreatedAt, &apiKey.LastUsedAt)
		if err != nil {
			return nil, fmt.Errorf("failed to scan api key: %w", err)
		}
		apiKeys = append(apiKeys, apiKey)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating api keys: %w", err)
	}

	return apiKeys, nil
}

// Delete removes an API key by ID (physical delete)
func (r *APIKeyRepository) Delete(id int) error {
	result, err := r.db.Exec("DELETE FROM api_keys WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete api key: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("api key not found")
	}

	return nil
}
