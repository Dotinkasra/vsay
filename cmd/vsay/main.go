package main

import (
	"log"
	"os"
	"slices"
	"vsay/cmd/sub"

	"github.com/urfave/cli/v2"
)

const defaultPort = 10101

func MakeFlags(scmd sub.Cmd) []cli.Flag {
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
			Value:   defaultPort,
		},
	}
	return slices.Concat(baseFlags, scmd.GetFlags())
}

func main() {
	say := sub.Say{}
	dict := sub.Dict{}
	install := sub.Install{}

	app := cli.NewApp()
	app.Name = "vsay"
	app.Usage = "Synthesized voice is played from the terminal."
	app.UseShortOptionHandling = true

	app.Commands = []*cli.Command{
		{
			Name:  "say",
			Usage: "Say something",
			Flags: MakeFlags(&say),
			Action: func(c *cli.Context) error {
				return say.Action(c)
			},
			Subcommands: []*cli.Command{
				{
					Name:    "ls",
					Aliases: []string{"l"},
					Usage:   "Show speakers",
					Flags:   MakeFlags(&say.ShowSpeaker),
					Action: func(c *cli.Context) error {
						return say.ShowSpeaker.Action(c)
					},
				},
			},
		},
		{
			Name:  "dict",
			Usage: "Show dictionary",
			Subcommands: []*cli.Command{
				{
					Name:    "add",
					Aliases: []string{"a"},
					Usage:   "Add word",
					Flags:   MakeFlags(&dict.AddDict),
					Action: func(c *cli.Context) error {
						return dict.AddDict.Action(c)
					},
				},
				{
					Name:    "delete",
					Aliases: []string{"r"},
					Usage:   "Remove word",
					Flags:   MakeFlags(&dict.DeleteDict),
					Action: func(c *cli.Context) error {
						return dict.DeleteDict.Action(c)
					},
				},
				{
					Name:    "ls",
					Aliases: []string{"l"},
					Usage:   "Show dictionary",
					Flags:   MakeFlags(&dict.ShowDict),
					Action: func(c *cli.Context) error {
						return dict.ShowDict.Action(c)
					},
				},
			},
		},
		{
			Name:  "install",
			Usage: "Show version",
			Flags: MakeFlags(&install),
			Action: func(c *cli.Context) error {
				return install.Action(c)
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
