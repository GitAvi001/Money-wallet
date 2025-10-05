package main

import (
	"log"
	"transaction-service/config"
	"transaction-service/handlers"
	"transaction-service/middleware"

	"github.com/gin-gonic/gin"
)

func main() {
	// Initialize database
	config.InitDB()
	defer config.CloseDB()

	// Run migrations
	config.RunMigrations()

	router := gin.Default()

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok", "service": "transaction-service"})
	})

	// Protected routes (require authentication)
	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		// Wallet routes
		protected.GET("/wallet", handlers.GetWallet)
		protected.POST("/wallet/add", handlers.AddFunds)
		protected.POST("/wallet/withdraw", handlers.WithdrawFunds)

		// Transaction routes
		protected.POST("/transfer", handlers.Transfer)
		protected.GET("/transactions", handlers.GetTransactions)
		protected.GET("/transactions/:id", handlers.GetTransactionByID)
	}

	port := config.GetEnv("PORT", "8082")
	log.Printf("Transaction Service starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start Transaction Service:", err)
	}
}
