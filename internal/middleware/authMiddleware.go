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
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized, missing auth token in header"})
			return
		}
		trimTokenString := strings.TrimPrefix(tokenString, "Bearer ")
		if trimTokenString == "" {
			fmt.Println("Token not found in headers")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			return
		}

		token, err := jwt.ValidateAccessToken(trimTokenString)
		if err != nil {
			fmt.Println("Token verification failed")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}

		if claims, err := jwt.GetClaimsFromAccessToken(token); err != nil {
			fmt.Println("Error getting claims from token")
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		} else {
			c.Set("userId", claims.UserID)
			c.Set("username", claims.Username)
			c.Set("email", claims.Email)
		}

		c.Next()
	}
}
