package logger

import (
	"os"

	"github.com/go-kit/kit/log"
)

// Logger specifies logging API.
type Logger interface {
	// Info logs on info level.
	Info(string)
	// Error logs on error level.
	Error(string)
}

type logger struct {
	kitLogger log.Logger
}

func (l logger) Info(s string) {
	l.kitLogger.Log("level", "DEBUG", "message", s)
}

func (l logger) Error(s string) {
	l.kitLogger.Log("level", "ERROR", "message", s)
}

func New() Logger {
	l := log.NewJSONLogger(log.NewSyncWriter(os.Stdout))
	l = log.With(l, "ts", log.DefaultTimestampUTC)
	l = log.With(l, "caller", log.DefaultCaller)
	return logger{l}
}
