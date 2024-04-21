package main

import (
	"whispr-golang/routes"

	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.New()
	router.Use(gin.Logger())

	routes.AuthRouter(router)
	routes.UserRouter(router)
	routes.ChatRouter(router)

	router.Run(":8080")

}
