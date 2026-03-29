package handler

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"purchase-system/config"
	"purchase-system/model"
	"purchase-system/middleware"
	"purchase-system/service"
)

func Login(db *sql.DB, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// ログイン処理
		authService := service.NewAuthService(db)
		session, err := authService.Login(req.Password, cfg.Session.TimeoutMinutes)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid password"})
			return
		}

		// セッションCookieを設定
		c.SetCookie(
			middleware.SessionCookieName,
			session.Token,
			cfg.Session.TimeoutMinutes*60,
			"/",
			"localhost",
			false,
			true,
		)

		response := model.LoginResponse{
			Token:   session.Token,
			Message: "Login successful",
		}
		c.JSON(http.StatusOK, response)
	}
}

func Logout(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		// Cookieからトークンを取得
		token, err := c.Cookie(middleware.SessionCookieName)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "No session to logout"})
			return
		}

		// セッション削除
		authService := service.NewAuthService(db)
		if err := authService.Logout(token); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout"})
			return
		}

		// Cookieをクリア
		c.SetCookie(
			middleware.SessionCookieName,
			"",
			-1,
			"/",
			"localhost",
			false,
			true,
		)

		c.JSON(http.StatusOK, gin.H{"message": "Logout successful"})
	}
}
