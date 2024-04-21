package models

type Message struct {
	RecipientID string `json:"recipientid"`
	Content     string `json:"content"`
}