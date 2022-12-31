package launcher

import (
	"encoding/json"
	"time"
)

type Log struct {
	Datetime time.Time `json:"time"`
	Source   string    `json:"source"`
	Message  string    `json:"message"`
}

const (
	app = "app"
	sys = "sys"
)

func NewSysLog(message string) []byte {
	log := Log{
		Datetime: time.Now(),
		Source:   sys,
		Message:  message,
	}
	msg, _ := json.Marshal(log)
	return msg
}

func NewAppLog(message string) []byte {
	log := Log{
		Datetime: time.Now(),
		Source:   app,
		Message:  message,
	}
	msg, _ := json.Marshal(log)
	return msg
}
