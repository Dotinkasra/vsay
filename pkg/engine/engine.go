package engine

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"strconv"
	"vsay/pkg/engine/dictionary"
	"vsay/pkg/engine/speaker"
	"vsay/pkg/util"

	"github.com/fatih/color"
)

type Engine struct {
	Host string
	Port int
}

func (e *Engine) MyHost() string {
	return e.Host + ":" + strconv.Itoa(e.Port)
}

func (e *Engine) ShowSpeakers() []speaker.Speaker {
	uri, _ := url.JoinPath(e.MyHost(), "speakers")
	body, err := util.HttpGet(uri)
	if err != nil {
		log.Panic(err)
	}

	var speakers []speaker.Speaker
	if err := json.Unmarshal(body, &speakers); err != nil {
		log.Panic(err)
	}
	return speakers
}

func (e *Engine) ShowUserDict() map[string]dictionary.Dictionary {
	uri, _ := url.JoinPath(e.MyHost(), "user_dict")
	body, err := util.HttpGet(uri)
	if err != nil {
		log.Panic(err)
	}
	var userDict map[string]dictionary.Dictionary
	if err := json.Unmarshal(body, &userDict); err != nil {
		log.Panic(err)
	}
	return userDict
}

func (e *Engine) DeleteDict(uuid string) error {
	uri, err := url.JoinPath(e.MyHost(), "user_dict_word", uuid)
	if err != nil {
		color.Red(fmt.Sprintln("Error: ホスト名かポートを間違えている可能性があります。"))
		log.Panic(err)
	}
	body, err := util.HttpDelete(uri, nil)
	if err != nil {
		color.Red(fmt.Sprintln("Error: "))
		log.Panic(err)
	}
	fmt.Println(string(body))
	return err

}
