package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Message struct {
	ID         string    `gorm:"id;primaryKey"`
	SenderID   string    `gorm:"sender_id"`
	ReceiverID string    `gorm:"receiver_id"`
	Content    string    `gorm:"content"`
	IsRead     bool      `gorm:"is_read"`
	CreatedAt  time.Time `gorm:"created_at"`
	UpdatedAt  time.Time `gorm:"updated_at"`
	DeletedAt  time.Time `gorm:"deleted_at"`
}

type MessageReq struct {
	SenderID   string `json:"sender_id"`
	ReceiverID string `json:"receiver_id"`
	Content    string `json:"content"`
	IsRead     bool   `json:"is_read"`
}

type MessageResp struct {
	ID         string    `json:"id"`
	SenderID   string    `json:"sender_id"`
	ReceiverID string    `json:"receiver_id"`
	Content    string    `json:"content"`
	IsRead     bool      `json:"is_read"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (mReq *MessageReq) ToMessage() *Message {
	return &Message{
		ID:         primitive.NewObjectID().Hex(),
		SenderID:   mReq.SenderID,
		ReceiverID: mReq.ReceiverID,
		Content:    mReq.Content,
		IsRead:     mReq.IsRead,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
}

func (msg *Message) ToMessageResp() *MessageResp {
	return &MessageResp{
		ID:         primitive.NewObjectID().Hex(),
		SenderID:   msg.SenderID,
		ReceiverID: msg.ReceiverID,
		Content:    msg.Content,
		IsRead:     msg.IsRead,
		CreatedAt:  msg.CreatedAt,
		UpdatedAt:  msg.UpdatedAt,
	}
}
