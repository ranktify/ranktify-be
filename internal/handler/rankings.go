package handler

import (
	"net/http"
	"strconv"

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

func (h *RankingsHandler) RankSong(c *gin.Context) {
	rawUserID, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	songID, err := strconv.ParseUint(c.Param("song_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid song ID"})
		return
	}
	rank, err := strconv.Atoi(c.Param("rank"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rank"})
		return
	}
	userID := rawUserID.(uint64)
	statusCode, content := h.Service.RankSong(songID, userID, rank)
	c.JSON(statusCode, content)
}

func (h *RankingsHandler) DeleteRanking(c *gin.Context) {
	rankingID, err := strconv.ParseUint(c.Param("ranking_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ranking ID"})
		return
	}
	statusCode, content := h.Service.DeleteRanking(rankingID)
	c.JSON(statusCode, content)
}

func (h *RankingsHandler) UpdateRanking(c *gin.Context) {
	rankingID, err := strconv.ParseUint(c.Param("ranking_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ranking ID"})
		return
	}
	rank, err := strconv.Atoi(c.Param("rank"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rank"})
		return
	}
	statusCode, content := h.Service.UpdateRanking(rankingID, rank)
	c.JSON(statusCode, content)
}
