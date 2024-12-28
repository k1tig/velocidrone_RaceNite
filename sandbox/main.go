package main

import "fmt"

type racer string
type colors []string

type races struct {
	heat  int
	racer string
	time  float32
}

// allows for logic such as; if len(bracketGroup.races) {submitSplit = flase}
type bracketGroup struct {
	split  int      //
	name   []colors // name of bracket group
	races  []races
	racers []racer // master list of racers to check before tallying points before split.
}

func main() {
	fmt.Println("Yes")
}
