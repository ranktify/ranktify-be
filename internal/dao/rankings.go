package dao

import (
	"context"
	"database/sql"
	"time"

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
			s.album,  -- Select other song details you need
			s.genre,  -- Select other song details you need
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
