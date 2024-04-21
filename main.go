package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"whispr-golang/helpers"
	"whispr-golang/routes"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	id		 string
	conn 	 *websocket.Conn
	send 	 chan []byte
}

type Message struct {
	RecipientID 	string `json:"targetid"`
	Content 	string `json:"content"`
}

type ClientManager struct {
	clients		map[string]*Client
	register	chan *Client
	unregister	chan *Client
}

func (manager *ClientManager) start() {
	for {
		select {
		case client := <-manager.register:
			manager.clients[client.id] = client
			fmt.Println("Client connected:", client.id)
		case client := <- manager.unregister:
			if _, ok := manager.clients[client.id]; ok {

				delete(manager.clients, client.id)
				close(client.send)
				fmt.Println("Client disconnected:", client.id)
			}
		}
	}
}

func handleConnections(w http.ResponseWriter, r *http.Request, manager *ClientManager) {

	clientToken := r.Header.Get("token")
	if clientToken == "" {
		http.Error(w, "No Authorization header provided", http.StatusUnauthorized)
		return
	}

	claims, validateErr := helpers.ValidateToken(clientToken)
	if validateErr != "" {
		http.Error(w, validateErr, http.StatusInternalServerError)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Could not open WebSocket connection", http.StatusBadRequest)
		return
	}

	client := &Client{id: claims.Uid, conn: conn, send: make(chan []byte)}
	manager.register <- client

	go client.readPump(manager)
	go client.writePump()

}

func (client *Client) readPump(manager *ClientManager) {
	defer func() {
		manager.unregister <- client
		client.conn.Close()
	}()
	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			manager.unregister <- client
			client.conn.Close()
			break
		}

		

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			fmt.Println("Error decoding message", err.Error())
			continue
		}


		if msg.RecipientID != "" {
			recipientClient, ok := manager.clients[msg.RecipientID]
			if !ok {
				fmt.Println("recipient client is not connected")
				continue
			}
			select {
			case recipientClient.send <- message:
			default:
				close(recipientClient.send)
				delete(manager.clients, recipientClient.id)
			}
		} 
	}
}

func (client *Client) writePump() {
	defer func() {
		client.conn.Close()
	}()
	for message := range client.send {
		if err := client.conn.WriteMessage(websocket.TextMessage, message); err != nil {
			return
		}
	}
}

func main() {

	router := gin.Default()
	router.Use(gin.Logger())

	manager := ClientManager{
		clients: make(map[string]*Client),
		register: make(chan *Client),
		unregister: make(chan *Client),
	}

	go manager.start()

	router.GET("/ws", func(ctx *gin.Context) {
		handleConnections(ctx.Writer, ctx.Request, &manager)
	})

	routes.AuthRouter(router)
	routes.UserRouter(router)
	//routes.ChatRouter(router)

	router.Run(":8080")

}