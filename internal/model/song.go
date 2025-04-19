package model

import "time"

type Song struct {
	SongID      uint64     `json:"song_id"`
	SpotifyID   string     `json:"spotify_id"`
	Title       string     `json:"title"`
	Artist      *string    `json:"artist,omitempty"`
	Album       *string    `json:"album,omitempty"`
	ReleaseDate *time.Time `json:"release_date,omitempty"`
	Genre       *string    `json:"genre,omitempty"`
	CoverURI    *string    `json:"cover_uri,omitempty"`
	PreviewURI  *string    `json:"preview_uri,omitempty"`
	CreatedAt   time.Time  `json:"created_at"`
}
