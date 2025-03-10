package middleware

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ranktify/ranktify-be/internal/jwt"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := c.Cookie("jwt_token")
		if err != nil {
			fmt.Println("Token missing in cookies")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		token, err := jwt.VerifyToken(tokenString)
		if err != nil {
			fmt.Println("Token verification failed")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
			c.Abort()
			return
		}

		fmt.Printf("Token verified: %+v\\n", token.Claims)
		c.Next()
	}
}
