package middleware

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"purchase-system/repository"
)

func APIKeyAuth(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		key := c.GetHeader("X-API-Key")
		if key == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Missing X-API-Key header"})
			c.Abort()
			return
		}

		apiKeyRepo := repository.NewAPIKeyRepository(db)
		apiKey, err := apiKeyRepo.GetByRawKey(key)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			c.Abort()
			return
		}

		if apiKey == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid API key"})
			c.Abort()
			return
		}

		c.Set("api_key", apiKey)
		c.Next()
	}
}
