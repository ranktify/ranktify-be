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
	requestID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}
	senderID, err := strconv.ParseUint(c.Param("sender_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user/sender ID"})
		return
	}
	receiverID, err := strconv.ParseUint(c.Param("receiver_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user/receiver ID"})
		return
	}
	err = h.DAO.AcceptFriendRequest(senderID, receiverID)
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
	requestID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}
	_, err = strconv.ParseUint(c.Param("sender_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user/sender ID"})
		return
	}
	_, err = strconv.ParseUint(c.Param("receiver_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user/receiver ID"})
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
