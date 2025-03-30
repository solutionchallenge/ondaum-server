package memdb

import (
	"context"
	"sync"

	sqlhook "github.com/qustavo/sqlhooks/v2"
)

// Ref. https://github.com/dolthub/go-memdb-server/issues/1306
type MemdbThreadSafetyHook struct {
	mutex sync.Mutex
}

var _ sqlhook.Hooks = (*MemdbThreadSafetyHook)(nil)

func NewMemdbThreadSafetyHook() MemdbThreadSafetyHook {
	return MemdbThreadSafetyHook{}
}

func (hook *MemdbThreadSafetyHook) Before(ctx context.Context, _ string, _ ...any) (context.Context, error) {
	hook.mutex.Lock()
	return ctx, nil
}

func (hook *MemdbThreadSafetyHook) After(ctx context.Context, _ string, _ ...any) (context.Context, error) {
	hook.mutex.Unlock()
	return ctx, nil
}
