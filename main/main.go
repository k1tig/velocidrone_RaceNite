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
	fmt.Printf("\n\nFMV Bracket List:\n")
	for x, i := range finalGroup {
		fmt.Printf("Group %s: ", groups[x])
		for _, racer := range i {
			fmt.Println(racer)
			fmt.Println("")
		}
	}
	fmvlist := rg.GetFMVvoice()

	for _, i := range fmvlist {
		fmt.Println(i.Racer)
	}

	time.Sleep(30 * time.Second) // just to keep terminal open
}
