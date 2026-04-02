package handler

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"purchase-system/repository"
)

func GetSummary(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		fromStr := c.DefaultQuery("from", "")
		toStr := c.DefaultQuery("to", "")

		// 期間解析
		var from, to time.Time
		var err error

		if fromStr == "" || toStr == "" {
			// デフォルト：今月全体
			now := time.Now()
			from = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
			to = from.AddDate(0, 1, -1).Add(time.Hour * 23).Add(time.Minute * 59).Add(time.Second * 59)
		} else {
			from, err = time.Parse("2006-01-02", fromStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid from date format (use YYYY-MM-DD)"})
				return
			}

			to, err = time.Parse("2006-01-02", toStr)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid to date format (use YYYY-MM-DD)"})
				return
			}

			// to を終日に設定
			to = to.Add(time.Hour * 23).Add(time.Minute * 59).Add(time.Second * 59)
		}

		// サマリ取得
		summaryRepo := repository.NewSummaryRepository(db)
		report, err := summaryRepo.GetSummary(from, to)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		c.JSON(http.StatusOK, report)
	}
}
