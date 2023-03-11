package internal

import "time"

const (
	StatusError = 500
)

type Problem struct {
	ID         string
	Filename   string
	Text       string
	UploadedAt time.Time
}
