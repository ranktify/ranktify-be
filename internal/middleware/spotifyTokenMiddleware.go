package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func SpotifyTokenMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		rawSpotifyToken := c.GetHeader("Spotify-Token")
		if rawSpotifyToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Spotify-Token header is required"})
			return
		}
		accessToken := strings.TrimPrefix(rawSpotifyToken, "Bearer ")
		if accessToken == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Spotify-Token header is required"})
			return
		}

		c.Set("spotifyToken", accessToken)
		c.Next()
	}
}
