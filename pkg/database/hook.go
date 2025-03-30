package database

import (
	"context"
	"time"

	"github.com/lithammer/shortuuid"
	sqlhook "github.com/qustavo/sqlhooks/v2"
	"github.com/solutionchallenge/ondaum-server/pkg/utils"
)

type hookContextKey string

type QueryLoggingHook struct {
	LogLevel utils.LogLevel
}

var _ sqlhook.Hooks = (*QueryLoggingHook)(nil)

const (
	dbQueryLoggingTimestampCtxKey  hookContextKey = "db-query-log-timestamp"
	dbQueryLoggingIdentifierCtxKey hookContextKey = "db-query-log-uuid"
	dbQueryLoggingPlaceholder      string         = "-----"
)

func NewQueryLoggingHook(logLevel utils.LogLevel) QueryLoggingHook {
	return QueryLoggingHook{
		LogLevel: logLevel,
	}
}

func (hook *QueryLoggingHook) Before(ctx context.Context, query string, args ...any) (context.Context, error) {
	queryUUID := shortuuid.New()

	requestID := utils.GetRequestID(ctx)
	utils.Log(hook.LogLevel).Ctx(ctx).RID(requestID).Send("[%s:%s]> %s %+v", queryUUID, dbQueryLoggingPlaceholder, query, args)

	wrappedContext := utils.WithValue(ctx, dbQueryLoggingTimestampCtxKey, time.Now())
	wrappedContext = utils.WithValue(wrappedContext, dbQueryLoggingIdentifierCtxKey, queryUUID)
	return wrappedContext, nil
}

func (hook *QueryLoggingHook) After(ctx context.Context, query string, args ...any) (context.Context, error) {
	queryUUID, ok := utils.GetValue[hookContextKey, string](ctx, dbQueryLoggingIdentifierCtxKey)
	if !ok {
		return ctx, nil
	}

	beginTimestamp, ok := utils.GetValue[hookContextKey, time.Time](ctx, dbQueryLoggingTimestampCtxKey)
	if !ok {
		return ctx, nil
	}

	requestID := utils.GetRequestID(ctx)
	utils.Log(hook.LogLevel).Ctx(ctx).RID(requestID).Send("[%s:%s]< %s %+v", queryUUID, time.Since(beginTimestamp), query, args)
	return ctx, nil
}
