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

func (dao *RankingsDao) CheckIfSongIsRanked(spotifyId string, userID uint64) (bool, error) {
	var exists bool
	err := dao.DB.QueryRow(`
        SELECT EXISTS (
            SELECT 1
              FROM rankings r
              JOIN songs    s ON s.song_id = r.song_id
             WHERE s.spotify_id = $1
               AND r.user_id   = $2
        );
    `, spotifyId, userID).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (dao *RankingsDao) GetFriendsRankedSongsWithNoUserRank(userID uint64) ([]map[string]interface{}, error) {
	query := `
		SELECT DISTINCT
			r.song_id,
			r.ranking_id,
			r.user_id,
			r.rank,
			s.spotify_id,
			s.title,
			s.artist,
			s.album,
			s.release_date,
			s.genre,
			s.cover_uri,
			s.preview_uri,
			s.created_at
		FROM friends f
		JOIN users u ON
			(f.user_id = $1 AND u.id = f.friend_id)
			OR (f.friend_id = $1 AND u.id = f.user_id)
		JOIN rankings r ON r.user_id = u.id
		JOIN songs s ON s.song_id = r.song_id
		WHERE r.song_id NOT IN (
			SELECT song_id FROM rankings WHERE user_id = $1
		)
		LIMIT 5;
	`
	rows, err := dao.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	for rows.Next() {
		// Create a slice of interface{}'s to hold each column value
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))

		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		// Create map for this row
		rowMap := make(map[string]interface{})
		for i, col := range columns {
			rowMap[col] = values[i]
		}

		results = append(results, rowMap)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

