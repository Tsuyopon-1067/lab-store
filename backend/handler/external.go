package handler

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"purchase-system/model"
	"purchase-system/repository"
)

func GetUserBalance(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		barcode := c.Param("barcode")
		if barcode == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Barcode is required"})
			return
		}

		userRepo := repository.NewUserRepository(db)
		user, err := userRepo.GetByBarcode(barcode)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get user"})
			return
		}

		if user == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		// Get balance using SummaryRepository
		summaryRepo := repository.NewSummaryRepository(db)
		now := time.Now()
		firstDay := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
		lastDay := firstDay.AddDate(0, 1, -1).Add(time.Hour*23 + time.Minute*59 + time.Second*59)

		report, err := summaryRepo.GetSummary(firstDay, lastDay)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get balance"})
			return
		}

		// Find the specific user's summary
		var userSummary *model.SummaryItem
		if report != nil && report.Users != nil {
			for _, s := range report.Users {
				if s.UserID == user.ID {
					userSummary = s
					break
				}
			}
		}

		if userSummary == nil {
			// User has no transactions this month
			c.JSON(http.StatusOK, model.UserBalanceResponse{
				Barcode:          user.Barcode,
				UserName:         user.Name,
				PurchaseTotal:    0,
				PurchasePaid:     0,
				PurchaseUnpaid:   0,
				RestockTotal:     0,
				RestockSettled:   0,
				RestockUnclaimed: 0,
				NetBalance:       0,
			})
			return
		}

		c.JSON(http.StatusOK, model.UserBalanceResponse{
			Barcode:          user.Barcode,
			UserName:         user.Name,
			PurchaseTotal:    userSummary.PurchaseTotal,
			PurchasePaid:     userSummary.PurchasePaid,
			PurchaseUnpaid:   userSummary.PurchaseUnpaid,
			RestockTotal:     userSummary.RestockTotal,
			RestockSettled:   userSummary.RestockSettled,
			RestockUnclaimed: userSummary.RestockUnclaimed,
			NetBalance:       userSummary.NetBalance,
		})
	}
}

func GetMonthlyReport(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		fromStr := c.DefaultQuery("from", "")
		toStr := c.DefaultQuery("to", "")

		var from, to time.Time

		// Parse from/to dates (YYYY-MM-DD format)
		if fromStr == "" || toStr == "" {
			// Default to current month
			now := time.Now()
			from = time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())
			to = from.AddDate(0, 1, -1).Add(time.Hour*23 + time.Minute*59 + time.Second*59)
		} else {
			var err error
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

			to = to.Add(time.Hour*23 + time.Minute*59 + time.Second*59)
		}

		summaryRepo := repository.NewSummaryRepository(db)
		report, err := summaryRepo.GetSummary(from, to)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get summary"})
			return
		}

		if report != nil {
			report.Period.From = from.Format("2006-01-02")
			report.Period.To = to.Format("2006-01-02")
		}

		c.JSON(http.StatusOK, report)
	}
}

func ListProductsExternal(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		productRepo := repository.NewProductRepository(db)
		products, err := productRepo.ListAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to list products"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"products": products,
		})
	}
}
