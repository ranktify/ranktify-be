package model

type FriendRequests struct {
	RequestID   uint64  `json:"request_id"`
	SenderID    uint64  `json:"sender_id"`
	ReceiverID  uint64  `json:"receiver_id"`
	RequestDate string  `json:"request_date"`
	Status      *string `json:"status"`
}
