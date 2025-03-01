package route

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/ranktify/ranktify-be/internal/dao"
	"github.com/ranktify/ranktify-be/internal/handler"
)

func UserRoutes(group *gin.RouterGroup, db *sql.DB) {
	userDAO := dao.NewUserDAO(db)
	userHandler := handler.NewUserHandler(userDAO)

	users := group.Group("/user")
	{
		users.POST("/login", userHandler.ValidateUser)
		users.POST("/register", userHandler.CreateUser)
	}
}
