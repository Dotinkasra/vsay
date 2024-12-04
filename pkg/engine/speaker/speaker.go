package speaker

import (
	"bytes"
	"encoding/json"
	"net/url"
	"strconv"
	"vsay/pkg/util"
)

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

func (s *Style) CreateAudioQuery(host string, text string) AudioQuery {
	uri_param := url.Values{}
	uri_param.Set("text", text)
	uri_param.Set("speaker", strconv.Itoa(s.Id))

	uri, _ := url.JoinPath(host, "audio_query")
	endpoint := uri + "?" + uri_param.Encode()

	body, err := util.HttpPost(endpoint, nil)
	if err != nil {
		panic(err)
	}

	var query AudioQuery
	if err := json.Unmarshal(body, &query); err != nil {
		panic(err)
	}
	return query
}

func (s *Style) GetAudio(host string, query AudioQuery) []byte {
	jsonQuery, _ := json.Marshal(query)
	uri_param := url.Values{}
	uri_param.Set("speaker", strconv.Itoa(s.Id))

	uri, _ := url.JoinPath(host, "synthesis")
	endpoint := uri + "?" + uri_param.Encode()

	body, _ := util.HttpPost(endpoint, bytes.NewBuffer(jsonQuery))
	return body
}
