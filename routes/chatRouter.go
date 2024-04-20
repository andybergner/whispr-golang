package routes

import (
	"whispr-golang/controllers"

	"github.com/gin-gonic/gin"
)


func ChatRouter(incomingRoutes *gin.Engine) {
	incomingRoutes.GET("/ws", controllers.WebsocketHandler())
}
