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

func (h *FriendHandler) ProcessFriendRequest(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("sender_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	requestID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request ID"})
		return
	}
	// Extract action from the request body (expecting JSON payload with {"action": "accept"} or {"action": "reject"})
	var requestBody struct {
		Action string `json:"action"`
	}
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}
	// Validate action (should be "accept" or "reject")
	if requestBody.Action != "accept" && requestBody.Action != "reject" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid action. Must be 'accept' or 'reject'"})
		return
	}
	// Verify the friend request first
	receiverID, err := h.DAO.VerifyFriendRequest(requestID, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Friend request not found or invalid"})
		return
	} else {
		if requestBody.Action == "accept" {
			// Accept the friend request and add them as friends
			err = h.DAO.AcceptFriendRequest(userID, receiverID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to accept the friend request"})
				return
			}
			// Delete the friend request after successful acceptance
			err = h.DAO.DeleteFriendRequest(requestID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete the friend request"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "Friend request accepted"})
		} else if requestBody.Action == "reject" {
			// Reject the friend request by deleting it
			err := h.DAO.DeleteFriendRequest(requestID)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reject the friend request"})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "Friend request rejected"})
		}
	}

}

func (h *FriendHandler) DeleteFriendRequest(c *gin.Context) {
	requestID, err := strconv.ParseUint(c.Param("request_id"), 10, 64) // Extract request ID
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
