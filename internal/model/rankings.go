package model

type Rankings struct {
	RankingID uint64 `json:"ranking_id"`
	SongID    uint64 `json:"song_id"`
	UserID    uint64 `json:"user_id"`
	Rank      int    `json:"rank"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}
