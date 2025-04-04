package main

import (
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
	router.GET("/list", getWs) // for testing live wsconnects

	router.Run("localhost:8080")
}
