package route

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func UserRoutes(group *gin.RouterGroup) {
	users := group.Group("/user")
	{
		users.GET("/login", func(c *gin.Context) {
			c.JSON(http.StatusAccepted, gin.H{"message": "HELLO RANKTIFY"})
		})
	}
}
