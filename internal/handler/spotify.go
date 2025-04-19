package handler

import (
	"context"
	"database/sql"
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
	"github.com/zmb3/spotify/v2"
	"golang.org/x/oauth2"
)

var tokenEndpoint string = "https://accounts.spotify.com/api/token"

var httpClient = &http.Client{
	Timeout: 10 * time.Second,
}

type SpotifyHandler struct {
	DAO *dao.SpotifyDAO
}

func NewSpotifyHandler(dao *dao.SpotifyDAO) *SpotifyHandler {
	return &SpotifyHandler{DAO: dao}
}

type SpotifyAuthCallbackResponse struct {
	Code   string `json:"code"`    // An authorization code that can be exchanged for an access token.
	State  string `json:"state"`   // The value of the state parameter supplied in the request.
	Err    string `json:"error"`   // The reason authorization failed, for example: "access_denied"
	UserID uint64 `json:"user_id"` // The ranktify user_id, to relate rt w/ user
}

type SpotifyAccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	Scope        string `json:"scope"`
}

type SpotifyRefreshTokenRequest struct {
	UserID uint64 `json:"user_id"`
}

func getSpotifyClientID() string {
	return os.Getenv("SPOTIFY_CLIENT_ID")
}

func getSpotifySecret() string {
	return os.Getenv("SPOTIFY_SECRET")
}

func getSpotifyRedirectURI() string {
	return os.Getenv("SPOTIFY_REDIRECT_URI")
}

func fetchSpotifyToken(ctx context.Context, formData url.Values) (*SpotifyAccessTokenResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenEndpoint, strings.NewReader(formData.Encode()))
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

// Uses zmb3 spotify wrapper
func SpotifyClientFromAccessToken(ctx context.Context, accessToken string) *spotify.Client {
	src := oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: accessToken,
		TokenType:   "Bearer",
	})

	return spotify.New(oauth2.NewClient(ctx, src))
}

// returns top N songs from the user using the CurrentUsersTopTracks from zmb3 client
// the PreviewURI, SongID, CreatedAt are ignored here
func getTopNSongs(ctx context.Context, client *spotify.Client, n int) ([]model.Song, error) {
	// We have tree options for terms: long_term, medium_term, and short_term
	// add spotify.Timerange(spotify.MediumTermRange) as parameter to "CurrentUsersTopTracks"
	results, err := client.CurrentUsersTopTracks(ctx, spotify.Limit(n))
	if err != nil {
		return nil, err
	}

	songs := make([]model.Song, 0, len(results.Tracks))

	for _, track := range results.Tracks {
		var artistNamePtr *string
		if len(track.Artists) > 0 {
			name := track.Artists[0].Name
			artistNamePtr = &name
		}

		var albumNamePtr *string
		if track.Album.Name != "" {
			album := track.Album.Name
			albumNamePtr = &album
		}

		var releaseDatePtr *time.Time
		if track.Album.ReleaseDate != "" {
			// Spotify gives release_date in YYYY[-MM[-DD]]
			// Try parsing full date
			if t, err := time.Parse("2006-01-02", track.Album.ReleaseDate); err == nil {
				releaseDatePtr = &t
			} else if t, err := time.Parse("2006-01", track.Album.ReleaseDate); err == nil {
				releaseDatePtr = &t
			} else if t, err := time.Parse("2006", track.Album.ReleaseDate); err == nil {
				releaseDatePtr = &t
			}
		}

		// Cover image URI (take first image if available)
		var coverURIPtr *string
		if len(track.Album.Images) > 0 {
			uri := track.Album.Images[0].URL
			coverURIPtr = &uri
		}

		songs = append(songs, model.Song{
			SpotifyID:   track.ID.String(),
			Title:       track.Name,
			Artist:      artistNamePtr,
			Album:       albumNamePtr,
			ReleaseDate: releaseDatePtr,
			Genre:       nil, // genre not provided by track
			CoverURI:    coverURIPtr,
			// PreviewURI, SongID, CreatedAt are ignored here
		})
	}

	return songs, nil
}

// Receives the auth code to perform the final step of authorization code
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

	formData := url.Values{}
	formData.Set("grant_type", "authorization_code")
	formData.Set("code", authCallbackResponse.Code)
	formData.Set("redirect_uri", getSpotifyRedirectURI())

	tokenResponse, err := fetchSpotifyToken(c.Request.Context(), formData)
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
	var refreshRequest SpotifyRefreshTokenRequest
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

	response, err := fetchSpotifyToken(c.Request.Context(), formData)
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

	client := SpotifyClientFromAccessToken(c.Request.Context(), accessToken)

	songs, err := getTopNSongs(c.Request.Context(), client, 50)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, songs)
}
