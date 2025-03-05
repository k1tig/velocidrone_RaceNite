package racetools

import (
	"fmt"
	"os"

	"github.com/gocarina/gocsv"
)

type VdPilot struct { //struct to recieve data from velocidrone csv
	VelocidronName string `csv:"Player Name"`
	QualifyingTime string `csv:"Lap Time"`
	X_Pos          string `csv:"-"`
	ModelName      string `csv:"Model Name"`
	X_Country      string `csv:"-"`
}

type FmvVoicePilot struct {
	RacerName      string `csv:"Display Name"`
	VdName         string `csv:"VdName"`
	QualifyingTime string
	ModelName      string
	Id             string `csv:"ID"`
	Status         string // for checkin placeholder
}

type DiscordIds struct {
	DiscordId string `csv:"DiscordId"`
	VdName    string `csv:"VdName"`
}

// take a list of racers and returns group sets of racers

func GetVdRacers(filename string) []*VdPilot {

	var Clients = []*VdPilot{}
	var OkRaceClass = []*VdPilot{}

	raceFile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer raceFile.Close()
	if err := gocsv.UnmarshalFile(raceFile, &Clients); err != nil { // Load clients from file
		panic(err)
	}
	//use OkRaceClass as a landing spot to filter Velocidrone list by spec
	OkRaceClass = append(OkRaceClass, Clients...)

	if _, err := raceFile.Seek(0, 0); err != nil { // Go to the start of the file
		panic(err)
	}
	return OkRaceClass
}

// For using the voice chat in FMV discord as base group for pairing.
func GetFMVvoice(fileCsv string) []*FmvVoicePilot {

	var FmvRacers = []*FmvVoicePilot{}

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

func GetDiscordId(discordCsv string) []*DiscordIds {
	var DiscordRecords = []*DiscordIds{}

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

func RaceArray(vdList [][]string) [][][]string {
	//makes a group of groups with the total amount of racers not exceeding a +1 differential
	var maxGroupsize = 8
	var grouplength int
	var totalGroups int
	var modulus int

	racers := len(vdList)
	if racers > 40 {
		maxGroupsize = 10
	}

	for i := 1; i <= maxGroupsize; i++ {
		if racers/i <= maxGroupsize {
			totalGroups = i
			modulus = racers % i
			if modulus == 0 {
				grouplength = racers / i
			} else {
				grouplength = (racers - modulus) / i
			}
			break
		}
	}

	var groupStructure = make([][][]string, totalGroups)
	var c int
	x := modulus

	for i := 1; i <= totalGroups; i++ {
		if x > 0 { // distribues the modulus between the lower teir groups
			racers := vdList[c : i*(grouplength+1)]
			groupStructure[i-1] = racers
			x--
			c += grouplength + 1
		} else { // groups that don't take a modulus
			racers := vdList[c : c+grouplength]
			groupStructure[i-1] = racers
			c += grouplength
		}
	}
	return groupStructure
}

func BindLists(vdl []*VdPilot, fmvl []*FmvVoicePilot, dcl []*DiscordIds) []*FmvVoicePilot {
	//var bound []*FmvVoicePilot

	for _, f := range fmvl {
		for _, d := range dcl {
			if d.DiscordId == f.Id {
				d.VdName = f.VdName

			}
		}
	}
	for _, fmv := range fmvl {
		for _, v := range vdl {
			if v.VelocidronName == fmv.VdName || fmv.RacerName == v.VelocidronName {
				fmv.QualifyingTime = v.QualifyingTime
				fmv.ModelName = v.ModelName
				break
			}
			var fmvNul FmvVoicePilot
			if fmv.QualifyingTime == fmvNul.QualifyingTime {
				fmv.QualifyingTime = "CHECK IN Please!"
			}

		}

	}
	return fmvl
}
