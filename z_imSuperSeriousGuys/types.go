package main

import "fmt"

var fmvTag string = `



_____ __  ____     __                     
|  ___|  \/  \ \   / /                     
| |_  | |\/| |\ \ / /                      
|  _| | |  | | \ V /                       
|_|__ |_|  |_|  \_/  _   _ _ _       _ _ _ 
|   _ \ __ _ ___ ___| \ | (_) |_ ___| | | |
| |_)  / _ |/ __/ _ \  \| | | __/ _ \ | | |
|  _ < (_| | (_|  __/ |\  | | ||  __/_|_|_|
|_| \_\__,_|\___\___|_| \_|_|\__\___(_|_|_)

   
   `

type Styles struct{}

// Pilot stuff
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

type listRacer struct {
	name, time, craft string
}

// for list interface
func (i listRacer) Title() string { return i.name }
func (i listRacer) Description() string {
	description := fmt.Sprintf("%s | %s", i.time, i.craft)
	return description
}
func (i listRacer) FilterValue() string { return i.name }
