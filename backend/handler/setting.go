package handler

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"purchase-system/model"
	"purchase-system/repository"
)

func GetSettings(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		repo := repository.NewSettingRepository(db)
		setting, err := repo.Get()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get settings"})
			return
		}
		c.JSON(http.StatusOK, setting)
	}
}

func UpdateSettings(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.UpdateSettingRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Enforce safe minimums
		if req.BackupIntervalMinutes < 10 {
			req.BackupIntervalMinutes = 10
		}
		if req.SessionTimeoutMinutes < 1 {
			req.SessionTimeoutMinutes = 1
		}

		repo := repository.NewSettingRepository(db)
		setting, err := repo.Update(&req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update settings"})
			return
		}
		c.JSON(http.StatusOK, setting)
	}
}
