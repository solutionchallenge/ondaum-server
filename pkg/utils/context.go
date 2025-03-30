package utils

import (
	"context"
	"time"
)

const (
	CtxKeyForRequestID = "ctxval_request_id"
)

func SleepWith(ctx context.Context, delay time.Duration) error {
	timer := time.NewTimer(delay)
	defer timer.Stop()

	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
		return nil
	}
}

func WithValue[K any, V any](ctx context.Context, key K, value V) context.Context {
	return context.WithValue(ctx, key, value)
}

func GetValue[K any, V any](ctx context.Context, key K) (V, bool) {
	value, ok := ctx.Value(key).(V)
	return value, ok
}

func GetRequestID(ctx context.Context, defval ...string) string {
	value, ok := GetValue[string, string](ctx, CtxKeyForRequestID)
	if !ok {
		if len(defval) > 0 && defval[0] != "" {
			return defval[0]
		}
		return ""
	}
	return value
}
