package future

import (
	"context"
	"time"
)

type Transaction interface {
	Complete(ctx context.Context) error
	Fail(ctx context.Context, errorMessage string) error
}

type Core interface {
	Create(ctx context.Context, actionType JobType, actionParams string, triggerAfter time.Duration, actionIdentifier ...string) (*Job, error)
	Update(ctx context.Context, ID string, actionType JobType, actionParams string, triggerAfter ...time.Duration) error
	Cancel(ctx context.Context, ID string) error
	Reschdule(ctx context.Context, ID string, triggerAfter time.Duration, onlyPending bool) error
	Inspect(ctx context.Context, ID string) (*Job, error)
	FindBy(ctx context.Context, actionIdentifier string) (*Job, error)
	RunNext(ctx context.Context, ignoreTriggerAfter bool) (*Job, Transaction, error)
	DeletePermanently(ctx context.Context, ID string) error
}
