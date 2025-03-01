package main

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ranktify/ranktify-be/config"
	"github.com/ranktify/ranktify-be/internal/route"
)

func main() {
	router := gin.Default()

	// below is the setup cors for browser testing, this does not work with mobile devices
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8081"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour, // Preflight request caching
	}))

	router.RemoveExtraSlash = true
	db := config.SetupConnection()

	mainGroup := router.Group("/ranktify")
	{
		route.UserRoutes(mainGroup, db)
	}

	router.Run("0.0.0.0:8080")
}
