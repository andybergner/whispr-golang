package models

import "time"

type Message struct {
	Recipient_id string    `json:"recipient_id"`
	Sender_id    string    `json:"sender_id"`
	Content      string    `json:"content"`
	Created_at   time.Time `json:"created_at"`
}
