package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"purchase-system/model"
	"purchase-system/repository"
)

func CreateRestock(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.CreateRestockRequest
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

		// 商品情報を検証
		productRepo := repository.NewProductRepository(db)
		var restockItems []struct {
			ProductID int
			Quantity  int
			UnitPrice int
		}

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

			restockItems = append(restockItems, struct {
				ProductID int
				Quantity  int
				UnitPrice int
			}{
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				UnitPrice: item.UnitPrice,
			})
		}

		// 仕入れを作成
		restockRepo := repository.NewRestockRepository(db)
		restock, items, err := restockRepo.Create(user.ID, req.TotalAmount, req.Note, restockItems)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create restock"})
			return
		}

		// 詳細情報を構築
		var restockDetails []*model.RestockItemDetail
		for _, item := range items {
			product, _ := productRepo.GetByID(item.ProductID)
			restockDetails = append(restockDetails, &model.RestockItemDetail{
				Item:        item,
				ProductName: product.Name,
				ProductID:   product.ID,
				Subtotal:    item.Quantity * item.UnitPrice,
			})
		}

		response := model.RestockResponse{
			RestockID:   restock.ID,
			UserID:      user.ID,
			UserName:    user.Name,
			Items:       restockDetails,
			TotalAmount: req.TotalAmount,
			RestockedAt: restock.RestockedAt,
			Note:        restock.Note,
		}

		c.JSON(http.StatusCreated, response)
	}
}

func ListRestocks(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Query("user_id")
		limit := c.DefaultQuery("limit", "100")
		offset := c.DefaultQuery("offset", "0")

		limitInt, _ := strconv.Atoi(limit)
		offsetInt, _ := strconv.Atoi(offset)

		query := "SELECT id, user_id, total_amount, restocked_at, note, deleted_at, created_at, updated_at FROM restocks WHERE deleted_at IS NULL ORDER BY restocked_at DESC"
		args := []interface{}{}

		if userID != "" {
			query = "SELECT id, user_id, total_amount, restocked_at, note, deleted_at, created_at, updated_at FROM restocks WHERE user_id = ? AND deleted_at IS NULL ORDER BY restocked_at DESC"
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

		var restocks []*model.Restock
		for rows.Next() {
			restock := &model.Restock{}
			if err := rows.Scan(&restock.ID, &restock.UserID, &restock.TotalAmount, &restock.RestockedAt, &restock.Note, &restock.DeletedAt, &restock.CreatedAt, &restock.UpdatedAt); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
				return
			}
			restocks = append(restocks, restock)
		}

		c.JSON(http.StatusOK, restocks)
	}
}

func UpdateRestock(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid restock ID"})
			return
		}

		var req model.UpdateRestockRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// 商品情報を検証
		productRepo := repository.NewProductRepository(db)
		var restockItems []struct {
			ProductID int
			Quantity  int
			UnitPrice int
		}

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

			restockItems = append(restockItems, struct {
				ProductID int
				Quantity  int
				UnitPrice int
			}{
				ProductID: item.ProductID,
				Quantity:  item.Quantity,
				UnitPrice: item.UnitPrice,
			})
		}

		// 仕入れを更新
		restockRepo := repository.NewRestockRepository(db)
		if err := restockRepo.Update(id, req.TotalAmount, req.Note, restockItems); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update restock"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Restock updated"})
	}
}

func DeleteRestock(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid restock ID"})
			return
		}

		restockRepo := repository.NewRestockRepository(db)
		if err := restockRepo.Delete(id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete restock"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Restock deleted"})
	}
}
