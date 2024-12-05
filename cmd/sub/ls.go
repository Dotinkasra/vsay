package sub

import (
	"fmt"
	"vsay/pkg/engine"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

func GetLsFlags() []cli.Flag {
	lsFlags := []cli.Flag{
		&cli.BoolFlag{
			Name:    "speaker",
			Aliases: []string{"s"},
			Usage:   "show speakers",
			Value:   true,
		},
		&cli.BoolFlag{
			Name:    "dict",
			Aliases: []string{"d"},
			Usage:   "show dictionaries",
			Value:   false,
		},
	}
	return lsFlags
}

func showSpeakers(e engine.Engine) {
	for i, s := range e.ShowSpeakers() {
		color.Red(fmt.Sprintf("%d: %s\n", i, s.Name))
		for j, style := range s.Styles {
			color.Green(fmt.Sprintf("\t%d: %d: %s\n", j, style.Id, style.Name))
		}
	}
}

func showDict(e engine.Engine) error {
	for id, d := range e.ShowUserDict() {
		color.Red(fmt.Sprintf("%s:\n", id))
		color.Green(fmt.Sprintf("\tID: %v\n\t単語: %v\n\t読み: %v\n\tアクセント: %v\n", d.ContextID, d.Surface, d.Yomi, d.AccentType))
	}
	return nil
}

func Ls(c *cli.Context, e engine.Engine) error {
	if c.Bool("speaker") {
		fmt.Printf("Speakers\n")
		showSpeakers(e)
		fmt.Println()
	}
	if c.Bool("dict") {
		fmt.Printf("User dictionary\n")
		showDict(e)
		fmt.Println()
	}
	return nil
}
