package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ranktify/ranktify-be/config"
	"github.com/ranktify/ranktify-be/internal/route"
)

func main() {
	router := gin.Default()
	router.RemoveExtraSlash = true
	db := config.SetupConnection()

	mainGroup := router.Group("/ranktify")
	{
		route.UserRoutes(mainGroup, db)
	}

	router.Run("localhost:8080")
}
