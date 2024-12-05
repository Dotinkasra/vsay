package sub

import (
	"fmt"
	"vsay/pkg/engine"

	"github.com/fatih/color"
)

func ShowDict(e engine.Engine) error {
	for id, d := range e.ShowUserDict() {
		color.Red(fmt.Sprintf("%s:\n", id))
		color.Green(fmt.Sprintf("\tID: %v\n\t単語: %v\n\t読み: %v\n\tアクセント: %v\n", d.ContextID, d.Surface, d.Yomi, d.AccentType))
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
