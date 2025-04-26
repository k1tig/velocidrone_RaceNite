package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func initBracket(c *gin.Context) {
	var newRecord raceRecord
	newBracketOK := true

	if err := c.BindJSON(&newRecord); err != nil {
		return
	}
	for _, i := range records {
		if i.Id == newRecord.Id {
			c.IndentedJSON(http.StatusOK, gin.H{"message": " error, id already exists"}) // not correct status
			newBracketOK = false
			break
		}
	}
	if newBracketOK {
		records = append(records, newRecord)
		c.IndentedJSON(http.StatusCreated, records)
	}

}

func getBracketById(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusOK, gin.H{"message": " cannot convert recoed id to int"})
		return
	}
	for _, i := range records {
		if i.Id == id {
			c.IndentedJSON(http.StatusOK, i)
			return
		}
	}

}

func getBrackets(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, records)
}

// //////////tests if conn is active/////////
func getWs(c *gin.Context) {
	for i := range clients {
		i.hub.printClients()
	}

	c.IndentedJSON(http.StatusOK, nil)
}

func (h *Hub) printClients() {
	for client := range h.clients {
		log.Println("Client Connected: ", client.conn.RemoteAddr())
	}
}

func editBracket(c *gin.Context) {
	var bracket raceRecord
	if err := c.BindJSON(&bracket); err != nil {
		return
	}
	//log.Printf("Recieved JSON: %v", bracket)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.IndentedJSON(http.StatusOK, gin.H{"message": " cannot convert recoed id to int"})
		return
	}
	for x, b := range records {
		if b.Id == id {
			bracketUpdated := false
			for i, oringalRacer := range b.Pilots {
				for _, editRacer := range bracket.Pilots {
					if editRacer.VdName == oringalRacer.VdName {
						records[x].Pilots[i] = editRacer
						if !bracketUpdated {
							bracketUpdated = true
						}
					}
				}
			}
			if bracketUpdated {
				message, err := json.Marshal(bracket)
				if err != nil {
					log.Printf("Error json") //name
					//fmt.Println(params)
				}

				for client := range clients {
					message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
					//client.send <- message
					//client.Send(websocket.TextMessage, message)
					client.conn.SetWriteDeadline(time.Now().Add(writeWait))

					w, err := client.conn.NextWriter(websocket.TextMessage)
					if err != nil {
						return
					}
					w.Write(message)
					// Add queued chat messages to the current websocket message.
					n := len(client.send)
					for i := 0; i < n; i++ {
						w.Write(newline)
						w.Write(<-client.send)
					}
					if err := w.Close(); err != nil {
						return
					}
				}
				c.IndentedJSON(http.StatusOK, gin.H{"message:": " Update Successfull"})
				return
			}
		} else {
			c.IndentedJSON(http.StatusOK, gin.H{"message": "no update to brackets"})
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"message": "bracket not found"})
}

func serveWs(hub *Hub, c *gin.Context) {
	// Upgrade the HTTP connection to a WebSocket connection
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	log.Println("Client connected:", conn.RemoteAddr())

	client := &Client{hub: hub, conn: conn, send: make(chan []byte)}
	client.hub.register <- client ///////was it this???
	clients[client] = true
	go client.writePump()
	go client.readPump()

}
