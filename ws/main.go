package main

import (
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type Client struct {
}

func main() {
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial("ws://192.168.68.83:60003/velocidrone", nil)
	if err != nil {
		fmt.Errorf("Error: %s", err)
	}
	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		fmt.Printf("Received: %s\n", message)
	}

}
