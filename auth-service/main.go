package main

import (
	"auth-service/config"
	"auth-service/handlers"
	"auth-service/middleware"
	"log"

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
		c.JSON(200, gin.H{"status": "ok", "service": "auth-service"})
	})

	// Public routes
	router.POST("/register", handlers.Register)
	router.POST("/login", handlers.Login)
	router.GET("/verify-email", handlers.VerifyEmail)
	router.POST("/send-verification", handlers.SendVerificationEmail)

	// Protected routes
	protected := router.Group("/")
	protected.Use(middleware.AuthMiddleware())
	{
		protected.GET("/me", handlers.GetCurrentUser)
		protected.GET("/users", handlers.GetAllUsers)
	}

	port := config.GetEnv("PORT", "8081")
	log.Printf("Auth Service starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start Auth Service:", err)
	}
}
