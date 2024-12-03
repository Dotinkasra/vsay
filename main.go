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

func (s *Style) getAudioQuery(host string, text string) json.RawMessage {
	uri_param := url.Values{}
	uri_param.Set("text", text)
	uri_param.Set("speaker", strconv.Itoa(s.Id))

	uri, err := url.JoinPath(host, "audio_query")
	if err != nil {
		panic(err)
	}

	endpoint := uri + "?" + uri_param.Encode()
	fmt.Println(uri)

	resp, err := http.Post(endpoint, "application/json", nil)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var jsonData json.RawMessage
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

func (e *Engine) getAudio(audio_query json.RawMessage) {
	uri_param := url.Values{}
	uri_param.Set("speaker", "888753765")

	uri, _ := url.JoinPath(e.getHost(), "synthesis")
	endpoint := uri + "?" + uri_param.Encode()

	resp, _ := http.Post(endpoint, "application/json", bytes.NewBuffer(audio_query))
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

func showSpeakers(e *Engine) {
	for i, s := range e.getSpeakers() {
		color.Red(fmt.Sprintf("%d: %s\n", i, s.Name))
		for j, style := range s.Styles {
			color.Green(fmt.Sprintf("  %d: %d: %s\n", j, style.Id, style.Name))
			if style.Id == 888753765 {
				audio_query := style.getAudioQuery(e.getHost(), "こんにちは")
				e.getAudio(audio_query)
			}
		}
	}
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
				engine := Engine{Host: host, Port: port}
				showSpeakers(&engine)
				return nil
			},
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
