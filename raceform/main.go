package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
)

type racer string
type Racers []racer

var split1, split2, split3 bool
var names = Racers{"asiy", "MOTORDRONEX", "eedok", "uGellin", "MGescapades", "RoflCopter!", "AP3X"}

func main() {
	log.SetReportTimestamp(false)

	var (
		group string
		racer []string
	)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Options(huh.NewOptions("Gold", "Magenta", "Teal", "Orange", "Green")...).
				Value(&group).
				Title("Split 1 - Heat 1").
				Height(8),
			huh.NewMultiSelect[string]().
				Value(&racer).
				Height(8).
				TitleFunc(func() string {
					switch group {
					case "Gold":
						return "Gold Results"
					case "Magenta":
						return "Magenta Results"
					default:
						return "-"
					}
				}, &group).
				OptionsFunc(func() []huh.Option[string] {
					s := racers[group]
					// simulate API call
					time.Sleep(10 * time.Millisecond)
					return huh.NewOptions(s...)
				}, &group /* only this function when `group` changes */),
			huh.NewConfirm().
				Title("Submit Entries").
				Affirmative("Submit").
				Negative("Cancel").
				Validate(func(v bool) error {
					if !v {
						return fmt.Errorf("Welp, finish up then")
					}
					return nil
				}),
		),
	)

	err := form.Run()
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%s, %s\n", racer, group)
}

var racers = map[string][]string{
	"Gold": {
		"asiy: 80.639",
		"MOTODRONX: 94.757",
		"eedok: 82873",
		"UGeLLin: 83.046",
		"andyy: 87.688",
		"Mgescapades: 84894",
		"RoflCopter: 84.813",
		"AP3X: 87.058",
	},
	"Magenta": {
		"Not Sure: 85.257",
		"Barnyard: 87.378",
		"Mayan_Hawk: 110.308",
		".MrE.: 87.112",
		"XaeroFPV: 89.162",
		"Treeseeker: MIA",
	},
	"Teal":   {},
	"Orange": {},
	"Green":  {},
}
