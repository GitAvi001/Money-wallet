package main

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	router := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:3000", "http://localhost:5173"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Origin", "Content-Type", "Authorization"}
	config.AllowCredentials = true
	router.Use(cors.New(config))

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "api-gateway"})
	})

	// Get service URLs from environment variables
	authServiceURL := getEnv("AUTH_SERVICE_URL", "http://localhost:8081")
	transactionServiceURL := getEnv("TRANSACTION_SERVICE_URL", "http://localhost:8082")

	// Transaction service routes (registering first - view the terminal when start the services)
	router.GET("/api/wallet", createSimpleProxy(transactionServiceURL, "/wallet"))
	router.POST("/api/wallet/add", createSimpleProxy(transactionServiceURL, "/wallet/add"))
	router.POST("/api/wallet/withdraw", createSimpleProxy(transactionServiceURL, "/wallet/withdraw"))
	router.POST("/api/transactions/transfer", createSimpleProxy(transactionServiceURL, "/transfer"))
	router.GET("/api/transactions", createSimpleProxy(transactionServiceURL, "/transactions"))
	router.GET("/api/transactions/:id", createSimpleProxy(transactionServiceURL, "/transactions"))

	// Auth service routes
	authGroup := router.Group("/api/auth")
	{
		authGroup.Any("/*path", proxyHandler(authServiceURL, ""))
	}

	// Start server
	port := getEnv("PORT", "8080")
	log.Printf("API Gateway starting on port %s", port)
	log.Printf("Auth Service URL: %s", authServiceURL)
	log.Printf("Transaction Service URL: %s", transactionServiceURL)

	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start API Gateway:", err)
	}
}

func proxyHandler(targetURL string, basePath ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		remote, err := url.Parse(targetURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Service configuration error"})
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(remote)
		proxy.Director = func(req *http.Request) {
			req.Header = c.Request.Header
			req.Host = remote.Host
			req.URL.Scheme = remote.Scheme
			req.URL.Host = remote.Host

			// For auth routes, use the path parameter
			path := c.Param("path")
			if path != "" {
				req.URL.Path = path
			} else {
				req.URL.Path = c.Request.URL.Path
			}
		}

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func createSimpleProxy(targetURL, targetPath string) gin.HandlerFunc {
	return func(c *gin.Context) {
		remote, err := url.Parse(targetURL)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Service configuration error"})
			return
		}

		proxy := httputil.NewSingleHostReverseProxy(remote)
		proxy.Director = func(req *http.Request) {
			req.URL.Scheme = remote.Scheme
			req.URL.Host = remote.Host
			req.URL.Path = targetPath
			req.Host = remote.Host

			// Copy headers
			for key, values := range c.Request.Header {
				for _, value := range values {
					req.Header.Set(key, value)
				}
			}

			// For routes with :id parameter
			if id := c.Param("id"); id != "" {
				req.URL.Path = targetPath + "/" + id
			}
		}

		proxy.ServeHTTP(c.Writer, c.Request)
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
