package model

type FriendRequests struct {
	ID          uint64 `json:"id"`
	SenderID    uint64 `json:"sender_id"`
	ReceiverID  uint64 `json:"receiver_id"`
	RequestDate string `json:"request_date"`
	Status      string `json:"status"`
}
