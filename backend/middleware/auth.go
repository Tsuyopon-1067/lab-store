package middleware

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"purchase-system/service"
)

const SessionCookieName = "admin_session"

// AdminAuth はセッションCookieを検証するミドルウェア
func AdminAuth(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Cookieからトークンを取得
		token, err := c.Cookie(SessionCookieName)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Session required"})
			c.Abort()
			return
		}

		// セッション検証
		authService := service.NewAuthService(db)
		session, err := authService.VerifySession(token)
		if err != nil || session == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired session"})
			c.Abort()
			return
		}

		c.Set("session", session)
		c.Next()
	}
}
