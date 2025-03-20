package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ranktify/ranktify-be/internal/jwt"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")
		if tokenString == "" {
			fmt.Println("Token missing in headers")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}
		trimTokenString := strings.TrimPrefix(tokenString, "Bearer ")
		if trimTokenString == tokenString {
			fmt.Println("Token not found in headers")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		token, err := jwt.ValidateAccessToken(trimTokenString)
		if err != nil {
			fmt.Println("Token verification failed")
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		fmt.Printf("Token verified: %+v\\n", token.Claims)
		c.Next()
	}
}
