package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"purchase-system/config"
	"purchase-system/db"
	"purchase-system/handler"
	"purchase-system/middleware"
)

func main() {
	// 設定読み込み
	configPath := "config.yaml"
	if _, err := os.Stat(configPath); err != nil {
		// リポジトリルートから実行されている場合
		configPath = "../config.yaml"
	}
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}
	cfg.SetDefaults()

	// DB初期化
	database, err := db.Init(cfg.Database.Path)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer database.Close()

	// Ginエンジン初期化
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	// CORS設定
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-User-Barcode")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// 認証不要なエンドポイント
	r.GET("/api/users/barcode/:code", handler.GetUserByBarcode(database))
	r.GET("/api/products/barcode/:code", handler.GetProductByBarcode(database))
	r.GET("/api/products", handler.ListProducts(database))
	r.POST("/api/purchases", handler.CreatePurchase(database))
	r.POST("/api/restocks", handler.CreateRestock(database))

	// 一般利用者向け（バーコード認証）
	me := r.Group("/api/me", middleware.BarcodeAuth(database))
	{
		me.GET("/balance", handler.GetMyBalance(database))
		me.GET("/purchases", handler.GetMyPurchases(database))
		me.GET("/restocks", handler.GetMyRestocks(database))
	}

	// 管理者用エンドポイント
	admin := r.Group("/api")
	{
		admin.GET("/users", handler.ListUsers(database))
		admin.POST("/users", handler.CreateUser(database))
		admin.PUT("/users/:id", handler.UpdateUser(database))
		admin.DELETE("/users/:id", handler.DeleteUser(database))

		admin.POST("/products", handler.CreateProduct(database))
		admin.PUT("/products/:id", handler.UpdateProduct(database))
		admin.POST("/products/:id/change-price", handler.ChangePrice(database))
		admin.GET("/products/:id/prices", handler.GetPriceHistory(database))

		admin.GET("/purchases", handler.ListPurchases(database))

		admin.GET("/restocks", handler.ListRestocks(database))
		admin.PUT("/restocks/:id", handler.UpdateRestock(database))
		admin.DELETE("/restocks/:id", handler.DeleteRestock(database))

		admin.POST("/payments", handler.CreatePayment(database))
		admin.GET("/payments", handler.ListPayments(database))
		admin.PUT("/payments/:id", handler.UpdatePayment(database))
		admin.DELETE("/payments/:id", handler.DeletePayment(database))
	}

	// ヘルスチェック
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// サーバー起動
	port := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Starting server on %s", port)
	if err := r.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
