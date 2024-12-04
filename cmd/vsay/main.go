package main

import (
	"log"
	"os"
	"vsay/pkg/engine"
	"vsay/pkg/engine/speaker"
	"vsay/pkg/vsay"

	"github.com/urfave/cli/v2"
)

//var port int
//var AivisSpeechEndpoint string = "http://localhost:" + strconv.Itoa(port)

func main() {
	var host string
	var port int

	app := cli.NewApp()
	app.Name = "vsay"
	app.Usage = "Test app for speakers"
	app.Flags = []cli.Flag{
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
			Action: func(c *cli.Context) error {
				e := engine.Engine{Host: host, Port: port}
				return vsay.Ls(e)
			},
		},
		{
			Name:    "say",
			Aliases: []string{"s"},
			Usage:   "Say something",
			Flags: []cli.Flag{
				&cli.IntFlag{
					Name:     "id",
					Usage:    "Speaker ID",
					Required: false,
				},
				&cli.IntFlag{
					Name:     "speaker-number",
					Aliases:  []string{"sn"},
					Usage:    "Speaker Number",
					Required: false,
				},
				&cli.IntFlag{
					Name:     "style-number",
					Aliases:  []string{"style"},
					Usage:    "Text to say",
					Required: false,
				},
				&cli.IntFlag{
					Name:     "accent",
					Usage:    "Text to say",
					Value:    -1,
					Required: false,
				},
				&cli.Float64Flag{
					Name:     "speed",
					Usage:    "Text to say",
					Value:    1.0,
					Required: false,
				},
				&cli.Float64Flag{
					Name:     "intonation",
					Usage:    "Text to say",
					Value:    1.0,
					Required: false,
				},
				&cli.Float64Flag{
					Name:     "tempo",
					Usage:    "Text to say",
					Value:    1.0,
					Required: false,
				},
				&cli.Float64Flag{
					Name:     "pitch",
					Usage:    "Text to say",
					Value:    0.0,
					Required: false,
				},
			},
			Action: func(c *cli.Context) error {
				e := engine.Engine{Host: host, Port: port}
				speakers := e.GetSpeakers()
				var style speaker.Style
				if c.Int("id") == 0 {
					sp := speakers[c.Int("speaker-number")]
					st := sp.Styles[c.Int("style-number")]
					style = st

				} else {
					for _, sp := range speakers {
						for _, st := range sp.Styles {
							if st.Id == c.Int("id") {
								style = st
							}
						}
					}
				}

				query := style.GetAudioQuery(e.GetHost(), c.Args().First())
				query.SpeedScale = c.Float64("speed")
				query.IntonationScale = c.Float64("intonation")
				query.TempoDynamicsScale = c.Float64("tempo")
				query.PitchScale = c.Float64("pitch")
				if c.Int("accent") != -1 {
					query.AccentPhrases[0].Accent = c.Int("accent")
				}
				style.GetAudio(e.GetHost(), query)

				return nil
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
