package worker

import (
	"context"
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type Logger struct{}

func NewLogger() *Logger {
	return &Logger{}
}

func (logger *Logger) print(level zerolog.Level, args ...interface{}) {
	log.
		WithLevel(level).
		Msg(fmt.Sprint(args...))
}

func (logger *Logger) Printf(ctx context.Context, format string, v ...interface{}) {
	log.WithLevel(zerolog.DebugLevel).Msgf(format, v...)
}

func (logger *Logger) Debug(args ...interface{}) {
	logger.print(zerolog.DebugLevel, args...)
}

func (logger *Logger) Info(args ...interface{}) {
	logger.print(zerolog.InfoLevel, args...)
}

func (logger *Logger) Warn(args ...interface{}) {
	logger.print(zerolog.WarnLevel, args...)
}

func (logger *Logger) Error(args ...interface{}) {
	logger.print(zerolog.ErrorLevel, args...)
}

func (logger *Logger) Fatal(args ...interface{}) {
	logger.print(zerolog.FatalLevel, args...)
}
