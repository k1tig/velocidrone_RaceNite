package main

import "github.com/charmbracelet/huh"

var races [][]string

type racerList struct {
	racers []struct {
		name string
		id   string
	}
}

type Model struct {
	form *huh.Form
}

func newModel() Model {
	var racedata = [][]string{{"EeDocking", "69.420"}, {"MrE", "69.69"}, {"Demic", "180.00"}, {"PrettySure", "65.10"}}
	races = racedata
	m := Model{}

	m.form = huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Key("yes").
				Title("X").
				Description("y"),
		),
	).
		WithWidth(45).
		WithShowHelp(false).
		WithShowErrors(false)

	return m
}

func main() {

}

/*


information to be sent:

raceAPI send all the racers and their times
select times to be sent

*/
