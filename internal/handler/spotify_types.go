package handler

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
	PreviewURL   *string             `json:"preview_url"` // Use pointer for nullable fields
}

type SpotifyArtist struct {
	Name string `json:"name"`
}

type SpotifyAlbum struct {
	Name         string          `json:"name"`
	ReleaseDate  string          `json:"release_date"`
	Images       []SpotifyImage  `json:"images"`
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
	PreviewURL string `json:"previewUrl"` // The original preview_url from API or ""
	AudioURI   string `json:"audioUri"`   // The final URI (scraped or preview_url)
}

type SearchResponse struct {
	Success bool           `json:"success"`
	Results []SearchResult `json:"results"`
	Error   string         `json:"error,omitempty"` // Only include error if present
}
