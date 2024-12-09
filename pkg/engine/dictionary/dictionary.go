package dictionary

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"unsafe"
	"vsay/pkg/engine"
	"vsay/pkg/util"

	"github.com/fatih/color"
)

type Dictionary struct {
	Surface               string `json:"surface"`
	Priority              int    `json:"priority"`
	ContextID             int    `json:"context_id"`
	PartOfSpeech          string `json:"part_of_speech"`
	PartOfSpeechDetail1   string `json:"part_of_speech_detail_1"`
	PartOfSpeechDetail2   string `json:"part_of_speech_detail_2"`
	PartOfSpeechDetail3   string `json:"part_of_speech_detail_3"`
	InflectionalType      string `json:"inflectional_type"`
	InflectionalForm      string `json:"inflectional_form"`
	Stem                  string `json:"stem"`
	Yomi                  string `json:"yomi"`
	Pronunciation         string `json:"pronunciation"`
	AccentType            int    `json:"accent_type"`
	MoraCount             int    `json:"mora_count"`
	AccentAssociativeRule string `json:"accent_associative_rule"`
}

type WordType int

const (
	PROPERNOUN WordType = iota + 1
	COMMONNOUN
	VERB
	ADJECTIVE
	SUFFIX
)

type DictRequest struct {
	Surface       string
	Pronunciation string
	AccentType    int
	WordType      *WordType
	Priority      *int
}

func (w WordType) String() string {
	return [...]string{"PROPER_NOUN", "COMMON_NOUN", "VERB", "ADJECTIVE", "SUFFIX"}[w]
}

func (d *DictRequest) RegisterUserDict(e engine.Engine) (string, error) {
	uri, _ := url.JoinPath(e.MyHost(), "user_dict_word")
	urlParam := url.Values{}
	urlParam.Set("surface", d.Surface)
	urlParam.Set("pronunciation", d.Pronunciation)
	urlParam.Set("accent_type", strconv.Itoa(d.AccentType))
	if d.WordType != nil {
		urlParam.Set("word_type", d.WordType.String())
	}
	if d.Priority != nil {
		urlParam.Set("priority", strconv.Itoa(*d.Priority))
	}
	endpoint := uri + "?" + urlParam.Encode()
	resp, err := util.HTTPPost(endpoint, nil)
	if err != nil {
		return "", err
	}

	return *(*string)(unsafe.Pointer(&resp)), nil
}

func ShowUserDict(myHost string) map[string]Dictionary {
	uri, _ := url.JoinPath(myHost, "user_dict")
	body, err := util.HTTPGet(uri)
	if err != nil {
		log.Panic(err)
	}
	var userDict map[string]Dictionary
	if err = json.Unmarshal(body, &userDict); err != nil {
		log.Panic(err)
	}
	return userDict
}

func DeleteDict(myHost string, uuid string) error {
	uri, err := url.JoinPath(myHost, "user_dict_word", uuid)
	if err != nil {
		color.Red(fmt.Sprintln("Error: ホスト名かポートを間違えている可能性があります。"))
		log.Panic(err)
	}
	body, err := util.HTTPDelete(uri, nil)
	if err != nil {
		color.Red(fmt.Sprintln("Error: "))
		log.Panic(err)
	}
	fmt.Println(string(body))
	return err
}
