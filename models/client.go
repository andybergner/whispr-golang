package models

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID   string
	Conn *websocket.Conn
	Send chan []byte
}

func (client *Client) ReadPump(manager *ClientManager) {
	defer func() {
		manager.Unregister <- client
		client.Conn.Close()
	}()
	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			manager.Unregister <- client
			client.Conn.Close()
			break
		}

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			fmt.Println("Error decoding message", err.Error())
			continue
		}

		msg.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		msg.Sender_id = client.ID

		msgBytes, _ := json.Marshal(msg)

		if msg.Recipient_id != "" {
			recipientClient, ok := manager.Clients[msg.Recipient_id]
			if !ok {
				fmt.Println("recipient client is not connected")
				continue
			}
			select {
			case recipientClient.Send <- msgBytes:
			default:
				close(recipientClient.Send)
				delete(manager.Clients, recipientClient.ID)
			}
		}
	}
}

func (client *Client) WritePump() {
	defer func() {
		client.Conn.Close()
	}()
	for message := range client.Send {
		if err := client.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
			return
		}
	}
}
