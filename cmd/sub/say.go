package sub

import (
	"vsay/pkg/audio"
	"vsay/pkg/engine"
	"vsay/pkg/engine/speaker"

	"github.com/urfave/cli/v2"
)

func GetSayFlags() []cli.Flag {
	sayFlags := []cli.Flag{
		&cli.IntFlag{
			Name:     "id",
			Usage:    "Style `ID`. This takes priority over the speaker number option.",
			Required: false,
		},
		&cli.IntFlag{
			Name:     "number",
			Usage:    "The speaker number as displayed by the `ls` command.",
			Aliases:  []string{"n"},
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
	}
	return sayFlags
}

func Say(c *cli.Context, e engine.Engine) error {
	speakers := e.ShowSpeakers()
	var style speaker.Style
	if c.Int("id") == 0 {
		sp := speakers[c.Int("number")]
		st := sp.Styles[c.Int("style")]
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
}
