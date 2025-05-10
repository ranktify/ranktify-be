package route

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/ranktify/ranktify-be/internal/dao"
	"github.com/ranktify/ranktify-be/internal/handler"
	"github.com/ranktify/ranktify-be/internal/middleware"
)

func StreakRoutes(group *gin.RouterGroup, db *sql.DB) {
	dao := dao.NewStreaksDAO(db)
	handler := handler.NewStreaksHandler(dao)

	streaks := group.Group("/streaks")
	{
		streaks.Use(middleware.AuthMiddleware())
		streaks.GET("", handler.GetStreaksByUserID)
	}
}
