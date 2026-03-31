package handler

import (
	"database/sql"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"purchase-system/model"
	"purchase-system/repository"
	"purchase-system/service"
)

func CreatePurchase(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.CreatePurchaseRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// ユーザー特定
		userRepo := repository.NewUserRepository(db)
		user, err := userRepo.GetByBarcode(req.UserBarcode)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		if user == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		if user.IsActive != 1 {
			c.JSON(http.StatusForbidden, gin.H{"error": "User is inactive"})
			return
		}

		// 商品情報と価格を取得
		productRepo := repository.NewProductRepository(db)
		var purchaseItems []struct {
			ProductID int
			Price     int
			Quantity  int
		}
		totalAmount := 0

		for _, item := range req.Items {
			product, err := productRepo.GetByID(item.ProductID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
				return
			}
			if product == nil {
				c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
				return
			}

			// 現在価格を取得
			var price int
			err = db.QueryRow(
				"SELECT price FROM product_prices WHERE product_id = ? AND valid_to IS NULL",
				item.ProductID,
			).Scan(&price)
			if err == sql.ErrNoRows {
				c.JSON(http.StatusNotFound, gin.H{"error": "Product price not found"})
				return
			}
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
				return
			}

			purchaseItems = append(purchaseItems, struct {
				ProductID int
				Price     int
				Quantity  int
			}{
				ProductID: item.ProductID,
				Price:     price,
				Quantity:  item.Quantity,
			})
			totalAmount += price * item.Quantity
		}

		// 購入を作成
		purchaseRepo := repository.NewPurchaseRepository(db)
		purchase, items, err := purchaseRepo.Create(user.ID, purchaseItems)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create purchase"})
			return
		}

		// 詳細情報を構築
		summaryService := service.NewSummaryService(db)
		balance, err := summaryService.GetUserBalance(user.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to calculate balance"})
			return
		}

		var purchaseDetails []*model.PurchaseItemDetail
		for _, item := range items {
			product, _ := productRepo.GetByID(item.ProductID)
			purchaseDetails = append(purchaseDetails, &model.PurchaseItemDetail{
				ID:          item.ID,
				ProductID:   item.ProductID,
				Quantity:    item.Quantity,
				UnitPrice:   item.UnitPrice,
				ProductName: product.Name,
				Subtotal:    item.Quantity * item.UnitPrice,
			})
		}

		response := model.PurchaseResponse{
			PurchaseID:  purchase.ID,
			UserID:      user.ID,
			UserName:    user.Name,
			Items:       purchaseDetails,
			TotalAmount: totalAmount,
			PurchasedAt: purchase.PurchasedAt,
			UpdatedBalance: &model.BalanceSummary{
				UserID:           user.ID,
				UserName:         user.Name,
				PurchaseUnpaid:   balance.PurchaseUnpaid,
				RestockUnclaimed: balance.RestockUnclaimed,
				NetBalance:       balance.NetBalance,
			},
		}

		c.JSON(http.StatusCreated, response)
	}
}

func ListPurchases(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Query("user_id")
		limit := c.DefaultQuery("limit", "100")
		offset := c.DefaultQuery("offset", "0")

		limitInt, _ := strconv.Atoi(limit)
		offsetInt, _ := strconv.Atoi(offset)

		query := "SELECT p.id, p.user_id, u.name, p.purchased_at FROM purchases p JOIN users u ON p.user_id = u.id WHERE p.deleted_at IS NULL ORDER BY p.purchased_at DESC"
		args := []interface{}{}

		if userID != "" {
			query = "SELECT p.id, p.user_id, u.name, p.purchased_at FROM purchases p JOIN users u ON p.user_id = u.id WHERE p.user_id = ? AND p.deleted_at IS NULL ORDER BY p.purchased_at DESC"
			args = append(args, userID)
		}

		query += " LIMIT ? OFFSET ?"
		args = append(args, limitInt, offsetInt)

		rows, err := db.Query(query, args...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		defer rows.Close()

		var purchases []*model.AdminPurchaseWithItems
		productRepo := repository.NewProductRepository(db)
		purchaseRepo := repository.NewPurchaseRepository(db)

		for rows.Next() {
			var id, userID int
			var userName string
			var purchasedAt time.Time
			if err := rows.Scan(&id, &userID, &userName, &purchasedAt); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
				return
			}

			// Get purchase items
			items, err := purchaseRepo.GetItemsByPurchaseID(id)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
				return
			}

			var adminItems []*model.AdminPurchaseItem
			totalAmount := 0
			for _, item := range items {
				product, err := productRepo.GetByID(item.ProductID)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
					return
				}

				subtotal := item.Quantity * item.UnitPrice
				totalAmount += subtotal

				adminItems = append(adminItems, &model.AdminPurchaseItem{
					ProductID:   item.ProductID,
					ProductName: product.Name,
					Quantity:    item.Quantity,
					UnitPrice:   item.UnitPrice,
					Subtotal:    subtotal,
				})
			}

			purchases = append(purchases, &model.AdminPurchaseWithItems{
				ID:          id,
				UserID:      userID,
				UserName:    userName,
				PurchasedAt: purchasedAt,
				TotalAmount: totalAmount,
				Items:       adminItems,
			})
		}

		c.JSON(http.StatusOK, purchases)
	}
}

func DeletePurchase(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid purchase ID"})
			return
		}

		repo := repository.NewPurchaseRepository(db)
		if err := repo.Delete(id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete purchase"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Purchase deleted"})
	}
}
