package route

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/ranktify/ranktify-be/internal/dao"
	"github.com/ranktify/ranktify-be/internal/handler"
)

func ApiRoutes(router *gin.RouterGroup, db *sql.DB) {
	tokensHandler := handler.NewTokensHandler(dao.NewTokensDAO(db), dao.NewUserDAO(db))

	api := router.Group("/api")
	{
		api.POST("/refresh", tokensHandler.Refresh)
	}
}
