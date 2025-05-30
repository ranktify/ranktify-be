package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ranktify/ranktify-be/internal/dao"
	"github.com/ranktify/ranktify-be/internal/model"
)

type ImpressionHandler struct {
	DAO *dao.ImpressionDAO
}

func NewImpressionHandler(dao *dao.ImpressionDAO) *ImpressionHandler {
	return &ImpressionHandler{DAO: dao}
}

func (h *ImpressionHandler) GetImpressionByLabel(c *gin.Context) {
	label := c.Param("label")
	if label == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "label is required"})
		return
	}
	imp, err := h.DAO.GetImpressionStatsByLabel(label)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	resp := struct {
		model.ImpressionStats
		Ctr float64 `json:"ctr"`
	}{
		ImpressionStats: *imp,
		Ctr:             imp.CTRPercent(),
	}

	c.JSON(http.StatusOK, resp)
}

func (h *ImpressionHandler) UpsertImpressionStats(c *gin.Context) {
	label := c.Param("label")
	if label == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "label is required"})
		return
	}
	var imp model.ImpressionStats
	if err := c.ShouldBindJSON(&imp); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	imp.ImpressionLabel = label

	if err := h.DAO.UpsertImpressionStats(&imp); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Impression stats updated successfully"})
}
