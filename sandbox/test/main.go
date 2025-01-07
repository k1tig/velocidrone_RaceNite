package main

import (
	"encoding/csv"
	"os"
)

type Pilot struct {
	pilotName, qTime, totalScore string
	h1, h2, h3, h4, h5           string
	h6, h7, h8, h9, h10          string
	p1, p2, p3, p4, p5           string
	p6, p7, p8, p9, p10          string
	s2, s3                       string
}

var fmvVoiceChat = []string{"Knee", "IQ0", "asiy", "eedok", "kalli", "dapaca", "uGeLLin", "AP3X", "SITHironoid"}
var qTimes = []string{"52.33", "54.239", "59.551", "58.475", "62.731", "56.913", "58.518", "69.356", "60.754"}
var pilots []Pilot

func makePilots() {
	for i := 0; i < len(fmvVoiceChat)-1; i++ {
		var x Pilot
		x.pilotName = fmvVoiceChat[i]
		x.qTime = qTimes[i]
		pilots = append(pilots, x)
	}
}

func main() {
	// Create a new CSV file
	file, err := os.Create("../try.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Create a CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header row
	headers := []string{"Pilot Name", "Qualifying Time", "Total Score",
		"H1 Time", "H1 Pts", "H2 Time", "H2 Pts", "H3 Time", "H3 Pts", "H4 Time", "H4 Pts",
		"H5 Time", "H5 Pts", "H6 Time", "H6 Pts", "H7 Time", "H7 Pts", "H8 Time", "H8 Pts",
		"H9 Time", "H9 Pts", "H10 Time", "H10 Pts"}

	writer.Write(headers)

	// Write data rows
	data := [][]string{}
	makePilots()
	for _, i := range pilots {
		racer := []string{i.pilotName, i.qTime, i.totalScore,
			i.h1, i.p1, i.h2, i.p2, i.h3, i.p3, i.h4, i.p4, i.h5, i.p5,
			i.h6, i.p6, i.h7, i.p7, i.h8, i.p8, i.h9, i.p9, i.h10, i.p10}
		data = append(data, racer)
	}

	writer.WriteAll(data)
}
