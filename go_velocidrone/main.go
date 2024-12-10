package main

import (
	"fmt"
	"os"
	"time"

	"github.com/gocarina/gocsv"
)

var running int

type Client struct { // Our example struct, you can use "-" to ignore a field
	PlayerName string `csv:"Player Name"`
	LapTime    string `csv:"Lap Time"`
	X_Pos      string `csv:"-"`
	ModelName  string `csv:"Model Name"`
	X_Country  string `csv:"-"`
}

type Racer struct {
	DisplayName string `csv:"Display Name"`
	VdName      string `csv:"Display Name"` // placeholder to update name linked to Velocidrone CSV PlayerName
}

func main() {
	running = 1
	for running == 1 {
		raceFile, err := os.OpenFile("race.csv", os.O_RDWR|os.O_CREATE, os.ModePerm)
		if err != nil {
			panic(err)
		}

		checkinFile, err := os.OpenFile("checkin.csv", os.O_RDONLY|os.O_CREATE, os.ModePerm)
		if err != nil {
			panic(err)
		}

		defer raceFile.Close()
		clients := []*Client{}
		racers := []*Racer{}

		checkedIn := []*Client{}

		if err := gocsv.UnmarshalFile(raceFile, &clients); err != nil { // Load clients from file
			panic(err)
		}

		if err := gocsv.UnmarshalFile(checkinFile, &racers); err != nil {
			panic(err)
		}

		for _, racer := range racers {
			for _, client := range clients {
				if racer.DisplayName == client.PlayerName {
					if client.ModelName == "TBS Spec" || client.ModelName == "Twig XL 3" {
						//fmt.Printf("Racer: %s\nTime: %v\nModel: %s\n\n", racer.DisplayName, client.LapTime, client.ModelName)
						checkedIn = append(checkedIn, client)
						break //to put the time with the faster quad
					}
				}
			}
		}

		if _, err := raceFile.Seek(0, 0); err != nil { // Go to the start of the file
			panic(err)
		}
		if _, err := checkinFile.Seek(0, 0); err != nil { // Go to the start of the file
			panic(err)
		}

		missingRacers := []*Racer{}
		for _, racer := range racers {
			racerMissing := true
			for _, notCheckedIn := range checkedIn {
				if notCheckedIn.PlayerName == racer.DisplayName {
					racerMissing = false
					break
				}
			}
			// create func to check list of VD entries
			if racerMissing {
				fmt.Printf("Pair Racer '%s': ", racer.DisplayName)
				var input string
				//input to pair racer name with VD name, needs cleaning
				_, err := fmt.Scan(&input)
				if err != nil {
					fmt.Println("error:", err)
					return
				}
				for _, client := range clients {
					if input == client.PlayerName {
						racer.VdName = input
						racerMissing = false
						break
					}
				}
				if racerMissing {
					fmt.Printf("\n!Could not pair input of '%s' with racer : %s!\n\n", input, racer.DisplayName)
					missingRacers = append(missingRacers, racer)
				}
				//missingRacers = append(missingRacers, racer)
			}
		}
		fmt.Printf("\n\n\n")
		fmt.Printf("Total Racers: %d\n", len(racers))
		fmt.Printf("Total Racers Checked in: %d\n", len(checkedIn))
		fmt.Printf("Racers to be paired: %d\n\n", len(missingRacers))

		for _, missing := range missingRacers {
			fmt.Println("Missing Racer: ", missing.DisplayName)
		}

		time.Sleep(100 * time.Second)
		running = 0
	}
}

/* err = gocsv.MarshalFile(&clients, clientsFile) // Use this to save the CSV back to the file
if err != nil {
	panic(err)
} */
