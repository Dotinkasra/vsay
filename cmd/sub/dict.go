package sub

import (
	"fmt"
	"vsay/pkg/engine"

	"github.com/fatih/color"
)

func showDict(e engine.Engine) error {
	for i, s := range e.ShowSpeakers() {
		fmt.Printf("%d: %s\n", i, s.Name)
		for j, style := range s.Styles {
			fmt.Printf("  %d: %d: %s\n", j, style.Id, style.Name)
		}
	}
	return nil
}

func Dict(e engine.Engine) error {
	for i, s := range e.ShowSpeakers() {
		color.Red(fmt.Sprintf("%d: %s\n", i, s.Name))
		for j, style := range s.Styles {
			color.Green(fmt.Sprintf("  %d: %d: %s\n", j, style.Id, style.Name))
		}
	}
	return nil
}
