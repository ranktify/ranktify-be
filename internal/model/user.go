package model

import "time"

type User struct {
	Id                       uint64    `json:"id"`
	Username                 string    `json:"username"`
	Password                 string    `json:"password"`
	FirstName                *string   `json:"first_name,omitempty"`
	LastName                 *string   `json:"last_name,omitempty"`
	Email                    string    `json:"email"`
	Role                     *string   `json:"role,omitempty"`
	SpotifyID                *string   `json:"spotify_id,omitempty"`
	SpotifyDisplayName       *string   `json:"spotify_display_name,omitempty"`
	SpotifyProfileURI        *string   `json:"spotify_profile_uri,omitempty"`
	SpotifyProfilePictureURI *string   `json:"spotify_profile_picture_uri,omitempty"`
	CreatedAt                time.Time `json:"created_at"`
}
