package spotify

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/ranktify/ranktify-be/internal/model"
	"github.com/zmb3/spotify/v2"
)

var genres = []string{
	// US market
	"pop",
	"hip%20hop",
	"rock",
	"country",
	"edm",
	"latin",
	"r%20b",
	"reggae",
	"jazz",
	"classical",
	// Latin Market
	"reggaeton",
	"salsa",
	"bachata",
	"merengue",
	"cumbia",
	"ranchera",
	"mariachi",
	"vallenato",
	"tango",
	"bolero",
	"latin%20pop",
	"latin%20rock",
	"latin%20trap",
}

var (
	scdnMP3PreviewRegex = regexp.MustCompile(`https://p\.scdn\.co/mp3-preview/[^"' >)]+`)

	httpClient = &http.Client{
		Timeout: 10 * time.Second,
	}
)

func GetRandomGenre() string {
	index := rand.Intn(len(genres))

	return genres[index]
}

// returns top N songs from the user using the CurrentUsersTopTracks from zmb3 client
// the PreviewURI, SongID, CreatedAt are ignored here
func GetTopNSongs(ctx context.Context, client *spotify.Client, n int) ([]model.Song, error) {
	// We have tree options for terms: long_term, medium_term, and short_term
	// add spotify.Timerange(spotify.MediumTermRange) as parameter to "CurrentUsersTopTracks"
	results, err := client.CurrentUsersTopTracks(ctx, spotify.Limit(n), spotify.Timerange(spotify.LongTermRange))
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

		artistNames := make([]string, len(track.Artists))
		for j, a := range track.Artists {
			artistNames[j] = a.Name
		}
		prevURI := ScrapePreviewURI(ctx, client, track.Name, strings.Join(artistNames, ","))
		songs = append(songs, model.Song{
			SpotifyID:   track.ID.String(),
			Title:       track.Name,
			Artist:      artistNamePtr,
			Album:       albumNamePtr,
			ReleaseDate: releaseDatePtr,
			Genre:       nil, // genre not provided by track
			CoverURI:    coverURIPtr,
			PreviewURI:  &prevURI,
			// SongID, CreatedAt are ignored here
		})
	}

	return songs, nil
}

func extractSCDNLink(ctx context.Context, pageURL string) (string, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", pageURL, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request for %s: %w", pageURL, err)
	}
	// Spotify might block default Go user agent, pretend to be a browser
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	resp, err := httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to fetch %s: %w", pageURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Specifically check for 404 Not Found, which might be common
		if resp.StatusCode == http.StatusNotFound {
			log.Printf("Info: Spotify page not found (404) for %s", pageURL)
			return "", nil
		}
		// Handle potential rate limiting on scraping
		if resp.StatusCode == http.StatusTooManyRequests {
			log.Printf("Warning: Rate limited while scraping %s", pageURL)
			return "", fmt.Errorf("rate limited while scraping %s", pageURL)
		}
		return "", fmt.Errorf("failed to fetch %s: status code %d", pageURL, resp.StatusCode)
	}

	htmlBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response body from %s: %w", pageURL, err)
	}
	html := string(htmlBytes)

	// Find the *first* match for the specific regex pattern
	match := scdnMP3PreviewRegex.FindString(html)

	// Return the match (which will be an empty string if no match was found)
	return match, nil
}

func ScrapePreviewURI(ctx context.Context, client *spotify.Client, trackTitle string, trackArtist string) string {
	query := fmt.Sprintf("track:%s artist:%s", trackTitle, trackArtist)

	opts := []spotify.RequestOption{
		spotify.Limit(1), // Request only the top result
		spotify.Market(spotify.CountryUSA),
	}

	results, err := client.Search(ctx, query, spotify.SearchTypeTrack, opts...)
	if err != nil {
		log.Printf("Error searching Spotify: %v", err)
		return ""
	}

	if results.Tracks == nil || len(results.Tracks.Tracks) == 0 {
		log.Println("No tracks found matching the query.")
		// Return success, but with an empty results slice
		return ""
	}

	track := results.Tracks.Tracks[0]

	var spotifyURL string
	if url, ok := track.ExternalURLs["spotify"]; ok {
		spotifyURL = url
	}
	if spotifyURL == "" {
		log.Println("No Spotify URL found for the track.")
		return ""
	}

	audioURI, err := extractSCDNLink(ctx, spotifyURL)
	if err != nil {
		log.Printf("Error extracting SCDN links: %v", err)
		return ""
	}

	return audioURI
}

func StoreSongs(body []byte) ([]model.Song, error) {
	var searchResp SpotifySearchResponse

	err := json.Unmarshal(body, &searchResp)
	if err != nil {
		return nil, err
	}

	var songs []model.Song

	for _, item := range searchResp.Tracks.Items {
		var releaseDate *time.Time
		if item.Album.ReleaseDate != "" {
			t, err := time.Parse("2006-01-02", item.Album.ReleaseDate)
			if err == nil {
				releaseDate = &t
			}
		}

		var artistName *string
		if len(item.Artists) > 0 {
			artistName = &item.Artists[0].Name
		}

		var coverURI *string
		if len(item.Album.Images) > 0 {
			coverURI = &item.Album.Images[0].URL
		}

		song := model.Song{
			SpotifyID:   item.ID,
			Title:       item.Name,
			Artist:      artistName,
			Album:       &item.Album.Name,
			ReleaseDate: releaseDate,
			CoverURI:    coverURI,
			PreviewURI:  item.PreviewURL,
			CreatedAt:   time.Now(),
		}

		songs = append(songs, song)
	}
	return songs, nil
}

func GetRandomSongs(ctx context.Context, accessToken string, limit int) ([]model.Song, error) {
	wildcards := []string{
		"%25a%25", "a%25",
		"%25e%25", "e%25",
		"%25i%25", "i%25",
		"%25o%25", "o%25",
		"%25u%25", "u%25",
	}

	query := wildcards[rand.Intn(len(wildcards))]
	offset := rand.Intn(limit) + 1
	market := "US"

	url := fmt.Sprintf("https://api.spotify.com/v1/search?q=%s&offset=%d&limit=%d&type=track&market=%s", query, offset, limit, market)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	accessToken = strings.TrimPrefix(accessToken, "Bearer ")
	req.Header.Add("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch songs, status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	songs, err := StoreSongs(body)
	if err != nil {
		return nil, fmt.Errorf("failed to store songs")
	}

	return songs, nil
}

func GetRandomSongsByGenre(ctx context.Context, accessToken string, limit int, genre string) ([]model.Song, error) {
	wildcards := []string{
		"%25a%25", "a%25",
		"%25e%25", "e%25",
		"%25i%25", "i%25",
		"%25o%25", "o%25",
		"%25u%25", "u%25",
	}

	query := wildcards[rand.Intn(len(wildcards))]
	offset := rand.Intn(limit) + 1
	market := "US"

	url := fmt.Sprintf("https://api.spotify.com/v1/search?q=%s%%20genre:%%22%s%%22&offset=%d&limit=%d&type=track&market=%s", query, genre, offset, limit, market)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	accessToken = strings.TrimPrefix(accessToken, "Bearer ")
	req.Header.Add("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch songs, status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	songs, err := StoreSongs(body)
	if err != nil {
		return nil, fmt.Errorf("failed to store songs")
	}
	return songs, nil
}
