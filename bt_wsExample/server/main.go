// server.go
package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	hub := newHub()
	go hub.run()
	router := mux.NewRouter()
	router.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})

	log.Println("Server started on :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

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

var records []raceRecord
