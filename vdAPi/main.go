package main

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// add grey group for catch-ups or late qualify

type Pilot struct {
	DiscordName    string      `csv:"Display Name" json:"displayname"` // discord
	VdName         string      `csv:"Player Name" json:"vdname"`
	QualifyingTime string      `csv:"Lap Time" json:"qualifytime"`
	ModelName      string      `csv:"Model Name" json:"modelname"`
	Id             string      `csv:"ID" json:"id"`
	Status         bool        `json:"status"` //used for checkin placeholder
	Points         float64     `json:"points"`
	RaceTimes      [10]float64 `json:"racetimes"`
}

type raceRecord struct {
	Id         int     `json:"id"`
	RoomPhrase string  `json:"roomphrase"`
	Round      int     `json:"round"`
	Turn       int     `json:"turn"`
	Pilots     []Pilot `json:"pilots"`
}

/*type bracketStatus struct {
	rev int
}*/

var records []raceRecord

func main() {
	hub := newHub()
	go hub.run()
	router := gin.Default()
	router.GET("/ws", func(c *gin.Context) {
		serveWs(hub, c)
	})
	router.GET("/brackets", getBrackets)
	router.GET("/brackets/:id", getBracketById)
	router.POST("/brackets", initBracket)
	router.PUT("/brackets/:id", editBracket)

	router.Run("localhost:8080")
}

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

func editBracket(c *gin.Context) {
	var bracket raceRecord
	if err := c.BindJSON(&bracket); err != nil {
		return
	}
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
				c.IndentedJSON(http.StatusOK, records)
			} else {
				c.IndentedJSON(http.StatusOK, gin.H{"message": "no update to brackets"})
			}
			break
		}
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "bracket not found"})
	}

}
