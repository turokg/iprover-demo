package internal

import (
	"time"
)

const (
	StatusError = 500
)

type Problem struct {
	ID         string
	Filename   string
	Text       string
	UploadedAt time.Time
}

type LaunchArgs struct {
	ProblemID   string
	ProblemText string
	Params      map[string]string
}

// Source сообщения
const (
	Process = "Process"
	System  = "System"
)

type LogMessage struct {
	Datetime time.Time `json:"time"`
	Source   string    `json:"source"`
	Message  string    `json:"text"`
}
