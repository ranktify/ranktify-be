package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ranktify/ranktify-be/internal/dao"
)

type FriendHandler struct {
	DAO *dao.FriendsDAO
}

func NewFriendHandler(dao *dao.FriendsDAO) *FriendHandler {
	return &FriendHandler{DAO: dao}
}

func (h *FriendHandler) GetFriends(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	friends, err := h.DAO.GetFriends(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve friends"})
		return
	}
	if len(friends) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("No friends found for user with id %d", userID)})
		return
	}
	c.JSON(http.StatusOK, gin.H{"friends": friends})
}

func (h *FriendHandler) GetFriendRequests(c *gin.Context) {
	rawUserID, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	userID := rawUserID.(uint64)
	friendRequests, friendRequestCount, err := h.DAO.GetFriendRequests(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve friend requests"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"friend_request":       friendRequests,
		"friend_request_count": friendRequestCount})
}

func (h *FriendHandler) GetFriendRequestsSent(c *gin.Context) {
	rawUserID, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	userID := rawUserID.(uint64)
	friendRequests, friendRequestCount, err := h.DAO.GetFriendRequestsSent(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve friend requests sent by current user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"friend_request":       friendRequests,
		"friend_request_count": friendRequestCount})
}

func (h *FriendHandler) DeleteFriendByID(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	friendID, err := strconv.ParseUint(c.Param("friend_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid friend ID"})
		return
	}
	err = h.DAO.DeleteFriendByID(userID, friendID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": fmt.Sprintf("No friendship found between user %d and %d", userID, friendID)})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete friend"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": "Friend deleted successfully"})
}

func (h *FriendHandler) SendFriendRequest(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	friendID, err := strconv.ParseUint(c.Param("receiver_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid friend ID"})
		return
	}
	err = h.DAO.SendFriendRequest(userID, friendID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send friend request"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Friend request sent successfully"})
}

func (h *FriendHandler) AcceptFriendRequest(c *gin.Context) {
	requestID, err := strconv.ParseUint(c.Param("request_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}
	friendRequest, err := h.DAO.GetFriendRequestsByRequestID(requestID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get the friend request"})
		return
	}
	err = h.DAO.AcceptFriendRequest(friendRequest.SenderID, friendRequest.ReceiverID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to accept the friend request"})
		return
	}
	err = h.DAO.DeleteFriendRequest(requestID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete the friend request"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Friend request accepted"})

}

func (h *FriendHandler) DeclineFriendRequest(c *gin.Context) {
	requestID, err := strconv.ParseUint(c.Param("request_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}
	err = h.DAO.DeleteFriendRequest(requestID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decline the friend request"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Friend request declined "})

}

func (h *FriendHandler) DeleteFriendRequest(c *gin.Context) {
	requestID, err := strconv.ParseUint(c.Param("request_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}
	err = h.DAO.DeleteFriendRequest(requestID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Friend request not found or already canceled"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to cancel friend request"})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Friend request canceled successfully"})
}

func (h *FriendHandler) GetTop5TracksAmongFriends(c *gin.Context) {
	userIDAny, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unauthorized"})
		return
	}
	userID := userIDAny.(uint64)

	topTracks, err := h.DAO.GetTopNTracksAmongFriends(c.Request.Context(), userID, 5)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve top tracks among friends"})
		return
	}
	if topTracks == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "No top tracks found among friends"})
		return
	}

	c.JSON(http.StatusOK, topTracks)
}
