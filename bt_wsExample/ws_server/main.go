// server.go
package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins
	},
}

func handleConnections(c *gin.Context) {
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()
	log.Println("Client connected:", ws.RemoteAddr())

	for {
		messageType, p, err := ws.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		log.Printf("Received message from %s: %s\n", ws.RemoteAddr(), p)

		err = ws.WriteMessage(messageType, p)
		if err != nil {
			log.Println(err)
			return
		}
	}

}

func LoggingMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start time
		start := time.Now()

		// Process the request
		c.Next()

		// Stop time
		end := time.Now()

		// Log the request details
		log.Printf("%s %s %s %s %v %s", c.ClientIP(), c.Request.Method, c.Request.Method, c.Request.URL.Path, c.Writer.Status(), end.Sub(start))
	}
}

func main() {
	router := gin.Default()
	router.Use(LoggingMiddleware())
	router.GET("/ws", func(c *gin.Context) {
		handleConnections(c)
	})
	router.Run("localhost:8080")
}
