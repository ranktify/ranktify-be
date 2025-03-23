package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ranktify/ranktify-be/internal/dao"
	"github.com/ranktify/ranktify-be/internal/jwt"
)

type TokensHandler struct {
	DAO     *dao.TokensDAO
	UserDAO *dao.UserDAO
}

func NewTokensHandler(dao *dao.TokensDAO, UserDAO *dao.UserDAO) *TokensHandler {
	return &TokensHandler{DAO: dao, UserDAO: UserDAO}
}

func (h *TokensHandler) Logout(c *gin.Context) {
	// TODO: Implement deleting the refresh token associated with the session: user_id from the access token?
}

func (h *TokensHandler) Refresh(c *gin.Context) {
	parsedRefreshToken, err := jwt.ParseRefreshTokenClaims(c.GetHeader("refresh_token"))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	// get the refresh token store in the database and validate
	rtDb, err := h.DAO.GetJWTRefreshTokenByJTI(parsedRefreshToken.JTI)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if rtDb.RefreshToken != parsedRefreshToken.RefreshToken {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	// Get the user
	user, err := h.UserDAO.GetUserByID(rtDb.UserID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	accessTokenString, rtString := jwt.CreateTokens(*user)
	// Store the refresh token in the database
	newRTString, err := jwt.ParseRefreshTokenClaims(rtString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// TODO: The below DAO operations can be changed to transactions for production code
	// since they are write operations

	err = h.DAO.SaveJWTRefreshToken(newRTString)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	// Delete the old refresh token
	err = h.DAO.DeleteJWTRefreshTokenByJTI(parsedRefreshToken.JTI)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"access_token": accessTokenString, "refresh_token": rtString})
}
