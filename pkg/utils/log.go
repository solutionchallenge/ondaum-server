package utils

import (
	"context"
	"fmt"
	"path"

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

func (logger Logger) CID(cid string) Logger {
	logger.Instance = logger.Instance.Str("correlation_id", cid)
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

func (logger Logger) BT(depth ...int) Logger {
	skip := 1
	if len(depth) > 0 && depth[0] > 0 {
		skip += depth[0]
	}
	filepath, filename, line := GetCallerInfo(skip)
	info := fmt.Sprintf("%s:%d", path.Join(filepath, filename), line)
	logger.Instance = logger.Instance.Str("caller", info)
	return logger
}

func (logger Logger) Send(fmtstr string, data ...any) {
	logger.Instance.Msgf(fmtstr, data...)
}
