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
	"purchase-system/service"
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
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With, X-User-Barcode, X-API-Key")
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
	r.GET("/api/products/search", handler.SearchProducts(database))
	r.GET("/api/products", handler.ListProducts(database))
	r.POST("/api/purchases", handler.CreatePurchase(database))
	r.POST("/api/restocks", handler.CreateRestock(database))

	// 管理者認証
	r.POST("/api/auth/login", handler.Login(database, cfg))

	// 一般利用者向け（バーコード認証）
	me := r.Group("/api/me", middleware.BarcodeAuth(database))
	{
		me.GET("/balance", handler.GetMyBalance(database))
		me.GET("/purchases", handler.GetMyPurchases(database))
		me.GET("/restocks", handler.GetMyRestocks(database))
	}

	// 外部API（APIキー認証）
	v1 := r.Group("/v1", middleware.APIKeyAuth(database))
	{
		v1.GET("/users/:barcode/balance", handler.GetUserBalance(database))
		v1.GET("/reports/monthly", handler.GetMonthlyReport(database))
		v1.GET("/products", handler.ListProductsExternal(database))
	}

	// 管理者用エンドポイント
	admin := r.Group("/api", middleware.AdminAuth(database))
	{
		admin.POST("/auth/logout", handler.Logout(database))

		admin.GET("/users", handler.ListUsers(database))
		admin.POST("/users", handler.CreateUser(database))
		admin.PUT("/users/:id", handler.UpdateUser(database))
		admin.DELETE("/users/:id", handler.DeleteUser(database))

		admin.POST("/products", handler.CreateProduct(database))
		admin.PUT("/products/:id", handler.UpdateProduct(database))
		admin.POST("/products/:id/change-price", handler.ChangePrice(database))
		admin.GET("/products/:id/prices", handler.GetPriceHistory(database))

		admin.GET("/purchases", handler.ListPurchases(database))
		admin.GET("/purchases/summary", handler.GetSummary(database))

		admin.GET("/restocks", handler.ListRestocks(database))
		admin.PUT("/restocks/:id", handler.UpdateRestock(database))
		admin.DELETE("/restocks/:id", handler.DeleteRestock(database))

		admin.POST("/payments", handler.CreatePayment(database))
		admin.GET("/payments", handler.ListPayments(database))
		admin.PUT("/payments/:id", handler.UpdatePayment(database))
		admin.DELETE("/payments/:id", handler.DeletePayment(database))

		admin.POST("/restock-payments", handler.CreateRestockPayment(database))
		admin.GET("/restock-payments", handler.ListRestockPayments(database))
		admin.PUT("/restock-payments/:id", handler.UpdateRestockPayment(database))
		admin.DELETE("/restock-payments/:id", handler.DeleteRestockPayment(database))

		// APIキー管理
		admin.POST("/api-keys", handler.CreateAPIKey(database))
		admin.GET("/api-keys", handler.ListAPIKeys(database))
		admin.DELETE("/api-keys/:id", handler.DeleteAPIKey(database))

		// バックアップ
		admin.POST("/backup", handler.CreateBackup(database, cfg))
		admin.GET("/backup/list", handler.ListBackups(database, cfg))
	}

	// ヘルスチェック
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// 定期バックアップスケジューラー起動
	service.StartDailyBackupScheduler(database, cfg)

	// サーバー起動
	port := fmt.Sprintf(":%d", cfg.Server.Port)
	log.Printf("Starting server on %s", port)
	if err := r.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
