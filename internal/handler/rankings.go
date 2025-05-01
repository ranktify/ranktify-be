package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ranktify/ranktify-be/internal/service"
)

type RankingsHandler struct {
	Service *service.RankingsService
}

func NewRankingsHandler(service *service.RankingsService) *RankingsHandler {
	return &RankingsHandler{Service: service}
}

func (h *RankingsHandler) GetRankedSongs(c *gin.Context) {
	rawUserID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	userID := rawUserID.(uint64)
	statusCode, content := h.Service.GetRankedSongs(userID)
	c.JSON(statusCode, content)
}

func (h *RankingsHandler) GetFriendsRankedSongs(c *gin.Context) {
	rawUserID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	userID := rawUserID.(uint64)
	statusCode, content := h.Service.GetFriendsRankedSongs(userID)
	c.JSON(statusCode, content)
}

func (h *RankingsHandler) GetFriendsRankedSongsWithNoUserRank(c *gin.Context) {
	rawUserID, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	userID := rawUserID.(uint64)
	statusCode, content := h.Service.GetFriendsRankedSongsWithNoUserRank(userID)
	c.JSON(statusCode, content)
}
