package engine

import (
	"net/url"
	"strconv"
)

type Engine struct {
	Host string
	Port int
}

func (e *Engine) MyHost() string {
	myHost, err := url.Parse(e.Host)
	if err != nil {
		return ""
	}
	if myHost.Scheme == "" {
		myHost.Scheme = "http"
	}
	return myHost.String() + ":" + strconv.Itoa(e.Port)
}
