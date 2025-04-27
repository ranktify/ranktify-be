package dao

import (
	"database/sql"

	"github.com/ranktify/ranktify-be/internal/model"
)

type RankingsDao struct {
	DB *sql.DB
}

func NewRankingsDAO(db *sql.DB) *RankingsDao {
	return &RankingsDao{DB: db}
}

func (dao *RankingsDao) GetRankedSongs(userID uint64) ([]model.Rankings, error) {
	query := `
		SELECT ranking_id, song_id, user_id, rank
		FROM rankings
		WHERE user_id = $1
	`
	rows, err := dao.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rankings []model.Rankings
	for rows.Next() {
		var ranking model.Rankings
		if err := rows.Scan(
			&ranking.RankingID,
			&ranking.SongID,
			&ranking.UserID,
			&ranking.Rank,
		); err != nil {
			return nil, err
		}
		rankings = append(rankings, ranking)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return rankings, nil
}

func (dao *RankingsDao) GetFriendsRankedSongs(userID uint64) ([]model.Rankings, error) {
	query := `
		SELECT r.ranking_id, r.song_id, r.user_id, r.rank
		FROM friends f
		JOIN users u ON (f.user_id = $1 AND u.id = f.friend_id)
					OR (f.friend_id = $1 AND u.id = f.user_id)
		JOIN rankings r ON r.user_id = u.id
		WHERE f.user_id = $1 OR f.friend_id = $1;

	`
	rows, err := dao.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rankings []model.Rankings
	for rows.Next() {
		var ranking model.Rankings
		if err := rows.Scan(
			&ranking.RankingID,
			&ranking.SongID,
			&ranking.UserID,
			&ranking.Rank,
		); err != nil {
			return nil, err
		}
		rankings = append(rankings, ranking)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return rankings, nil
}

func (dao *RankingsDao) CheckIfSongIsRanked(userID uint64, songID uint64) (bool, error) {
	var exists bool
	err := dao.DB.QueryRow(`
		SELECT EXISTS (
			SELECT 1 FROM rankings WHERE user_id = $1 AND song_id = $2
		)
	`, userID, songID).Scan(&exists)

	if err != nil {
		return false, err
	}

	return exists, nil
}
