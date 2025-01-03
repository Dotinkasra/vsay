package sub

import (
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"unsafe"
	"vsay/pkg/audio"
	"vsay/pkg/engine"
	"vsay/pkg/engine/speaker"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

type Say struct {
	Cmd
	ShowSpeaker
}

type ShowSpeaker struct {
	Cmd
}

func (scmd *Say) GetFlags() []cli.Flag {
	flags := []cli.Flag{
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
		&cli.BoolFlag{
			Name:     "b64",
			Aliases:  []string{"b"},
			Usage:    "Outputs audio as base64 encoding to Stdout.",
			Required: false,
		},
	}
	return flags
}

func (scmd *Say) Action(c *cli.Context) error {
	e := engine.Engine{Host: c.String("host"), Port: c.Int("port")}

	speakers := speaker.ShowSpeakers(e.MyHost())
	var style speaker.Style
	if c.Int("id") == 0 {
		sp := speakers[c.Int("number")]
		st := sp.Styles[c.Int("style")]
		style = st
	} else {
		for _, sp := range speakers {
			for _, st := range sp.Styles {
				if st.ID == c.Int("id") {
					style = st
				}
			}
		}
	}
	text := c.Args().First()
	if text == "" {
		stdin := os.Stdin
		t, err := io.ReadAll(stdin)
		if err != nil {
			log.Panic(err)
		}
		text = *(*string)(unsafe.Pointer(&t))
	}
	text = strings.TrimSpace(text)
	text = strings.TrimSuffix(text, "\n")

	query := style.CreateAudioQuery(e.MyHost(), text)
	query.SpeedScale = c.Float64("speed")
	query.IntonationScale = c.Float64("intonation")
	query.TempoDynamicsScale = c.Float64("tempo")
	query.PitchScale = c.Float64("pitch")
	if c.Int("accent") != -1 {
		query.AccentPhrases[0].Accent = c.Int("accent")
	}

	rawAudio := style.GetAudio(e.MyHost(), query)
	if !c.Bool("quiet") {
		audio.PlayAudio(rawAudio)
	}
	if c.String("save") != "" {
		audio.SaveAudio(rawAudio, c.String("save"))
	}
	if c.Bool("b64") {
		b64Audio := base64.StdEncoding.EncodeToString(rawAudio)
		fmt.Print(b64Audio)
	}
	return nil
}

func (scmd *ShowSpeaker) GetFlags() []cli.Flag {
	return []cli.Flag{}
}
func (scmd *ShowSpeaker) Action(c *cli.Context) error {
	e := engine.Engine{Host: c.String("host"), Port: c.Int("port")}
	for i, s := range speaker.ShowSpeakers(e.MyHost()) {
		color.Red(fmt.Sprintf("%d: %s\n", i, s.Name))
		for j, style := range s.Styles {
			color.Green(fmt.Sprintf("\t%d: %d: %s\n", j, style.ID, style.Name))
		}
	}
	return nil
}
