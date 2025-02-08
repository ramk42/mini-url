package logging

import (
	"github.com/rs/zerolog"
	"os"
)

var Logger zerolog.Logger

func InitLogger() {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	Logger = zerolog.New(os.Stderr).
		With().
		Timestamp().
		Str("service", "mini-url-service").
		Logger()
}
