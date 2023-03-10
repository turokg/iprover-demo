package internal

import "context"

type Logger interface {
	Info(ctx context.Context, msg string)
	Error(ctx context.Context, msg string, err error)
	WithField(key string, value interface{}) Logger
}
