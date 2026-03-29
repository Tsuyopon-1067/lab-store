package handler

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"purchase-system/middleware"
	"purchase-system/model"
	"purchase-system/repository"
	"purchase-system/service"
)

func GetMyBalance(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get(middleware.UserContextKey)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		currentUser := user.(*model.User)
		summaryService := service.NewSummaryService(db)
		balance, err := summaryService.GetUserBalance(currentUser.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get balance"})
			return
		}

		c.JSON(http.StatusOK, balance)
	}
}

func GetMyPurchases(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get(middleware.UserContextKey)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		currentUser := user.(*model.User)
		purchaseRepo := repository.NewPurchaseRepository(db)
		purchases, err := purchaseRepo.GetByUserID(currentUser.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		// 購入明細を取得
		var details []*model.PurchaseDetail
		productRepo := repository.NewProductRepository(db)
		for _, purchase := range purchases {
			items, err := purchaseRepo.GetItemsByPurchaseID(purchase.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
				return
			}

			var itemDetails []*model.PurchaseItemDetail
			totalAmount := 0
			for _, item := range items {
				product, _ := productRepo.GetByID(item.ProductID)
				subtotal := item.Quantity * item.UnitPrice
				totalAmount += subtotal
				itemDetails = append(itemDetails, &model.PurchaseItemDetail{
					Item:        item,
					ProductName: product.Name,
					ProductID:   product.ID,
					Subtotal:    subtotal,
				})
			}

			details = append(details, &model.PurchaseDetail{
				Purchase:    purchase,
				Items:       itemDetails,
				TotalAmount: totalAmount,
			})
		}

		c.JSON(http.StatusOK, details)
	}
}

func GetMyRestocks(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, exists := c.Get(middleware.UserContextKey)
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
			return
		}

		currentUser := user.(*model.User)

		rows, err := db.Query(`
			SELECT id, user_id, total_amount, restocked_at, note
			FROM restocks
			WHERE user_id = ? AND deleted_at IS NULL
			ORDER BY restocked_at DESC
		`, currentUser.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		defer rows.Close()

		var restocks []map[string]interface{}
		for rows.Next() {
			var id, userID, totalAmount int
			var note sql.NullString
			var restockedAt interface{}

			if err := rows.Scan(&id, &userID, &totalAmount, &restockedAt, &note); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
				return
			}

			restocks = append(restocks, map[string]interface{}{
				"id":             id,
				"user_id":        userID,
				"total_amount":   totalAmount,
				"restocked_at":   restockedAt,
				"note":           note.String,
			})
		}

		c.JSON(http.StatusOK, restocks)
	}
}
