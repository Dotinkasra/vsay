package engine

import (
	"encoding/json"
	"log"
	"net/url"
	"strconv"
	"vsay/pkg/engine/speaker"
	"vsay/pkg/util"
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
