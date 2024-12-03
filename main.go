package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/faiface/beep"
	"github.com/faiface/beep/speaker"
	"github.com/faiface/beep/wav"
)

var AivisSpeechEndpoint string = "http://localhost:10101"

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

func getSpeakers() []Speaker {
	uri := AivisSpeechEndpoint + "/speakers"
	resp, _ := http.Get(uri)
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

func (s *Style) getAudioQuery(text string) json.RawMessage {
	uri_param := url.Values{}
	uri_param.Set("text", text)
	uri_param.Set("speaker", strconv.Itoa(s.Id))

	uri := AivisSpeechEndpoint + "/audio_query?" + uri_param.Encode()

	resp, _ := http.Post(uri, "application/json", nil)
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

func getAudio(audio_query json.RawMessage) {
	uri_param := url.Values{}
	uri_param.Set("speaker", "888753765")
	uri := AivisSpeechEndpoint + "/synthesis?" + uri_param.Encode()
	resp, _ := http.Post(uri, "application/json", bytes.NewBuffer(audio_query))
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
	for i, s := range getSpeakers() {
		fmt.Printf("%d: %s\n", i, s.Name)
		for _, style := range s.Styles {
			fmt.Printf("  %d: %s\n", style.Id, style.Name)
			if style.Id == 888753765 {
				audio_query := style.getAudioQuery("こんにちは")
				getAudio(audio_query)
			}
		}
	}
}
