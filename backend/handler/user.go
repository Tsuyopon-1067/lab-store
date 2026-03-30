package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"purchase-system/repository"
)

type CreateUserRequest struct {
	Name    string `json:"name" binding:"required"`
	Barcode string `json:"barcode" binding:"required"`
}

type UpdateUserRequest struct {
	Name     string `json:"name" binding:"required"`
	IsActive *int   `json:"is_active"`
}

func GetUserByBarcode(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		barcode := c.Param("code")
		repo := repository.NewUserRepository(db)
		user, err := repo.GetByBarcode(barcode)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		if user == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
			return
		}
		c.JSON(http.StatusOK, user)
	}
}

func ListUsers(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		repo := repository.NewUserRepository(db)
		users, err := repo.ListAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		c.JSON(http.StatusOK, users)
	}
}

func CreateUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		repo := repository.NewUserRepository(db)
		user, err := repo.Create(req.Name, req.Barcode)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create user"})
			return
		}
		c.JSON(http.StatusCreated, user)
	}
}

func UpdateUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		var req UpdateUserRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		repo := repository.NewUserRepository(db)
		user, err := repo.Update(id, req.Name)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
			return
		}

		// is_active が送信された場合のみ更新
		if req.IsActive != nil {
			if err := repo.SetActive(id, *req.IsActive); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update is_active"})
				return
			}
			// SetActive後に最新データを取得
			user, err = repo.GetByID(id)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated user"})
				return
			}
		}

		c.JSON(http.StatusOK, user)
	}
}

func DeleteUser(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
			return
		}

		repo := repository.NewUserRepository(db)
		if err := repo.SetActive(id, 0); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
	}
}
