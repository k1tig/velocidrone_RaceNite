package main

import (
	"os"

	"github.com/gocarina/gocsv"
)

type Client struct { // Our example struct, you can use "-" to ignore a field
	PlayerName string `csv:"Player Name"`
	LapTime    string `csv:"Lap Time"`
	X_Pos      string `csv:"-"`
	ModelName  string `csv:"Model Name"`
	X_Country  string `csv:"-"`
}

var clients = []*Client{}
var checkedIn = []*Client{}

func getRacers() {
	raceFile, err := os.OpenFile("race.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}

	defer raceFile.Close()

	if err := gocsv.UnmarshalFile(raceFile, &clients); err != nil { // Load clients from file
		panic(err)
	}

	for _, client := range clients {
		if client.ModelName == "TBS Spec" || client.ModelName == "Twig XL 3" {
			//fmt.Printf("Racer: %s\nTime: %v\nModel: %s\n\n", racer.DisplayName, client.LapTime, client.ModelName)
			checkedIn = append(checkedIn, client)
			break //to put the time with the faster quad
		}
	}

	if _, err := raceFile.Seek(0, 0); err != nil { // Go to the start of the file
		panic(err)
	}
}

/* err = gocsv.MarshalFile(&clients, clientsFile) // Use this to save the CSV back to the file
if err != nil {
	panic(err)
} */

func main() {
	getRacers()

}
