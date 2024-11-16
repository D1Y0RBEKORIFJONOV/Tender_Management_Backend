package entity

import "time"

type (
	MessageBid struct {
		CreateAt   time.Time `json:"create_at" bson:"create_at"`
		SenderName string    `json:"sender_name" bson:"sender_name"`
		Status     string    `json:"status" bson:"status"`
	}
	CreateNotification struct {
		UserId string `json:"user_id" bson:"user_id"`
	}
	Notification struct {
		UserId   string       `json:"user_id" bson:"user_id"`
		Offset   int64        `json:"offset" bson:"offset"`
		Messages []MessageBid `json:"messages" bson:"messages"`
	}
	GetNotificationReq struct {
		UserId string `json:"user_id" bson:"user_id"`
		Offset int64  `json:"offset" bson:"offset"`
	}
	GetNotificationResp struct {
		Messages []MessageBid `json:"notifications" bson:"notifications"`
	}
	CreateMessageReq struct {
		Status     string `json:"status" bson:"status"`
		SenderName string `json:"sender_name" bson:"sender_name"`
	}
	AddNotificationReq struct {
		UserId        string            `json:"user_id" bson:"user_id"`
		CreateMessage *CreateMessageReq `json:"notification" bson:"notification"`
	}
)
