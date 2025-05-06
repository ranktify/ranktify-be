package dao

import (
	"context"
	"database/sql"
	"time"
	"fmt"

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

func (dao *RankingsDao) GetTopWeeklyRankedSongs(ctx context.Context) ([]model.Song, error) {
	query := `
		-- get all the rank songs and their average rank and rating 
		WITH RankedSongs AS (
			SELECT
				song_id,
				AVG(rank) AS avg_rank,
				COUNT(*)  AS ratings_count
			FROM rankings
			WHERE 
				updated_at >= $1 AND updated_at < $2
			GROUP BY song_id
			ORDER BY
				avg_rank DESC,      -- higher average rank is better
				ratings_count DESC -- Higher rating count breaks ties (more popular)
			LIMIT 5
		)
		SELECT
			s.song_id,
			s.spotify_id,
			s.title,
			s.artist,
			s.album,  
			s.release_date,
			s.genre,  
			s.cover_uri,
			s.preview_uri
		FROM songs s
		JOIN RankedSongs rs ON s.song_id = rs.song_id
		ORDER BY
			rs.avg_rank DESC,      
			rs.ratings_count DESC;
	`
	now := time.Now()
	// today at midnight
	todayMid := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

	// compute how many days since Monday:
	// (Weekday()+6)%7 maps Monday→0, Tuesday→1, … Sunday→6
	offset := (int(now.Weekday()) + 6) % 7

	// this week’s Monday 00:00
	startThisWeek := todayMid.AddDate(0, 0, -offset)

	// last week’s Monday 00:00
	startLastWeek := startThisWeek.AddDate(0, 0, -7)

	rows, err := dao.DB.QueryContext(ctx, query, startLastWeek, startThisWeek)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var songs []model.Song
	for rows.Next() {
		var song model.Song
		if err := rows.Scan(
			&song.SongID,
			&song.SpotifyID,
			&song.Title,
			&song.Artist,
			&song.Album,
			&song.ReleaseDate,
			&song.Genre,
			&song.CoverURI,
			&song.PreviewURI,
		); err != nil {
			return nil, err
		}
		songs = append(songs, song)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return songs, nil
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

func (dao *RankingsDao) RankSong(songID uint64, userID uint64, rank int) error {
	query := `
		INSERT INTO rankings (song_id, user_id, rank, created_at, updated_at)
		VALUES ($1, $2, $3, NOW(), NOW())
	`
	_, err := dao.DB.Exec(query, songID, userID, rank)
	if err != nil {
		return fmt.Errorf("error ranking song: %v", err)
	}
	return nil
}

func (dao *RankingsDao) DeleteRanking(rankingID uint64) error {
	query := `
		DELETE FROM rankings 
		WHERE ranking_id = $1;
	`
	result, err := dao.DB.Exec(query, rankingID)
	if err != nil {
		return fmt.Errorf("error deleting ranking: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected (Rankings): %v", err)
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (dao *RankingsDao) UpdateRanking(rankingID uint64, rank int) error {
	query := `
		UPDATE rankings
		SET rank = $2
		WHERE ranking_id = $1;
	`
	_, err := dao.DB.Exec(query, rankingID, rank)
	return err
}

func (dao *RankingsDao) StoreSongInDB(spotifyID string, title string, artist *string, album *string,
	releaseDate *time.Time, genre *string, coverURI *string, previewURI *string) error {
	query := `
		INSERT INTO songs (spotify_id, title, artist, album, release_date,
			genre, cover_uri, preview_uri, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW())
	`
	_, err := dao.DB.Exec(query, spotifyID, title, artist, album, releaseDate, genre, coverURI, previewURI)
	if err != nil {
		return fmt.Errorf("error storing song: %v", err)
	}
	return nil
}

func (dao *RankingsDao) GetSongBySpotifyID(spotifyID string) (model.Song, error) {
    var song model.Song
    const query = `
        SELECT
            song_id,
            spotify_id,
            title,
            artist,
            album,
            release_date,
            genre,
            cover_uri,
            preview_uri,
            created_at
        FROM songs
        WHERE spotify_id = $1
    `

    err := dao.DB.QueryRow(query, spotifyID).Scan(
        &song.SongID,
        &song.SpotifyID,
        &song.Title,
        &song.Artist,
        &song.Album,
        &song.ReleaseDate,
        &song.Genre,
        &song.CoverURI,
        &song.PreviewURI,
        &song.CreatedAt,
    )
    if err != nil {
        return model.Song{}, fmt.Errorf("GetSongBySpotifyID: %w", err)
    }

    return song, nil
}

