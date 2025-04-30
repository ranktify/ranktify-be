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
		rankings.GET("/:user_id", rankingsHandler.GetRankedSongs)
		rankings.GET("/friends-ranked-songs/:user_id", rankingsHandler.GetFriendsRankedSongs)
		rankings.GET("/friends-songs/:user_id", rankingsHandler.GetFriendsRankedSongsWithNoUserRank)
	}
}
