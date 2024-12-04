package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
)

//var port int
//var AivisSpeechEndpoint string = "http://localhost:" + strconv.Itoa(port)

type Engine struct {
	Host string
	Port int
}

type Speaker struct {
	Name               string  `json:"name"`
	Uuid               string  `json:"speaker_uuid"`
	Styles             []Style `json:"styles"`
	Version            string  `json:"version"`
	Supported_features struct {
		Permited_synchesis_morphing string `json:"permitted_synthesis_morphing"`
	} `json:"supported_features"`
}

type Style struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

type AudioQuery struct {
	AccentPhrases []struct {
		Moras []struct {
			Text            string  `json:"text"`
			Consonant       string  `json:"consonant"`
			ConsonantLength float64 `json:"consonant_length"`
			Vowel           string  `json:"vowel"`
			VowelLength     float64 `json:"vowel_length"`
			Pitch           float64 `json:"pitch"`
		} `json:"moras"`
		Accent          int  `json:"accent"`
		PauseMora       any  `json:"pause_mora"`
		IsInterrogative bool `json:"is_interrogative"`
	} `json:"accent_phrases"`
	SpeedScale         float64 `json:"speedScale"`
	IntonationScale    float64 `json:"intonationScale"`
	TempoDynamicsScale float64 `json:"tempoDynamicsScale"`
	PitchScale         float64 `json:"pitchScale"`
	VolumeScale        float64 `json:"volumeScale"`
	PrePhonemeLength   float64 `json:"prePhonemeLength"`
	PostPhonemeLength  float64 `json:"postPhonemeLength"`
	PauseLength        any     `json:"pauseLength"`
	PauseLengthScale   float64 `json:"pauseLengthScale"`
	OutputSamplingRate int     `json:"outputSamplingRate"`
	OutputStereo       bool    `json:"outputStereo"`
	Kana               string  `json:"kana"`
}

func (e *Engine) getHost() string {
	//return net.JoinHostPort(e.Host, strconv.Itoa(e.Port))
	return e.Host + ":" + strconv.Itoa(e.Port)
}

func (e *Engine) getSpeakers() []Speaker {
	uri, _ := url.JoinPath(e.getHost(), "speakers")

	resp, err := http.Get(uri)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var speakers []Speaker
	if err := json.Unmarshal(body, &speakers); err != nil {
		panic(err)
	}
	return speakers
}

func (s *Style) getAudioQuery(host string, text string) AudioQuery {
	uri_param := url.Values{}
	uri_param.Set("text", text)
	uri_param.Set("speaker", strconv.Itoa(s.Id))

	uri, err := url.JoinPath(host, "audio_query")
	if err != nil {
		panic(err)
	}

	endpoint := uri + "?" + uri_param.Encode()

	resp, err := http.Post(endpoint, "application/json", nil)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var jsonData AudioQuery
	if err := json.Unmarshal(body, &jsonData); err != nil {
		fmt.Printf("error decoding sakura response: %v\n", err)
		if e, ok := err.(*json.SyntaxError); ok {
			fmt.Printf("syntax error at byte offset %d\n", e.Offset)
		}
		fmt.Printf("sakura response: %q\n", body)
		panic(err)
	}
	return jsonData
}

func (s *Style) getAudio(host string, query AudioQuery) {
	jsonQuery, _ := json.Marshal(query)
	uri_param := url.Values{}
	uri_param.Set("speaker", strconv.Itoa(s.Id))

	uri, _ := url.JoinPath(host, "synthesis")
	endpoint := uri + "?" + uri_param.Encode()

	resp, _ := http.Post(endpoint, "application/json", bytes.NewBuffer(jsonQuery))
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	st, format, err := wav.Decode(bytes.NewReader(body))
	if err != nil {
		panic(err)
	}
	defer st.Close()

	speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))
	done := make(chan bool)
	speaker.Play(beep.Seq(st, beep.Callback(func() {
		done <- true
	})))
	<-done
}

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
				e := Engine{Host: host, Port: port}
				for i, s := range e.getSpeakers() {
					color.Red(fmt.Sprintf("%d: %s\n", i, s.Name))
					for j, style := range s.Styles {
						color.Green(fmt.Sprintf("  %d: %d: %s\n", j, style.Id, style.Name))
					}
				}
				return nil
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
				e := Engine{Host: host, Port: port}
				speakers := e.getSpeakers()
				var style Style
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

				query := style.getAudioQuery(e.getHost(), c.Args().First())
				query.SpeedScale = c.Float64("speed")
				query.IntonationScale = c.Float64("intonation")
				query.TempoDynamicsScale = c.Float64("tempo")
				query.PitchScale = c.Float64("pitch")
				if c.Int("accent") != -1 {
					query.AccentPhrases[0].Accent = c.Int("accent")
				}
				style.getAudio(e.getHost(), query)

				return nil
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
