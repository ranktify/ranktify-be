package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/ranktify/ranktify-be/internal/dao"
)

type RankingsHandler struct {
	DAO *dao.RankingsDao
}

func NewRankingsHandler(dao *dao.RankingsDao) *RankingsHandler {
	return &RankingsHandler{DAO: dao}
}

func (h *RankingsHandler) GetRankedSongs(c *gin.Context) {
	userID, err := strconv.ParseUint(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}
	rankings, err := h.DAO.GetRankedSongs(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve rankings"})
	}
	c.JSON(http.StatusOK, gin.H{"rankings": rankings})
}
