package route

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/ranktify/ranktify-be/internal/dao"
	"github.com/ranktify/ranktify-be/internal/handler"
	"github.com/ranktify/ranktify-be/internal/middleware"
)

func SongRecommendationRoutes(router *gin.RouterGroup, db *sql.DB) {
	rankingDAO := dao.NewRankingsDAO(db)
	songRecommendationHandler := handler.NewSongRecommendationHandler(rankingDAO)

	songRecommendation := router.Group("/song-recommendation")
	{
		songRecommendation.Use(middleware.SpotifyTokenMiddleware())
		songRecommendation.GET("/:user_id/:limit", songRecommendationHandler.SongRecommendation)
	}
}
