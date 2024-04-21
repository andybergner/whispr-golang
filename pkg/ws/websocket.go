package ws

import (
	"net/http"
	"whispr-golang/helpers"
	"whispr-golang/models"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleConnections(w http.ResponseWriter, r *http.Request, manager *models.ClientManager) {
	
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

	client := &models.Client{ID: claims.Uid, Conn: conn, Send: make(chan []byte)}
	manager.Register <- client

	go client.ReadPump(manager)
	go client.WritePump()
}
