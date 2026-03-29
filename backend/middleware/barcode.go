package middleware

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"purchase-system/repository"
)

const UserContextKey = "user"

// BarcodeAuth はX-User-Barcodeヘッダでユーザーを特定し、contextに格納するミドルウェア
func BarcodeAuth(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		barcode := c.GetHeader("X-User-Barcode")
		if barcode == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "X-User-Barcode header required"})
			c.Abort()
			return
		}

		repo := repository.NewUserRepository(db)
		user, err := repo.GetByBarcode(barcode)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			c.Abort()
			return
		}
		if user == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
			c.Abort()
			return
		}

		// ユーザーが有効か確認
		if user.IsActive != 1 {
			c.JSON(http.StatusForbidden, gin.H{"error": "User is inactive"})
			c.Abort()
			return
		}

		c.Set(UserContextKey, user)
		c.Next()
	}
}
