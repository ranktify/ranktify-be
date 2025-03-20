package model

import "time"

type JWTRefreshToken struct {
	ID           uint64 // Primary Key
	UserID       uint64 // References users(id)
	JTI          string // Unique Token Identifier
	RefreshToken string
	ExpiresAt    time.Time
	CreatedAt    time.Time
}

type SpotifyRefreshToken struct {
	ID        uint64    `json:"id"`
	UserID    uint64    `json:"user_id"`
	Token     string    `json:"token"`
	ExpiresAt time.Time `json:"expires_at"`
}
