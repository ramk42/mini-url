// Package logger provides a global logger instance that can be used throughout the application.
package logger

import (
	"github.com/rs/zerolog"
	"os"
	"sync"
)

var (
	once sync.Once
	log  zerolog.Logger
)

func Init(serviceName string) {
	once.Do(func() {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		log = zerolog.New(os.Stderr).
			With().
			Timestamp().
			Str("service", serviceName).
			Logger().
			Level(zerolog.InfoLevel)
	})
}

func Instance() zerolog.Logger {
	return log
}
