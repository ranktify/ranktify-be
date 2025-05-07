package route

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/ranktify/ranktify-be/internal/dao"
	"github.com/ranktify/ranktify-be/internal/handler"
)

func ImpressionRoutes(group *gin.RouterGroup, db *sql.DB) {
	impDAO := dao.NewImpressionDAO(db)
	impHandler := handler.NewImpressionHandler(impDAO)

	impression := group.Group("/impression")
	{
		impression.GET("/:label", impHandler.GetImpressionByLabel)
		impression.POST("/:label", impHandler.UpsertImpressionStats)
	}
}
