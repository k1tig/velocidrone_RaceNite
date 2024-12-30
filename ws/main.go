package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

/*
type racer struct {
	Finished string `json:"finished"`
	Gate     string `json:"gate"`
	Lap      string `json:"lap"`
	Position string `json:"position"`
	Time     string `json:"time"`
	Colour   string `json:"colour"`
}*/

var data map[string]interface{}

func main() {
	dialer := websocket.Dialer{}
	conn, _, err := dialer.Dial("ws://192.168.68.83:60003/velocidrone", nil)
	if err != nil {
		log.Panic(err)
	}
	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}

		if err := json.Unmarshal([]byte(message), &data); err != nil {
			log.Fatal(err)

		}
		fmt.Println(data["racestatus"])
	}
}
