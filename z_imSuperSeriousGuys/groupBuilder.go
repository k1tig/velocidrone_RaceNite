package main

import (
	"os"

	"github.com/gocarina/gocsv"
)

func GetVdRacers(filename string) []Pilot {
	var pilot = []Pilot{}
	var filteredPilots = []Pilot{}

	raceFile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer raceFile.Close()
	if err := gocsv.UnmarshalFile(raceFile, &pilot); err != nil { // Load clients from file
		panic(err)
	}
	//use OkRaceClass as a landing spot to filter Velocidrone list by spec
	filteredPilots = append(filteredPilots, pilot...)

	if _, err := raceFile.Seek(0, 0); err != nil { // Go to the start of the file
		panic(err)
	}
	return filteredPilots
}
