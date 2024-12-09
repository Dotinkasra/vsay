package engine

import (
	"strconv"
)

type Engine struct {
	Host string
	Port int
}

func (e *Engine) MyHost() string {
	return e.Host + ":" + strconv.Itoa(e.Port)
}
