package route

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/ranktify/ranktify-be/internal/dao"
	"github.com/ranktify/ranktify-be/internal/handler"
	"github.com/ranktify/ranktify-be/internal/middleware"
)

func ApiRoutes(router *gin.RouterGroup, db *sql.DB) {
	tokensHandler := handler.NewTokensHandler(dao.NewTokensDAO(db), dao.NewUserDAO(db))
	spotifyHandler := handler.NewSpotifyHandler(dao.NewSpotifyDAO(db))

	api := router.Group("/api")
	{
		// JWT
		api.POST("/refresh", tokensHandler.Refresh)

		// protected routes
		// Spotify
		api.Use(middleware.AuthMiddleware())
		api.POST("/callback", spotifyHandler.AuthCallback)
		api.POST("/spotify-refresh", spotifyHandler.RefreshAccessToken)
	}
}
