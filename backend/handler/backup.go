package handler

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"purchase-system/config"
	"purchase-system/model"
	"purchase-system/service"
)

func CreateBackup(db *sql.DB, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		svc := service.NewBackupService(db, cfg)
		backup, err := svc.CreateBackup()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create backup"})
			return
		}

		c.JSON(http.StatusCreated, backup)
	}
}

func ListBackups(db *sql.DB, cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		svc := service.NewBackupService(db, cfg)
		backups, err := svc.ListBackups()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list backups"})
			return
		}

		if backups == nil {
			backups = []*model.Backup{}
		}

		c.JSON(http.StatusOK, backups)
	}
}
