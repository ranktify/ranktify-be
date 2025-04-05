package route

import (
	"database/sql"

	"github.com/gin-gonic/gin"
	"github.com/ranktify/ranktify-be/internal/dao"
	"github.com/ranktify/ranktify-be/internal/handler"
)

func FriendRoutes(group *gin.RouterGroup, db *sql.DB) {
	friendDAO := dao.NewFriendsDAO(db)
	friendsHandler := handler.NewFriendHandler(friendDAO)

	friends := group.Group("/friends")
	{
		friends.GET("/:user_id/friend-list", friendsHandler.GetFriends)
		friends.POST("/:user_id/friend-requests/:receiver_id/send", friendsHandler.SendFriendRequest)
		friends.PUT("/:sender_id/friend-requests/:id/process", friendsHandler.ProcessFriendRequest)
		friends.DELETE("/:user_id/friend-requests/:request_id", friendsHandler.CancelFriendRequest)
		friends.DELETE("/:user_id/friends/:friend_id", friendsHandler.DeleteFriendByID)
	}
}
