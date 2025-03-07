package main

type Pilot struct {
	DiscordName    string `csv:"Display Name" json:"displayname"` // discord
	VdName         string `csv:"Player Name" json:"vdname"`
	QualifyingTime string `csv:"Lap Time" json:"qualifytime"`
	ModelName      string `csv:"Model Name" json:"modelname"`
}
