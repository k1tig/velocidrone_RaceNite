package main

import (
	"fmt"
	"time"

	rg "github.com/main/racegroup"
)

var groups = []string{"Yellow", "Magenta", "Cyan", "Gold", "Green"}
var vdList = []string{"jonE5", "Sam", "BagsFPV", "Treeseeker", "Blasta", "JimmyFPV", "MoarSparkles", "Krhom", "K1tig", "Dave",
	"Mr E", "AlsoMrE", "DeMic", "SillyYogurt", "HooHoo", "HooHa", "Max", "Timmy", "Tommy", "MeMaw",
	"jonE5", "Sam", "BagsFPV", "Treeseeker", "Blasta", "JimmyFPV", "MoarSparkles", "Krhom", "K1tig", "Dave",
	"Mr E", "AlsoMrE", "DeMic", "SillyYogurt", "HooHoo", "HooHa", "Max", "Timmy", "Tommy", "MeMaw",
	"Mr E", "AlsoMrE", "DeMic", "SillyYogurt", "HooHoo", "HooHa", "Max", "Timmy", "Tommy", "MeMaw"}

func main() {
	finalGroup := rg.RaceArray(vdList)
	fmt.Printf("\n\nFMV Bracket List:\n\n")
	printGroup(finalGroup)
	/*for x, i := range finalGroup {
		fmt.Printf("Group %s: ", groups[x])
		for _, racer := range i {
			fmt.Printf(" %s,", racer)
		}
		fmt.Println()
	}*/

	fmvcsv := rg.GetFMVvoice("checkin.csv")
	var fmvList []string
	for _, i := range fmvcsv {
		fmvList = append(fmvList, i.Racer)
	}

	fmvArray := rg.RaceArray(fmvList)
	fmt.Printf("\n\n\n FMV race list from CSV:\n\n")
	printGroup(fmvArray)

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
