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
)

func getSpotifyClientID() string {
	return os.Getenv("SPOTIFY_CLIENT_ID")
}

func getSpotifySecret() string {
	return os.Getenv("SPOTIFY_SECRET")
}

type SpotifyHandler struct {
	DAO *dao.TokensDAO
}

func NewSpotifyHandler(dao *dao.TokensDAO) *SpotifyHandler {
	return &SpotifyHandler{DAO: dao}
}

var BaseUrl string = "https://accounts.spotify.com"

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

type SpotifyAuthCallbackResponse struct {
	Code  string `json:"code"`  // An authorization code that can be exchanged for an access token.
	State string `json:"state"` // The value of the state parameter supplied in the request.
	Err   string `json:"error"` // The reason authorization failed, for example: "access_denied"
}

type SpotifyAccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
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
	fmt.Println("State:", authCallbackResponse.State)

	formData := url.Values{}
	formData.Set("grant_type", "authorization_code")
	formData.Set("code", authCallbackResponse.Code)
	formData.Set("redirect_uri", "exp://127.0.0.1:19000/")

	fmt.Println("Code", authCallbackResponse.Code)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, BaseUrl+"/api/token", strings.NewReader(formData.Encode()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	credentials := fmt.Sprintf("%s:%s", getSpotifyClientID(), getSpotifySecret())
	encodedCredentials := base64.StdEncoding.EncodeToString([]byte(credentials))
	req.Header.Add("Authorization", "Basic "+encodedCredentials)

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := httpClient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer resp.Body.Close()

	fmt.Println("Status code:", resp.StatusCode)
	if resp.StatusCode != http.StatusOK {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch access token"})
		return
	}

	var tokenResponse SpotifyAccessTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResponse); err != nil {
		fmt.Println(err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("Access Token: ", tokenResponse.AccessToken)
	fmt.Println("Refresh Token: ", tokenResponse.RefreshToken)
	fmt.Println("Expires in: ", tokenResponse.ExpiresIn)

	// TODO: insert the refresh token in the database
	c.JSON(http.StatusOK, gin.H{"access_token": tokenResponse.AccessToken})

}

func (h *SpotifyHandler) RequestAcessToken(c *gin.Context) {
	// Fetch the user spotify refresh token
	// request an access token w/ refresh token
	// return the access token to the client
}
