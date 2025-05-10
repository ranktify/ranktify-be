package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ranktify/ranktify-be/internal/dao"
)

type StreaksHandler struct {
	DAO *dao.StreaksDAO
}

func NewStreaksHandler(dao *dao.StreaksDAO) *StreaksHandler {
	return &StreaksHandler{DAO: dao}
}

func (h *StreaksHandler) GetStreaksByUserID(c *gin.Context) {
	rawUserID, ok := c.Get("userId")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := rawUserID.(uint64)
	// Get streaks from the database
	streaks, err := h.DAO.GetStreaksByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get streaks"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"streaks": streaks})
}
