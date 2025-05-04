package dao

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ranktify/ranktify-be/internal/model"
)

type FriendsDAO struct {
	DB *sql.DB
}

func NewFriendsDAO(db *sql.DB) *FriendsDAO {
	return &FriendsDAO{DB: db}
}

func (dao *FriendsDAO) GetFriends(id uint64) ([]model.User, error) {
	query := `
		SELECT u.id, u.username, u.first_name, u.last_name, 
		u.email, u.role, u.created_at
		FROM friends f
		JOIN users u ON (f.user_id = $1 AND u.id = f.friend_id)
		OR (f.friend_id = $1 AND u.id = f.user_id)
		WHERE f.user_id = $1 OR f.friend_id = $1;
		`

	rows, err := dao.DB.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friends []model.User
	for rows.Next() {
		var user model.User
		if err := rows.Scan(
			&user.Id,
			&user.Username,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.Role,
			&user.CreatedAt,
		); err != nil {
			return nil, err
		}
		friends = append(friends, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return friends, nil
}

func (dao *FriendsDAO) GetFriendRequests(receiverID uint64) ([]model.FriendRequests, int, error) {
	query := `
		SELECT request_id, sender_id, receiver_id, request_date, status
		FROM friend_requests
		WHERE receiver_id = $1
`
	rows, err := dao.DB.Query(query, receiverID)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var friendRequests []model.FriendRequests
	for rows.Next() {
		var friendRequest model.FriendRequests
		if err := rows.Scan(
			&friendRequest.RequestID,
			&friendRequest.SenderID,
			&friendRequest.ReceiverID,
			&friendRequest.RequestDate,
			&friendRequest.Status,
		); err != nil {
			return nil, 0, err
		}
		friendRequests = append(friendRequests, friendRequest)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	friendCount := len(friendRequests)

	return friendRequests, friendCount, nil
}

func (dao *FriendsDAO) DeleteFriendByID(userID uint64, friendID uint64) error {
	query := `
		DELETE FROM friends 
		WHERE 
		    (user_id = $1 AND friend_id = $2) 
		    OR 
		    (user_id = $2 AND friend_id = $1);
	`
	result, err := dao.DB.Exec(query, userID, friendID)
	if err != nil {
		return fmt.Errorf("error deleting friend: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

func (dao *FriendsDAO) SendFriendRequest(userID uint64, receiverID uint64) error {
	query := `
		INSERT INTO friend_requests (sender_id, receiver_id)
		VALUES ($1, $2);
	`
	_, err := dao.DB.Exec(query, userID, receiverID)
	if err != nil {
		return fmt.Errorf("error sending friend request: %v", err)
	}
	return nil
}

func (dao *FriendsDAO) AcceptFriendRequest(userID uint64, receiverID uint64) error {
	_, err := dao.DB.Exec("INSERT INTO friends (user_id, friend_id) VALUES ($1, $2)", userID, receiverID)
	return err
}

func (dao *FriendsDAO) DeleteFriendRequest(requestID uint64) error {
	query := `
		DELETE FROM friend_requests 
		WHERE request_id = $1;
	`
	result, err := dao.DB.Exec(query, requestID)
	if err != nil {
		return fmt.Errorf("error deleting friend request: %v", err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected: %v", err)
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}
	return nil
}

type topSongsStruct struct {
	model.Song
	AvgRank      float64 `json:"avg_rank"`
	RatingsCount int     `json:"rating_count"`
}

func (dao *FriendsDAO) GetTopNTracksAmongFriends(ctx context.Context, userID uint64, limit int) ([]topSongsStruct, error) {
	query := `
		SELECT
		s.song_id,
		s.title,
		s.artist,
		cover_uri,
		preview_uri,
		AVG(r.rank) AS avg_rank,
		COUNT(*)    AS ratings_count
		FROM friends f
		JOIN rankings r
		ON (
				(f.user_id   = $1 AND f.friend_id = r.user_id)
			OR (f.friend_id = $1 AND f.user_id   = r.user_id)
			)
		JOIN songs s
		ON s.song_id = r.song_id
		GROUP BY
		s.song_id,
		s.title,
		s.artist
		ORDER BY
		avg_rank      DESC,
		ratings_count DESC
		LIMIT $2;
    `

	rows, err := dao.DB.QueryContext(ctx, query, userID, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var topSongs []topSongsStruct
	for rows.Next() {

		var topSong topSongsStruct
		if err := rows.Scan(
			&topSong.SongID,
			&topSong.Title,
			&topSong.Artist,
			&topSong.CoverURI,
			&topSong.PreviewURI,
			&topSong.AvgRank,
			&topSong.RatingsCount,
		); err != nil {
			return nil, err
		}
		topSongs = append(topSongs, topSong)
	}
	return topSongs, nil
}
