package main

import (
	"fmt"
	"time"

	rg "github.com/main/racegroup"
)

var groups = []string{"Yellow", "Magenta", "Cyan", "Gold", "Green"}

func main() {
	fmvcsv := rg.GetFMVvoice("checkin.csv")
	var fmvList []string
	for _, i := range fmvcsv {
		fmvList = append(fmvList, i.Racer)
	}

	fmvArray := rg.RaceArray(fmvList)
	fmt.Printf("\n\n\n FMV race list from CSV:\n\n")
	printGroup(fmvArray)

	// select entry from VD list to write Client.LapTime to fmvcsv.Qualifying time
	//fmt.Printf("Name: %s\nQualifying Time: %s\n\n", fmvcsv[1].Racer, fmvcsv[1].QualifyingTime)

	time.Sleep(30 * time.Second) // just to keep terminal open
}

func printGroup(group [][]string) {
	for x, i := range group {
		fmt.Printf("Group %s:", groups[x])
		for _, racer := range i {
			fmt.Printf(" %s,", racer)
		}
		fmt.Println()
	}
}
