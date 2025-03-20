package route

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/ranktify/ranktify-be/internal/dao"
	"github.com/ranktify/ranktify-be/internal/handler"
	"github.com/ranktify/ranktify-be/internal/middleware"
	"github.com/ranktify/ranktify-be/internal/service"
)

func UserRoutes(group *gin.RouterGroup, db *sql.DB) {
	userService := service.NewUserService(
		dao.NewUserDAO(db),
		dao.NewTokensDAO(db),
	)
	userHandler := handler.NewUserHandler(userService)

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
