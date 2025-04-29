package handler

import (
	"database/sql"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/ranktify/ranktify-be/internal/dao"
	"github.com/ranktify/ranktify-be/internal/model"
	"github.com/ranktify/ranktify-be/internal/spotify"
)

type SpotifyHandler struct {
	DAO *dao.SpotifyDAO
}

func NewSpotifyHandler(dao *dao.SpotifyDAO) *SpotifyHandler {
	return &SpotifyHandler{DAO: dao}
}

// Receives the auth code to perform the final step of authorization code
func (h *SpotifyHandler) AuthCallback(c *gin.Context) {
	var authCallbackResponse spotify.SpotifyAuthCallbackResponse

	if err := c.ShouldBind(&authCallbackResponse); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if authCallbackResponse.Err != "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": authCallbackResponse.Err})
		return
	}

	formData := url.Values{}
	formData.Set("grant_type", "authorization_code")
	formData.Set("code", authCallbackResponse.Code)
	formData.Set("redirect_uri", spotify.GetSpotifyRedirectURI())

	tokenResponse, err := spotify.FetchSpotifyToken(c.Request.Context(), formData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	rt := model.SpotifyRefreshToken{
		UserID: authCallbackResponse.UserID,
		Token:  tokenResponse.RefreshToken,
	}
	if err := h.DAO.SaveRefreshToken(rt); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"access_token": tokenResponse.AccessToken})
}

// Returns an access token given a user_id
func (h *SpotifyHandler) RefreshAccessToken(c *gin.Context) {
	var refreshRequest spotify.SpotifyRefreshTokenRequest
	if err := c.ShouldBindJSON(&refreshRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	refreshToken, err := h.DAO.GetRefreshToken(refreshRequest.UserID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "refresh token not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	fmt.Println(refreshToken)

	formData := url.Values{}
	formData.Set("grant_type", "refresh_token")
	formData.Set("refresh_token", refreshToken)

	response, err := spotify.FetchSpotifyToken(c.Request.Context(), formData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Println(response)
	// If the api returned a new rt then update it in the database
	if response.RefreshToken != "" {
		err := h.DAO.UpdateRefreshToken(refreshRequest.UserID, response.RefreshToken)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// Return the new access token
	c.JSON(http.StatusOK, gin.H{"access_token": response.AccessToken})

}

// Get N songs to rank for now N is constant at 50 this can be changed, even paginated w/ offset
func (h *SpotifyHandler) GetSongsToRank(c *gin.Context) {
	rawToken := c.GetHeader("Spotify-Token")
	if rawToken == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "No access token provided"})
		return
	}
	accessToken := strings.TrimPrefix(rawToken, "Bearer ") // Maybe we can add this to the middleware

	client := spotify.SpotifyClientFromAccessToken(c.Request.Context(), accessToken)

	songs, err := spotify.GetTopNSongs(c.Request.Context(), client, 5)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, songs)
}
