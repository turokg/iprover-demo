package internal

import "time"

type Problem struct {
	Id         string
	Filename   string
	Text       string
	UploadedAt time.Time
}
