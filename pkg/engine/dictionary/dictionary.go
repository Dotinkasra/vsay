package dictionary

import (
	"net/url"
	"strconv"
	"unsafe"
	"vsay/pkg/util"
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
	PROPER_NOUN WordType = iota
	COMMON_NOUN
	VERB
	ADJECTIVE
	SUFFIX
)

type DictRequest struct {
	Surface       string
	pronunciation string
	accent_type   int
	word_type     *WordType
	priority      *int
}

func (w WordType) String() string {
	return [...]string{"PROPER_NOUN", "COMMON_NOUN", "VERB", "ADJECTIVE", "SUFFIX"}[w]
}

func (d *DictRequest) RegisterUserDict(host string) (string, error) {
	uri, _ := url.JoinPath(host, "user_dict_word")
	url_param := url.Values{}
	url_param.Set("surface", d.Surface)
	url_param.Set("pronunciation", d.pronunciation)
	url_param.Set("accent_type", strconv.Itoa(d.accent_type))
	if d.word_type != nil {
		url_param.Set("word_type", d.word_type.String())
	}
	if d.priority != nil {
		url_param.Set("priority", strconv.Itoa(*d.priority))
	}
	endpoint := uri + "?" + url_param.Encode()
	resp, err := util.HttpPost(endpoint, nil)
	if err != nil {
		return "", err
	}
	return *(*string)(unsafe.Pointer(&resp)), nil

}
