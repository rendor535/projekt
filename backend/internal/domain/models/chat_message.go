package models

import "time"

type ChatMessage struct {
	ID         int       `json:"id" db:"id"`
	SenderID   int       `json:"sender_id" db:"sender_id"`
	ReceiverID int       `json:"receiver_id" db:"receiver_id"`
	Content    string    `json:"content" db:"content"`
	IsRead     bool      `json:"is_read" db:"is_read"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
}

type ChatMessageWithDetails struct {
	ChatMessage
	SenderUsername   string `json:"sender_username"`
	ReceiverUsername string `json:"receiver_username"`
}
