package internal

import (
	"context"
	log "github.com/sirupsen/logrus"
	"os"
)

func NewLogger() Logger {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	log.SetLevel(log.InfoLevel)
	return &logger{}
}

type logger struct {
	entry *log.Entry
}

func (l *logger) Info(_ context.Context, msg string) {
	if l.entry != nil {
		l.entry.Info(msg)
		return
	}
	log.Info(msg)
}

func (l *logger) Error(_ context.Context, msg string, err error) {
	if l.entry != nil {
		l.entry.WithError(err).Error(msg)
		return
	}
	log.WithError(err).Error(msg)
}

func (l *logger) WithField(key string, value interface{}) Logger {
	l.entry = log.WithField(key, value)
	return l
}
