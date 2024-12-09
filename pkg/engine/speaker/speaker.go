package speaker

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
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
	Name              string  `json:"name"`
	UUID              string  `json:"speaker_uuid"`
	Styles            []Style `json:"styles"`
	Version           string  `json:"version"`
	SupportedFeatures struct {
		PermitedSynchesisMorphing string `json:"permitted_synthesis_morphing"`
	} `json:"supported_features"`
}

type Style struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

func (s *Style) CreateAudioQuery(host string, text string) AudioQuery {
	urlParam := url.Values{}
	urlParam.Set("text", text)
	urlParam.Set("speaker", strconv.Itoa(s.ID))

	uri, _ := url.JoinPath(host, "audio_query")
	endpoint := uri + "?" + urlParam.Encode()

	body, err := util.HTTPPost(endpoint, nil)
	if err != nil {
		log.Panic(err)
	}

	var query AudioQuery
	if err = json.Unmarshal(body, &query); err != nil {
		log.Panic(err)
	}
	return query
}

func (s *Style) GetAudio(host string, query AudioQuery) []byte {
	jsonQuery, _ := json.Marshal(query)
	urlParam := url.Values{}
	urlParam.Set("speaker", strconv.Itoa(s.ID))

	uri, _ := url.JoinPath(host, "synthesis")
	endpoint := uri + "?" + urlParam.Encode()

	body, err := util.HTTPPost(endpoint, bytes.NewBuffer(jsonQuery))
	if err != nil {
		log.Panic(err)
	}
	return body
}

func InitializeSpeaker(host string, sytleID int, skipReinit bool) error {
	urlParam := url.Values{}
	urlParam.Add("speaker", strconv.Itoa(sytleID))
	urlParam.Add("skip_reinit", strconv.FormatBool(skipReinit))

	uri, _ := url.JoinPath(host, "initialize_speaker")
	endpoint := uri + "?" + urlParam.Encode()

	_, err := util.HTTPGet(endpoint)
	if err != nil {
		return err
	}
	return nil
}

func InstallAivmModels(host string, filePath string, aivmURL string) ([]byte, error) {
	uri, _ := url.JoinPath(host, "aivm_models", "install")
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if filePath != "" {
		fileWriter, err := writer.CreateFormFile("file", filepath.Base(filePath))
		if err != nil {
			return nil, err
		}
		file, err := os.Open(filePath)
		if err != nil {
			return nil, err
		}
		defer file.Close()
		_, err = io.Copy(fileWriter, file)
		if err != nil {
			return nil, err
		}
	}

	if filePath != "" {
		// URLフィールドの設定
		encodedURL, err := url.Parse(aivmURL)
		if err != nil {
			return nil, err
		}
		err = writer.WriteField("url", encodedURL.String())
		if err != nil {
			return nil, err
		}
	}
	writer.Close()

	req, err := http.NewRequest(http.MethodPost, uri, body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Accept", "*/*")
	req.ContentLength = int64(body.Len())

	resp, err := util.HTTPPostWithMultipart(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func ShowSpeakers(myHost string) []Speaker {
	uri, _ := url.JoinPath(myHost, "speakers")
	body, err := util.HTTPGet(uri)
	if err != nil {
		log.Panic(err)
	}

	var speakers []Speaker
	if err = json.Unmarshal(body, &speakers); err != nil {
		log.Panic(err)
	}
	return speakers
}
