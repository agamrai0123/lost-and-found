package pkg

import (
	"sync"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"gopkg.in/natefinch/lumberjack.v2"
)

var oncelog sync.Once

func GetLogger() zerolog.Logger {
	oncelog.Do(func() {
		zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack
		zerolog.TimeFieldFormat = time.RFC3339Nano

		rotatingLog := &lumberjack.Logger{
			Filename:   AppConfig.Logging.Path,
			MaxSize:    AppConfig.Logging.MaxSize,
			MaxBackups: 10,
			MaxAge:     14, //days
			Compress:   true,
		}

		logger := zerolog.New(rotatingLog).
			Level(zerolog.Level(AppConfig.Logging.Level)).
			With().
			Timestamp().
			Str("service", "validate_session").
			Logger()

		log.Logger = logger
		log.Info().Msg("logger started for validate_session")
	})

	return log.Logger
}
