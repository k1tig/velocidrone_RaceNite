package main

import (
	"fmt"
	"os"

	"github.com/gocarina/gocsv"
)

func GetVdRacers(filename string) []Pilot {
	var pilot = []Pilot{}
	var filteredPilots = []Pilot{}

	vdTrackFile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer vdTrackFile.Close()
	if err := gocsv.UnmarshalFile(vdTrackFile, &pilot); err != nil { // Load clients from file
		panic(err)
	}
	//use OkRaceClass as a landing spot to filter Velocidrone list by spec
	filteredPilots = append(filteredPilots, pilot...)
	if _, err := vdTrackFile.Seek(0, 0); err != nil { // Go to the start of the file
		panic(err)
	}
	return filteredPilots
}

func GetFMVvoice(fileCsv string) []Pilot {
	var FmvRacers = []Pilot{}
	fmvVoiceFile, err := os.OpenFile(fileCsv, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer fmvVoiceFile.Close()
	if err := gocsv.UnmarshalFile(fmvVoiceFile, &FmvRacers); err != nil { // Load clients from file
		fmt.Printf("Something broke with FMV CSV: %v", err) //csv needs to be in same folder as main.go for now
	}
	if _, err := fmvVoiceFile.Seek(0, 0); err != nil { // Go to the start of the file
		panic(err)
	}
	return FmvRacers
}

func GetDiscordId(discordCsv string) []Pilot {
	var DiscordRecords = []Pilot{}
	discordIdFile, err := os.OpenFile(discordCsv, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer discordIdFile.Close()
	if err := gocsv.UnmarshalFile(discordIdFile, &DiscordRecords); err != nil { // Load clients from file
		fmt.Printf("Something broke with discord ID CSV: %v", err) //csv needs to be in same folder as main.go for now
	}
	if _, err := discordIdFile.Seek(0, 0); err != nil { // Go to the start of the file
		panic(err)
	}
	return DiscordRecords
}

func BindLists(velocidroneList, fmvList, discordCheatSheet []Pilot) []Pilot {
	var NulPilot Pilot
	//discordCheatSheet is the DIY list a group can keep of discord usr IDs and VD names in CSV
	//might break this with the index usage
	for fmvIndex, fmvPilot := range fmvList {
		for _, discord := range discordCheatSheet {
			if discord.Id == fmvPilot.Id {
				fmvList[fmvIndex].VdName = discord.VdName

			}
		}
	}
	for index, fmvPilot := range fmvList {
		for _, velocidronePilot := range velocidroneList {
			if velocidronePilot.VdName == fmvPilot.VdName || fmvPilot.DiscordName == velocidronePilot.VdName {
				fmvPilot.QualifyingTime = velocidronePilot.QualifyingTime
				fmvPilot.ModelName = velocidronePilot.ModelName
				fmvPilot.VdName = velocidronePilot.VdName
				fmvList[index] = fmvPilot
				break
			}
			if fmvPilot.QualifyingTime == NulPilot.QualifyingTime {

				fmvList[index].QualifyingTime = "CHECK IN Please!"
			}
		}
	}
	return fmvList
}
