package utils

import (
	"context"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

type LogLevel zerolog.Level

const (
	FatalLevel = LogLevel(zerolog.FatalLevel)
	ErrorLevel = LogLevel(zerolog.ErrorLevel)
	WarnLevel  = LogLevel(zerolog.WarnLevel)
	InfoLevel  = LogLevel(zerolog.InfoLevel)
	DebugLevel = LogLevel(zerolog.DebugLevel)
	TraceLevel = LogLevel(zerolog.TraceLevel)
)

type Logger struct {
	Instance *zerolog.Event
}

func Log(lv LogLevel) Logger {
	return Logger{
		Instance: log.WithLevel(zerolog.Level(lv)),
	}
}

func (logger Logger) Ctx(ctx context.Context) Logger {
	logger.Instance = logger.Instance.Ctx(ctx)
	return logger
}

func (logger Logger) Err(err error) Logger {
	logger.Instance = logger.Instance.Err(err)
	return logger
}

func (logger Logger) RID(rid string) Logger {
	logger.Instance = logger.Instance.Str("request_id", rid)
	return logger
}

func (logger Logger) UID(uid string) Logger {
	logger.Instance = logger.Instance.Str("user_id", uid)
	return logger
}

func (logger Logger) Send(fmt string, data ...any) {
	logger.Instance.Msgf(fmt, data...)
}
