package spotify

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/zmb3/spotify/v2"
	"golang.org/x/oauth2"
)

const tokenEndpoint string = "https://accounts.spotify.com/api/token"

func GetSpotifyClientID() string {
	return os.Getenv("SPOTIFY_CLIENT_ID")
}

func GetSpotifySecret() string {
	return os.Getenv("SPOTIFY_SECRET")
}

func GetSpotifyRedirectURI() string {
	return os.Getenv("SPOTIFY_REDIRECT_URI")
}

func FetchSpotifyToken(ctx context.Context, formData url.Values) (*SpotifyAccessTokenResponse, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, tokenEndpoint, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, err
	}

	credentials := fmt.Sprintf("%s:%s", GetSpotifyClientID(), GetSpotifySecret())
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
