package model

type Friends struct {
	UserId 	       uint64 `json:"user_id"`
	FriendId  	   uint64 `json:"friend_id"`
	FriendshipDate string `json:"friendship_date"`
}