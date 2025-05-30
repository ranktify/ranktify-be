package main

import (
	"flag"
	"os"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/ranktify/ranktify-be/config"
	"github.com/ranktify/ranktify-be/internal/jwt"
	"github.com/ranktify/ranktify-be/internal/route"
)

func main() {
	genTokensAndExit := flag.Bool("jwt", false, "Generate JWT tokens and terminates program")
	flag.Parse()

	if *genTokensAndExit {
		jwt.GenerateJWTKeys()
		return
	}

	router := gin.Default()

	// below is the cors setup for browser testing
	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8081"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour, // Preflight request caching
	}))

	// comment this out, if want to test in mobile devices
	// router.Use(cors.Default())
	router.RemoveExtraSlash = true
	db := config.SetupConnection()

	mainGroup := router.Group("/ranktify")
	{
		route.UserRoutes(mainGroup, db)
		route.FriendRoutes(mainGroup, db)
		route.ApiRoutes(mainGroup, db)
		route.RankingsRoutes(mainGroup, db)
		route.SongRecommendationRoutes(mainGroup, db)
		route.StreakRoutes(mainGroup, db)
		route.ImpressionRoutes(mainGroup, db)
	}
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := router.Run(":" + port); err != nil {
		panic(err)
	}
}
