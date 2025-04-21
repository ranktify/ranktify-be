package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ranktify/ranktify-be/internal/dao"
	"github.com/ranktify/ranktify-be/internal/model"
)

type RankingsHandler struct {
	RankingDAO *dao.RankingsDao
	FriendsDAO *dao.FriendsDAO
}

func NewRankingsHandler(rankingsDao *dao.RankingsDao, friendsDao *dao.FriendsDAO) *RankingsHandler {
	return &RankingsHandler{
		RankingDAO: rankingsDao,
		FriendsDAO: friendsDao,
	}
}

func (h *RankingsHandler) GetRankedSongs(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	rankings, err := h.RankingDAO.GetRankedSongs(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve rankings"})
	}
	c.JSON(http.StatusOK, gin.H{"rankings": rankings})
}

func (h *RankingsHandler) GetFriendsRankedSongs(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	friends, err := h.FriendsDAO.GetFriends(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve friends"})
		return
	}
	var friendIDs []uint64
	for _, friend := range friends {
		friendIDs = append(friendIDs, uint64(friend.Id))
	}
	var rankedSongs []model.Rankings
	for _, friend := range friendIDs {
		rankedSong, err := h.RankingDAO.GetRankedSongs(uint64(friend))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve ranked songs"})
			return
		}
		rankedSongs = append(rankedSongs, rankedSong...)
	}
	c.JSON(http.StatusOK, gin.H{"User's friends rankings": rankedSongs})
}
