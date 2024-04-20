package routes

import (
	"whispr-golang/controllers"
	"whispr-golang/middleware"

	"github.com/gin-gonic/gin"
)

func UserRouter(incomingRoutes *gin.Engine) {
	incomingRoutes.Use(middleware.Authenticate())
	incomingRoutes.GET("/users/:user_id", controllers.GetUser())
}