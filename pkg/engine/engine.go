package engine

import (
	"encoding/json"
	"net/url"
	"strconv"
	"vsay/pkg/engine/speaker"
	"vsay/pkg/util"
)

type Engine struct {
	Host string
	Port int
}

func (e *Engine) GetHost() string {
	return e.Host + ":" + strconv.Itoa(e.Port)
}

func (e *Engine) GetSpeakers() []speaker.Speaker {
	uri, _ := url.JoinPath(e.GetHost(), "speakers")
	body, err := util.HttpGet(uri)
	if err != nil {
		panic(err)
	}

	var speakers []speaker.Speaker
	if err := json.Unmarshal(body, &speakers); err != nil {
		panic(err)
	}
	return speakers
}
