package main

// Pilot stuff
type Pilot struct {
	DiscordName    string `csv:"Display Name" json:"displayname"` // discord
	VdName         string `csv:"Player Name" json:"vdname"`
	QualifyingTime string `csv:"Lap Time" json:"qualifytime"`
	ModelName      string `csv:"Model Name" json:"modelname"`
	Id             string `csv:"ID" json:"id"`
	Status         string //used for checkin placeholder
}

var NulPilot = Pilot{Status: "-"}

type Styles struct{}
