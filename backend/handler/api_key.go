package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"purchase-system/model"
	"purchase-system/repository"
)

func CreateAPIKey(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.CreateAPIKeyRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		apiKeyRepo := repository.NewAPIKeyRepository(db)
		apiKey, rawKey, err := apiKeyRepo.Create(req.Name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create API key"})
			return
		}

		response := model.CreateAPIKeyResponse{
			ID:         apiKey.ID,
			Name:       apiKey.Name,
			KeyPrefix:  apiKey.KeyPrefix,
			RawKey:     rawKey,
			IsActive:   apiKey.IsActive,
			CreatedAt:  apiKey.CreatedAt,
			LastUsedAt: apiKey.LastUsedAt,
		}

		c.JSON(http.StatusCreated, response)
	}
}

func ListAPIKeys(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		apiKeyRepo := repository.NewAPIKeyRepository(db)
		apiKeys, err := apiKeyRepo.ListAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list API keys"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"keys": apiKeys,
		})
	}
}

func DeleteAPIKey(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
			return
		}

		apiKeyRepo := repository.NewAPIKeyRepository(db)
		if err := apiKeyRepo.Delete(id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete API key"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "API key deleted"})
	}
}
