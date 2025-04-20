package route

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/ranktify/ranktify-be/internal/dao"
	"github.com/ranktify/ranktify-be/internal/handler"
)

func RankingsRoutes(group *gin.RouterGroup, db *sql.DB) {
	rankingsDAO := dao.NewRankingsDAO(db)
	rankingsHandler := handler.NewRankingsHandler(rankingsDAO)

	rankings := group.Group("/rankings")
	{
		rankings.GET("/:user_id", rankingsHandler.GetRankedSongs)
	}
}
