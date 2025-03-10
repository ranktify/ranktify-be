package route

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/ranktify/ranktify-be/internal/dao"
	"github.com/ranktify/ranktify-be/internal/handler"
	"github.com/ranktify/ranktify-be/internal/middleware"
)

func UserRoutes(group *gin.RouterGroup, db *sql.DB) {
	userDAO := dao.NewUserDAO(db)
	userHandler := handler.NewUserHandler(userDAO)

	users := group.Group("/user")
	{
		users.POST("/login", userHandler.ValidateUser)
		users.POST("/register", userHandler.CreateUser)

		//add authentication to the rest of the routes
		users.Use(middleware.AuthMiddleware())
		users.GET("/:id", userHandler.GetUserByID)
		users.GET("/", userHandler.GetAllUsers)
		users.PUT("/:id", userHandler.UpdateUserByID)
		users.DELETE("/:id", userHandler.DeleteUserByID)
	}
}
