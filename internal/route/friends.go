package route

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/ranktify/ranktify-be/internal/dao"
	"github.com/ranktify/ranktify-be/internal/handler"
	"github.com/ranktify/ranktify-be/internal/middleware"
)

func FriendRoutes(group *gin.RouterGroup, db *sql.DB) {
	friendDAO := dao.NewFriendsDAO(db)
	friendsHandler := handler.NewFriendHandler(friendDAO)

	friends := group.Group("/friends")
	{
		friends.Use(middleware.AuthMiddleware())
		friends.GET("/top-tracks", friendsHandler.GetTop5TracksAmongFriends)
		// Routes to manage Friends
		friends.GET("/:user_id", friendsHandler.GetFriends)  
		friends.DELETE("/:user_id/:friend_id", friendsHandler.DeleteFriendByID) 
		// Routes to manage Friend Requests
		friends.POST("/send/:user_id/:receiver_id", friendsHandler.SendFriendRequest)  
		friends.POST("/accept/:request_id", friendsHandler.AcceptFriendRequest) 
		friends.DELETE("/decline/:request_id", friendsHandler.DeclineFriendRequest) 
		friends.DELETE("/friend-request/:request_id", friendsHandler.DeleteFriendRequest)
		// User Notifications
		friends.GET("/friend-requests", friendsHandler.GetFriendRequests)
	}
}
