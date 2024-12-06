package sub

import (
	"fmt"
	"vsay/pkg/engine"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

type Ls struct {
	Cmd
}

func (scmd *Ls) GetFlags() []cli.Flag {
	lsFlags := []cli.Flag{
		&cli.BoolFlag{
			Name:    "speaker",
			Aliases: []string{"s"},
			Usage:   "show speakers",
			Value:   false,
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
			color.Green(fmt.Sprintf("\t%d: %d: %s\n", j, style.ID, style.Name))
		}
	}
}

func showDict(e engine.Engine) {
	for id, d := range e.ShowUserDict() {
		color.Red(fmt.Sprintf("%s:\n", id))
		color.Green(fmt.Sprintf("\tID: %v\n\t単語: %v\n\t読み: %v\n\tアクセント: %v\n", d.ContextID, d.Surface, d.Yomi, d.AccentType))
	}
}

func (scmd *Ls) Action(c *cli.Context) error {
	e := engine.Engine{Host: c.String("host"), Port: c.Int("port")}
	speakerFlag := c.Bool("speaker")
	dictFlag := c.Bool("dict")

	if !(speakerFlag || dictFlag) {
		speakerFlag = !speakerFlag
		dictFlag = !dictFlag
	}

	if speakerFlag {
		fmt.Printf("Speakers\n")
		showSpeakers(e)
		fmt.Println()
	}
	if dictFlag {
		fmt.Printf("User dictionary\n")
		showDict(e)
		fmt.Println()
	}
	return nil
}
