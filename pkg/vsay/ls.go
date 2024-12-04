package vsay

import (
	"fmt"
	"vsay/pkg/engine"

	"github.com/fatih/color"
)

func Ls(e engine.Engine) error {
	for i, s := range e.GetSpeakers() {
		color.Red(fmt.Sprintf("%d: %s\n", i, s.Name))
		for j, style := range s.Styles {
			color.Green(fmt.Sprintf("  %d: %d: %s\n", j, style.Id, style.Name))
		}
	}
	return nil
}
