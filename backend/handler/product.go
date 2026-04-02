package handler

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"purchase-system/repository"
)

type CreateProductRequest struct {
	Name          string `json:"name" binding:"required"`
	Barcode       string `json:"barcode" binding:"required"`
	Price         int    `json:"price" binding:"required"`
	Note          string `json:"note"`
	StockQuantity int    `json:"stock_quantity"`
}

type UpdateProductRequest struct {
	Name          string `json:"name" binding:"required"`
	Barcode       string `json:"barcode" binding:"required"`
	Note          string `json:"note"`
	IsActive      *int   `json:"is_active"`
	StockQuantity *int   `json:"stock_quantity"`
}

type ChangePriceRequest struct {
	Price int `json:"price" binding:"required"`
}

func GetProductByBarcode(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		barcode := c.Param("code")
		repo := repository.NewProductRepository(db)
		product, err := repo.GetByBarcode(barcode)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		if product == nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "Product not found"})
			return
		}
		c.JSON(http.StatusOK, product)
	}
}

func ListProducts(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		repo := repository.NewProductRepository(db)
		products, err := repo.ListAll()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		c.JSON(http.StatusOK, products)
	}
}

func SearchProducts(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		query := c.Query("q")
		if query == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Search query is required"})
			return
		}

		repo := repository.NewProductRepository(db)
		products, err := repo.Search(query)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		c.JSON(http.StatusOK, products)
	}
}

func CreateProduct(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req CreateProductRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		repo := repository.NewProductRepository(db)
		product, err := repo.Create(req.Name, req.Barcode, req.Note, req.Price, req.StockQuantity)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create product"})
			return
		}
		c.JSON(http.StatusCreated, product)
	}
}

func UpdateProduct(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
			return
		}

		var req UpdateProductRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		repo := repository.NewProductRepository(db)
		product, err := repo.Update(id, req.Name, req.Barcode, req.Note, req.StockQuantity)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update product"})
			return
		}

		// is_active が送信された場合のみ更新
		if req.IsActive != nil {
			if err := repo.SetActive(id, *req.IsActive); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update is_active"})
				return
			}
			// SetActive後に最新データを取得（productWithPrice）
			updatedProduct, err := repo.GetByBarcode(product.Barcode)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch updated product"})
				return
			}
			product = updatedProduct
		}

		c.JSON(http.StatusOK, product)
	}
}

func ChangePrice(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
			return
		}

		var req ChangePriceRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		repo := repository.NewProductRepository(db)
		if err := repo.ChangePrice(id, req.Price); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to change price"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Price updated"})
	}
}

func GetPriceHistory(db *sql.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid product ID"})
			return
		}

		repo := repository.NewProductRepository(db)
		prices, err := repo.GetPriceHistory(id)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
			return
		}
		c.JSON(http.StatusOK, prices)
	}
}
