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
		SELECT ranking_id, song_id, rank
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
			&ranking.Rank,
		); err != nil{
			return nil, err
		}
		rankings = append(rankings, ranking)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return rankings, nil
}
