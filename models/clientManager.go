package models

import "fmt"

type ClientManager struct {
	Clients    map[string]*Client
	Register   chan *Client
	Unregister chan *Client
}

func (manager *ClientManager) Start() {
	for {
		select {
		case client := <-manager.Register:
			manager.Clients[client.ID] = client
			fmt.Println("Client connected:", client.ID)
		case client := <- manager.Unregister:
			if _, ok := manager.Clients[client.ID]; ok {
				delete(manager.Clients, client.ID)
				close(client.Send)
				fmt.Println("Client disconnected:", client.ID)
			}
		}

	}
}