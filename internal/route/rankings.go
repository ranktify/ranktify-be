package route

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/ranktify/ranktify-be/internal/dao"
	"github.com/ranktify/ranktify-be/internal/handler"
	"github.com/ranktify/ranktify-be/internal/middleware"
	"github.com/ranktify/ranktify-be/internal/service"
)

func RankingsRoutes(group *gin.RouterGroup, db *sql.DB) {
	rankingsService := service.NewRankingsService(
		dao.NewRankingsDAO(db),
	)
	rankingsHandler := handler.NewRankingsHandler(rankingsService)

	rankings := group.Group("/rankings")
	{
		rankings.Use(middleware.AuthMiddleware())
		rankings.GET("/ranked-songs", rankingsHandler.GetRankedSongs)
		rankings.GET("/friends-ranked-songs", rankingsHandler.GetFriendsRankedSongs)
		rankings.GET("/friends-songs", rankingsHandler.GetFriendsRankedSongsWithNoUserRank)
		rankings.POST("/:song_id/:rank", rankingsHandler.RankSong)
		rankings.DELETE("/:ranking_id", rankingsHandler.DeleteRanking)
		rankings.PUT("/:ranking_id/:rank", rankingsHandler.UpdateRanking)
	}
}
