package dao

import (
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
		SELECT u.id, u.username, u.password, u.first_name, u.last_name, 
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
			&user.Password,
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

func (dao *FriendsDAO) VerifyFriendRequest(requestID uint64, userID uint64) (uint64, error) {
	var receiverID uint64
	err := dao.DB.QueryRow(`
		SELECT receiver_id 
		FROM friend_requests 
		WHERE id = $1 AND sender_id = $2
	`, requestID, userID).Scan(&receiverID)

	if err != nil {
		return 0, sql.ErrNoRows
	}
	return receiverID, nil
}

func (dao *FriendsDAO) AcceptFriendRequest(userID uint64, receiverID uint64) error {
	_, err := dao.DB.Exec("INSERT INTO friends (user_id, friend_id) VALUES ($1, $2)", userID, receiverID)
	return err
}

func (dao *FriendsDAO) DeleteFriendRequest(requestID uint64) error {
	query := `
		DELETE FROM friend_requests 
		WHERE id = $1;
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
