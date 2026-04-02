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

		// 購入履歴を取得
		var histories []*model.PurchaseHistory
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
					ID:          item.ID,
					ProductID:   item.ProductID,
					Quantity:    item.Quantity,
					UnitPrice:   item.UnitPrice,
					ProductName: product.Name,
					Subtotal:    subtotal,
				})
			}

			histories = append(histories, &model.PurchaseHistory{
				ID:          purchase.ID,
				PurchasedAt: purchase.PurchasedAt,
				TotalAmount: totalAmount,
				Items:       itemDetails,
			})
		}

		c.JSON(http.StatusOK, histories)
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
		restockRepo := repository.NewRestockRepository(db)
		restocks, err := restockRepo.GetByUserID(currentUser.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		// 仕入れ履歴を取得
		var histories []*model.RestockHistory
		productRepo := repository.NewProductRepository(db)
		for _, restock := range restocks {
			items, err := restockRepo.GetItemsByRestockID(restock.ID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
				return
			}

			var itemDetails []*model.RestockItemDetail
			for _, item := range items {
				product, _ := productRepo.GetByID(item.ProductID)
				itemDetails = append(itemDetails, &model.RestockItemDetail{
					ID:          item.ID,
					ProductID:   item.ProductID,
					Quantity:    item.Quantity,
					UnitPrice:   item.UnitPrice,
					ProductName: product.Name,
					Subtotal:    item.Quantity * item.UnitPrice,
				})
			}

			histories = append(histories, &model.RestockHistory{
				ID:          restock.ID,
				TotalAmount: restock.TotalAmount,
				RestockedAt: restock.RestockedAt,
				Note:        restock.Note,
				Items:       itemDetails,
			})
		}

		c.JSON(http.StatusOK, histories)
	}
}
