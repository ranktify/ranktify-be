package main

import (
	"github.com/gin-gonic/gin"
	"github.com/ranktify/ranktify-be/internal/route"
)

func main() {
	router := gin.Default()
	router.RemoveExtraSlash = true

	mainGroup := router.Group("")
	{
		route.UserRoutes(mainGroup)
	}

	router.Run("localhost:8080")
}
