package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"purchase-system/model"
	"purchase-system/repository"
)

func CreatePayment(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req model.CreatePaymentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// ユーザー存在確認
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

		// 支払い記録作成
		paymentRepo := repository.NewPaymentRepository(db)
		payment, err := paymentRepo.Create(req.UserID, req.Amount, req.PaidAt, req.Note, "admin")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create payment"})
			return
		}

		response := model.PaymentResponse{
			ID:        payment.ID,
			UserID:    payment.UserID,
			UserName:  user.Name,
			Amount:    payment.Amount,
			PaidAt:    payment.PaidAt,
			Note:      payment.Note,
			CreatedBy: payment.CreatedBy,
			UpdatedBy: payment.UpdatedBy,
			CreatedAt: payment.CreatedAt,
			UpdatedAt: payment.UpdatedAt,
		}

		c.JSON(http.StatusCreated, response)
	}
}

func ListPayments(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := c.Query("user_id")
		limit := c.DefaultQuery("limit", "100")
		offset := c.DefaultQuery("offset", "0")

		limitInt, _ := strconv.Atoi(limit)
		offsetInt, _ := strconv.Atoi(offset)

		var payments []*model.Payment
		var err error

		if userID != "" {
			paymentRepo := repository.NewPaymentRepository(db)
			userIDInt, _ := strconv.Atoi(userID)
			payments, err = paymentRepo.ListByUserID(userIDInt)
		} else {
			paymentRepo := repository.NewPaymentRepository(db)
			payments, err = paymentRepo.ListAll()
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
			payments = []*model.Payment{}
		}

		c.JSON(http.StatusOK, payments)
	}
}

func UpdatePayment(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment ID"})
			return
		}

		var req model.UpdatePaymentRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		paymentRepo := repository.NewPaymentRepository(db)
		if err := paymentRepo.Update(id, req.Amount, req.PaidAt, req.Note, "admin"); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update payment"})
			return
		}

		// 更新後のレコード取得
		payment, err := paymentRepo.GetByID(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}

		// ユーザー名を取得
		userRepo := repository.NewUserRepository(db)
		user, _ := userRepo.GetByID(payment.UserID)

		response := model.PaymentResponse{
			ID:        payment.ID,
			UserID:    payment.UserID,
			UserName:  user.Name,
			Amount:    payment.Amount,
			PaidAt:    payment.PaidAt,
			Note:      payment.Note,
			CreatedBy: payment.CreatedBy,
			UpdatedBy: payment.UpdatedBy,
			CreatedAt: payment.CreatedAt,
			UpdatedAt: payment.UpdatedAt,
		}

		c.JSON(http.StatusOK, response)
	}
}

func DeletePayment(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid payment ID"})
			return
		}

		paymentRepo := repository.NewPaymentRepository(db)
		if err := paymentRepo.Delete(id); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete payment"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Payment deleted"})
	}
}
