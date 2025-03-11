package main

import "fmt"

type Styles struct{}

// Pilot stuff
type Pilot struct {
	DiscordName    string `csv:"Display Name" json:"displayname"` // discord
	VdName         string `csv:"Player Name" json:"vdname"`
	QualifyingTime string `csv:"Lap Time" json:"qualifytime"`
	ModelName      string `csv:"Model Name" json:"modelname"`
	Id             string `csv:"ID" json:"id"`
	Status         string //used for checkin placeholder
}

var NulPilot Pilot

type listRacer struct {
	name, time, craft string
}

func (i listRacer) Title() string { return i.name }
func (i listRacer) Description() string {
	description := fmt.Sprintf("%s | %s", i.time, i.craft)
	return description
}
func (i listRacer) FilterValue() string { return i.name }
