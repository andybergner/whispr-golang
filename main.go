package main

import (
	"fmt"
	"net/http"
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

func main() {

	router := gin.Default()
	router.Use(gin.Logger())

	router.GET("/ws", func(ctx *gin.Context) {
		conns := make(map[*websocket.Conn]bool)
		conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			fmt.Println("upgrader error:", err.Error())
			return
		}
		defer conn.Close()
	
		fmt.Println("new incoming connection from client:", conn.RemoteAddr())
	
		conns[conn] = true
	
		for {
			_, msg, err := conn.ReadMessage()
			if err != nil {
				fmt.Println("read error:", err.Error())
				delete(conns, conn)
				return
			}
	
			fmt.Println(string(msg))
			if err := conn.WriteMessage(websocket.TextMessage, []byte("thank you for the msg!!!")); err != nil {
				fmt.Println("write error:", err.Error())
			}
		}
	})

	routes.AuthRouter(router)
	routes.UserRouter(router)
	//routes.ChatRouter(router)

	router.Run(":8080")

}