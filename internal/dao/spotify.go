package dao

import (
	"database/sql"

	"github.com/ranktify/ranktify-be/internal/model"
)

type SpotifyDAO struct {
	DB *sql.DB
}

func NewSpotifyDAO(db *sql.DB) *SpotifyDAO {
	return &SpotifyDAO{DB: db}
}

func (dao *SpotifyDAO) SaveRefreshToken(rt model.SpotifyRefreshToken) error {
	query := `
		INSERT INTO spotify_refresh_tokens (user_id, refresh_token, created_at)
		VALUES ($1, $2, NOW())
	`
	_, err := dao.DB.Exec(query, rt.UserID, rt.Token)
	return err
}

func (dao *SpotifyDAO) GetRefreshToken(userID string) (string, error) {
	var token string
	query := `
		SELECT refresh_token 
		FROM spotify_refresh_tokens
		WHERE user_id = $1
		LIMIT 1
	`

	err := dao.DB.QueryRow(query, userID).Scan(&token)
	if err != nil {
		return "", err
	}
	return token, nil
}
