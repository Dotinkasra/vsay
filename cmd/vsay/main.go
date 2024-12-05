package main

import (
	"log"
	"os"
	"slices"
	"vsay/cmd/sub"

	"github.com/urfave/cli/v2"
)

func MakeFlags(scmd sub.SubCommand) []cli.Flag {
	baseFlags := []cli.Flag{
		&cli.StringFlag{
			Name:  "host",
			Usage: "Host address",
			Value: "http://localhost",
		},
		&cli.IntFlag{
			Name:    "port",
			Usage:   "Port number",
			Aliases: []string{"p"},
			Value:   10101,
		},
	}
	return slices.Concat(baseFlags, scmd.GetFlags())
}

func main() {
	ls := sub.Ls{}
	say := sub.Say{}
	dict := sub.Dict{}

	app := cli.NewApp()
	app.Name = "vsay"
	app.Usage = "Synthesized voice is played from the terminal."
	app.UseShortOptionHandling = true

	app.Commands = []*cli.Command{
		{
			Name:    "ls",
			Aliases: []string{"l"},
			Usage:   "Show speakers",
			Flags:   MakeFlags(&ls),
			Action: func(c *cli.Context) error {
				return ls.Action(c)
			},
		},
		{
			Name:    "say",
			Aliases: []string{"s"},
			Usage:   "Say something",
			Flags:   MakeFlags(&say),
			Action: func(c *cli.Context) error {
				return say.Action(c)
			},
		},
		{
			Name:    "dict",
			Aliases: []string{"d"},
			Usage:   "Show dictionary",
			Subcommands: []*cli.Command{
				{
					Name:    "add",
					Aliases: []string{"a"},
					Usage:   "Add word",
					Flags:   MakeFlags(&dict.DictAdd),
					Action: func(c *cli.Context) error {
						return dict.DictAdd.Action(c)
					},
				},
				{
					Name:    "delete",
					Aliases: []string{"r"},
					Usage:   "Remove word",
					Flags:   MakeFlags(&dict.DictDelete),
					Action: func(c *cli.Context) error {
						return dict.DictDelete.Action(c)
					},
				},
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
