package spotify

// --- Structs for Spotify API AUTH Response ---

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

type SpotifyRefreshTokenRequest struct { // TODO" We can eliminate this struct since we are passing the userID in the JWT token
	UserID uint64 `json:"user_id"`
}

// --- Structs for Spotify API Response ---

type SpotifySearchResponse struct {
	Tracks SpotifyTracksObject `json:"tracks"`
}

type SpotifyTracksObject struct {
	Items []SpotifyTrack `json:"items"`
}

type SpotifyTrack struct {
	ID           string              `json:"id"`
	Name         string              `json:"name"`
	Artists      []SpotifyArtist     `json:"artists"`
	Album        SpotifyAlbum        `json:"album"`
	ExternalURLs SpotifyExternalURLs `json:"external_urls"`
	PreviewURL   *string             `json:"preview_url"` // Might be null
}

type SpotifyArtist struct {
	Name string `json:"name"`
}

type SpotifyAlbum struct {
	Name        string         `json:"name"`
	ReleaseDate string         `json:"release_date"`
	Images      []SpotifyImage `json:"images"`
}

type SpotifyImage struct {
	URL    string `json:"url"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

type SpotifyExternalURLs struct {
	Spotify string `json:"spotify"`
}

// --- Structs for Final Result ---

type SearchResult struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Artists    string `json:"artists"`
	ImageURL   string `json:"imageUrl"`
	SpotifyURL string `json:"spotifyUrl"`
	AudioURI   string `json:"audioUri"` // The final URI (scraped or preview_url)
}

type SearchResponse struct {
	Success bool           `json:"success"`
	Results []SearchResult `json:"results"`
	Error   string         `json:"error,omitempty"` // Only include error if present
}
