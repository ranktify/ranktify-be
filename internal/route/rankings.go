package route

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/ranktify/ranktify-be/internal/dao"
	"github.com/ranktify/ranktify-be/internal/handler"
	"github.com/ranktify/ranktify-be/internal/middleware"
)

func RankingsRoutes(group *gin.RouterGroup, db *sql.DB) {
	friendsDAO := dao.NewFriendsDAO(db)
	rankingsDAO := dao.NewRankingsDAO(db)
	rankingsHandler := handler.NewRankingsHandler(rankingsDAO, friendsDAO)

	rankings := group.Group("/rankings")
	{
		rankings.Use(middleware.AuthMiddleware())
		rankings.GET("/:user_id", rankingsHandler.GetRankedSongs)
		rankings.GET("/friends-ranked-songs/:user_id", rankingsHandler.GetFriendsRankedSongs)
	}
}
