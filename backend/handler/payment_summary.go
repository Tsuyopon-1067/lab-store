package handler

import (
	"bytes"
	"database/sql"
	"encoding/csv"
	"fmt"
	"net/http"
	"purchase-system/repository"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// GetUserPaymentSummary returns JSON with all active users' payment summaries
func GetUserPaymentSummary(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		repo := repository.NewSummaryRepository(db)
		balances, err := repo.GetAllUsersBalance()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user payment summary"})
			return
		}

		c.JSON(http.StatusOK, balances)
	}
}

// ExportUserPaymentSummaryCSV exports user payment summary as CSV
func ExportUserPaymentSummaryCSV(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		repo := repository.NewSummaryRepository(db)
		balances, err := repo.GetAllUsersBalance()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch user payment summary"})
			return
		}

		var buf bytes.Buffer
		writer := csv.NewWriter(&buf)

		// Write header with Japanese labels
		writer.Write([]string{"ユーザー名", "利用額", "仕入れ金額", "最終的な支払い金額"})

		for _, balance := range balances {
			writer.Write([]string{
				balance.UserName,
				strconv.Itoa(balance.PurchaseUnpaid),
				strconv.Itoa(balance.RestockUnclaimed),
				strconv.Itoa(balance.NetBalance),
			})
		}

		writer.Flush()
		if err := writer.Error(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "CSV generation error"})
			return
		}

		filename := fmt.Sprintf("user_payment_summary_%s.csv",
			time.Now().Format("2006-01-02"))
		c.Header("Content-Type", "text/csv; charset=utf-8")
		c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filename))
		c.Data(http.StatusOK, "text/csv", buf.Bytes())
	}
}
