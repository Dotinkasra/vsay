package main

import (
	"log"
	"os"
	"slices"
	"vsay/cmd/sub"
	"vsay/pkg/engine"

	"github.com/urfave/cli/v2"
)

//var port int
//var AivisSpeechEndpoint string = "http://localhost:" + strconv.Itoa(port)

func main() {
	var host string
	var port int

	app := cli.NewApp()
	app.Name = "vsay"
	app.Usage = "Synthesized voice is played from the terminal."
	app.UseShortOptionHandling = true

	baseFlags := []cli.Flag{
		&cli.StringFlag{
			Name:        "host",
			Usage:       "Host address",
			Value:       "http://localhost",
			Destination: &host,
		},
		&cli.IntFlag{
			Name:        "port",
			Usage:       "Port number",
			Aliases:     []string{"p"},
			Value:       10101,
			Destination: &port,
		},
	}

	app.Commands = []*cli.Command{
		{
			Name:    "ls",
			Aliases: []string{"l"},
			Usage:   "Show speakers",
			Flags:   slices.Concat(baseFlags, sub.GetLsFlags()),
			Action: func(c *cli.Context) error {
				e := engine.Engine{Host: host, Port: port}
				return sub.Ls(c, e)
			},
		},
		{
			Name:    "say",
			Aliases: []string{"s"},
			Usage:   "Say something",
			Flags:   slices.Concat(baseFlags, sub.GetSayFlags()),
			Action: func(c *cli.Context) error {
				e := engine.Engine{Host: host, Port: port}
				return sub.Say(c, e)
			},
		},
		{
			Name:    "dict",
			Aliases: []string{"d"},
			Usage:   "Show dictionary",
			Flags:   baseFlags,
			Action: func(c *cli.Context) error {
				//e := engine.Engine{Host: host, Port: port}
				return nil
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
