package main

import (
	"log"
	"os"
	"vsay/pkg/audio"
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
	app.Usage = "Synthesized voice is played from the terminal."
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
					Usage:    "Style `ID`. This takes priority over the speaker number option.",
					Required: false,
				},
				&cli.IntFlag{
					Name:     "number",
					Usage:    "The speaker number as displayed by the `ls` command.",
					Required: false,
				},
				&cli.IntFlag{
					Name:     "style",
					Usage:    "The style number as displayed by the `ls` command.",
					Required: false,
				},
				&cli.IntFlag{
					Name:     "accent",
					Usage:    "Specify the accent by its `index` in the string.",
					Value:    -1,
					Required: false,
				},
				&cli.Float64Flag{
					Name:     "speed",
					Usage:    "Set the speaking speed. Valid range: `0.5 to 2.0`.",
					Value:    1.0,
					Required: false,
				},
				&cli.Float64Flag{
					Name:     "intonation",
					Usage:    "Set the intonation, affecting the style's strength. Valid range: `0.0 to 2.0`.",
					Value:    1.0,
					Required: false,
				},
				&cli.Float64Flag{
					Name:     "tempo",
					Usage:    "Set the tempo. Valid range: `0.0 to 2.0`.",
					Value:    1.0,
					Required: false,
				},
				&cli.Float64Flag{
					Name:     "pitch",
					Usage:    "Set the pitch. Valid range: `-0.15 to 0.15`.",
					Value:    0.0,
					Required: false,
				},
				&cli.StringFlag{
					Name:     "save",
					Aliases:  []string{"s"},
					Usage:    "Specify the `PATH` to save the audio file.",
					Value:    "",
					Required: false,
				},
				&cli.BoolFlag{
					Name:     "quiet",
					Aliases:  []string{"q"},
					Usage:    "Don't play audio.",
					Required: false,
				},
			},
			Action: func(c *cli.Context) error {
				e := engine.Engine{Host: host, Port: port}
				speakers := e.ShowSpeakers()
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

				query := style.CreateAudioQuery(e.MyHost(), c.Args().First())
				query.SpeedScale = c.Float64("speed")
				query.IntonationScale = c.Float64("intonation")
				query.TempoDynamicsScale = c.Float64("tempo")
				query.PitchScale = c.Float64("pitch")
				if c.Int("accent") != -1 {
					query.AccentPhrases[0].Accent = c.Int("accent")
				}
				raw_audio := style.GetAudio(e.MyHost(), query)
				if !c.Bool("quiet") {
					audio.PlayAudio(raw_audio)
				}
				if c.String("save") != "" {
					audio.SaveAudio(raw_audio, c.String("save"))
				}
				return nil
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
