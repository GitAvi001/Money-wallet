package middleware

import (
	"auth-service/utils" //Use for JWT validation
	"net/http"           //Use for HTTP status codes
	"strings"            //Use for string operations

	"github.com/gin-gonic/gin" //web framework for go instead of httprouter
)

// AuthMiddleware checks if the user is authenticated before goes with desired endpoint
// INtercepts with requests before reach to the actual user
// Handler function is an implicit interface
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization") //Authorization header accepts with Bearer <Access token> format
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		token := tokenParts[1]
		claims, err := utils.ValidateJWT(token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort() //stop chain
			return
		}

		c.Set("user_id", claims.UserID)
		c.Set("email", claims.Email)
		c.Next() //continue chain
	}
}
