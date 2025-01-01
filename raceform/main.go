package main

import (
	"fmt"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/log"
	"golang.org/x/exp/maps"
)

func main() {
	var r = []string{"asiy", "MOTODRONEX", "eedok", "uGeLLin", "andyy", "MGescapades", "RoflCopter!", "AP3X"}
	var t1 = []string{"84.132", "84.135", "86.167", "84.132", "85.236", "88.968", "89.003", "92.542"}
	var t2 = []string{"83.142", "81.735", "85.437", "94.692", "89.111", "87.934", "85.303", "96.662"}

	var dirtyTime1 []string
	var dirtyTime2 []string
	for x := 0; x < len(r)-1; x++ {
		formRTimes := fmt.Sprintf("%s: %s", r[x], t1[x])
		dirtyTime1 = append(dirtyTime1, formRTimes)
	}
	for x := 0; x < len(r)-1; x++ {
		formRTimes := fmt.Sprintf("%s: %s", r[x], t2[x])
		dirtyTime2 = append(dirtyTime2, formRTimes)
	}

	heatDict := make(map[string][]string)
	heatDict["Gold S1R1"] = dirtyTime1
	heatDict["Gold S1R2"] = dirtyTime2
	log.SetReportTimestamp(false)

	var (
		group string
		racer []string
	)

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				//Options(huh.NewOptions("United States", "Canada", "Mexico")...).
				OptionsFunc(func() []huh.Option[string] {
					Keys := maps.Keys(heatDict)
					return huh.NewOptions(Keys...)
				}, &group /* only this function when `country` changes */).
				Value(&group).
				Title("Races").
				Height(5),

			huh.NewMultiSelect[string]().
				Value(&racer).
				Height(8).
				Title("Races").
				OptionsFunc(func() []huh.Option[string] {
					s := heatDict[group]
					time.Sleep(500 * time.Millisecond) // seems to be needed to allow optiont to load
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

	fmt.Printf("\n\nEntered: %s\n\n", group)
	for _, i := range racer {
		fmt.Println(i)
	}
	time.Sleep(10 * time.Second)
}
