package handler

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ranktify/ranktify-be/internal/dao"
	"github.com/ranktify/ranktify-be/internal/model"
)

func getSpotifyClientID() string {
	return os.Getenv("SPOTIFY_CLIENT_ID")
}

func getSpotifySecret() string {
	return os.Getenv("SPOTIFY_SECRET")
}

func getSpotifyRedirectURI() string {
	return os.Getenv("SPOTIFY_REDIRECT_URI")
}

type SpotifyHandler struct {
	DAO *dao.SpotifyDAO
}

func NewSpotifyHandler(dao *dao.SpotifyDAO) *SpotifyHandler {
	return &SpotifyHandler{DAO: dao}
}

var BaseUrl string = "https://accounts.spotify.com"

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

type SpotifyAuthCallbackResponse struct {
	Code   string `json:"code"`  // An authorization code that can be exchanged for an access token.
	State  string `json:"state"` // The value of the state parameter supplied in the request.
	Err    string `json:"error"` // The reason authorization failed, for example: "access_denied"
	UserID uint64 `json:"user_id"`
}

type SpotifyAccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

func fetchSpotifyToken(ctx context.Context, endpoint string, formData url.Values) (*SpotifyAccessTokenResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, BaseUrl+endpoint, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, err
	}

	credentials := fmt.Sprintf("%s:%s", getSpotifyClientID(), getSpotifySecret())
	encodedCredentials := base64.StdEncoding.EncodeToString([]byte(credentials))
	req.Header.Add("Authorization", "Basic "+encodedCredentials)

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch access token, status code not ok")
	}

	var tokenResponse SpotifyAccessTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		return nil, err
	}
	return &tokenResponse, nil
}

func (h *SpotifyHandler) AuthCallback(c *gin.Context) {
	var authCallbackResponse SpotifyAuthCallbackResponse

	if err := c.ShouldBind(&authCallbackResponse); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if authCallbackResponse.Err != "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": authCallbackResponse.Err})
		return
	}
	// TODO: As soon as the frontend establishes a state then validate the state provided to prevent cross-origin requests

	ctx := context.Background()
	formData := url.Values{}
	formData.Set("grant_type", "authorization_code")
	formData.Set("code", authCallbackResponse.Code)
	formData.Set("redirect_uri", getSpotifyRedirectURI())

	tokenResponse, err := fetchSpotifyToken(ctx, "/api/token", formData)
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

func (h *SpotifyHandler) refreshAccessToken(c *gin.Context) {
	// Fetch the user spotify refresh token
	// request an access token w/ refresh token
	// return the access token to the client
}
