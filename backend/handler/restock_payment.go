package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"purchase-system/model"
	"purchase-system/repository"
)

func CreateRestockPayment(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.CreateRestockPaymentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// ユーザー存在確認（仕入れ者）
		userRepo := repository.NewUserRepository(db)
		user, err := userRepo.GetByID(req.UserID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		if user == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}

		// 精算記録作成
		restockPaymentRepo := repository.NewRestockPaymentRepository(db)
		restockPayment, err := restockPaymentRepo.Create(req.UserID, req.Amount, req.SettledAt, req.Note, "admin")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create restock payment"})
			return
		}

		response := model.RestockPaymentResponse{
			ID:        restockPayment.ID,
			UserID:    restockPayment.UserID,
			UserName:  user.Name,
			Amount:    restockPayment.Amount,
			SettledAt: restockPayment.SettledAt,
			Note:      restockPayment.Note,
			CreatedBy: restockPayment.CreatedBy,
			UpdatedBy: restockPayment.UpdatedBy,
			CreatedAt: restockPayment.CreatedAt,
			UpdatedAt: restockPayment.UpdatedAt,
		}

		c.JSON(http.StatusCreated, response)
	}
}

func ListRestockPayments(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Query("user_id")
		limit := c.DefaultQuery("limit", "100")
		offset := c.DefaultQuery("offset", "0")

		limitInt, _ := strconv.Atoi(limit)
		offsetInt, _ := strconv.Atoi(offset)

		var payments []*model.RestockPayment
		var err error

		if userID != "" {
			restockPaymentRepo := repository.NewRestockPaymentRepository(db)
			userIDInt, _ := strconv.Atoi(userID)
			payments, err = restockPaymentRepo.ListByUserID(userIDInt)
		} else {
			restockPaymentRepo := repository.NewRestockPaymentRepository(db)
			payments, err = restockPaymentRepo.ListAll()
		}

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		// ページネーション適用
		if offsetInt < len(payments) {
			end := offsetInt + limitInt
			if end > len(payments) {
				end = len(payments)
			}
			payments = payments[offsetInt:end]
		} else {
			payments = []*model.RestockPayment{}
		}

		c.JSON(http.StatusOK, payments)
	}
}

func UpdateRestockPayment(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid restock payment ID"})
			return
		}

		var req model.UpdateRestockPaymentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		restockPaymentRepo := repository.NewRestockPaymentRepository(db)
		if err := restockPaymentRepo.Update(id, req.Amount, req.SettledAt, req.Note, "admin"); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update restock payment"})
			return
		}

		// 更新後のレコード取得
		restockPayment, err := restockPaymentRepo.GetByID(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		// ユーザー名を取得
		userRepo := repository.NewUserRepository(db)
		user, _ := userRepo.GetByID(restockPayment.UserID)

		response := model.RestockPaymentResponse{
			ID:        restockPayment.ID,
			UserID:    restockPayment.UserID,
			UserName:  user.Name,
			Amount:    restockPayment.Amount,
			SettledAt: restockPayment.SettledAt,
			Note:      restockPayment.Note,
			CreatedBy: restockPayment.CreatedBy,
			UpdatedBy: restockPayment.UpdatedBy,
			CreatedAt: restockPayment.CreatedAt,
			UpdatedAt: restockPayment.UpdatedAt,
		}

		c.JSON(http.StatusOK, response)
	}
}

func DeleteRestockPayment(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid restock payment ID"})
			return
		}

		restockPaymentRepo := repository.NewRestockPaymentRepository(db)
		if err := restockPaymentRepo.Delete(id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete restock payment"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Restock payment deleted"})
	}
}
