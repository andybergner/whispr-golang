package controllers

import (
	"whispr-golang/models"
	"whispr-golang/pkg/ws"

	"github.com/gin-gonic/gin"
)

func WebsocketHandler() gin.HandlerFunc {
	manager := models.ClientManager{
		Clients:    make(map[string]*models.Client),
		Register:   make(chan *models.Client),
		Unregister: make(chan *models.Client),
	}

	go manager.Start()

	return func(ctx *gin.Context) {
		ws.HandleConnections(ctx.Writer, ctx.Request, &manager)
	}
}
