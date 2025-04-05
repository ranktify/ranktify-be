package dao

import (
	"database/sql"
	"fmt"
)

type FriendsDAO struct {
	DB *sql.DB
}

func NewFriendsDAO(db *sql.DB) *FriendsDAO {
	return &FriendsDAO{DB: db}
}

func (dao *FriendsDAO) GetFriends(id uint64) ([]string, error) {
	query := `
		SELECT u.username 
		FROM users u
		JOIN friends f ON u.id = f.friend_id
		WHERE f.user_id = $1

		UNION 

		SELECT u.username 
		FROM users u
		JOIN friends f ON u.id = f.user_id
		WHERE f.friend_id = $1;
	`

	rows, err := dao.DB.Query(query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friends []string
	for rows.Next() {
		var username string
		if err := rows.Scan(&username); err != nil {
			return nil, err
		}
		friends = append(friends, username)
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

func (dao *FriendsDAO) ProcessFriendRequest(userID uint64, requestID uint64, action string) error {
	// Check if the friend request exists and belongs to the user
	var receiverID uint64
	err := dao.DB.QueryRow("SELECT receiver_id FROM friend_requests WHERE id = $1 AND sender_id = $2", requestID, userID).Scan(&receiverID)
	if err != nil {
		return sql.ErrNoRows // Request not found
	}
	// If the user accepts the request
	if action == "accept" {
		// Begin transaction
		tx, err := dao.DB.Begin()
		if err != nil {
			return err
		}
		// Insert into friends table
		_, err = tx.Exec("INSERT INTO friends (user_id, friend_id) VALUES ($1, $2)", userID, receiverID)
		if err != nil {
			tx.Rollback()
			return err
		}
		// Delete friend request
		_, err = tx.Exec("DELETE FROM friend_requests WHERE id = $1", requestID)
		if err != nil {
			tx.Rollback()
			return err
		}
		// Commit transaction
		return tx.Commit()
	}
	// If the user declines the request, just delete it
	_, err = dao.DB.Exec("DELETE FROM friend_requests WHERE id = $1", requestID)
	return err
}

func (dao *FriendsDAO) CancelFriendRequest(requestID uint64) error {
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
