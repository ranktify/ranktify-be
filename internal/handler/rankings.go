package handler

import (
	"database/sql"
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
	rawUserID, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	userID := rawUserID.(uint64)
	statusCode, content := h.Service.GetRankedSongs(userID)
	c.JSON(statusCode, content)
}

func (h *RankingsHandler) GetFriendsRankedSongs(c *gin.Context) {
	rawUserID, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	userID := rawUserID.(uint64)
	statusCode, content := h.Service.GetFriendsRankedSongs(userID)
	c.JSON(statusCode, content)
}

func (h *RankingsHandler) GetFriendsRankedSongsWithNoUserRank(c *gin.Context) {
	rawUserID, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	userID := rawUserID.(uint64)
	statusCode, content := h.Service.GetFriendsRankedSongsWithNoUserRank(userID)
	c.JSON(statusCode, content)
}

func (h *RankingsHandler) GetTopWeeklyTracks(c *gin.Context) {
	songs, err := h.Service.GetTopWeeklyRankedSongs(c.Request.Context())
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "No tracks found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve top weekly tracks"})
		}
		return
	}
	c.JSON(http.StatusOK, songs)
}
