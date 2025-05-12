package utils

import (
	"context"

	"github.com/cenkalti/backoff/v4"
)

func Retry(ctx context.Context, retryLimit int, action func() error, condition func(err error) bool) error {
	var err error
	backoffContext := backoff.WithContext(backoff.NewExponentialBackOff(), ctx)
	for retryCount := 0; retryLimit < 0 || retryCount <= retryLimit; retryCount++ {
		err = action()
		if err == nil {
			return nil
		}
		if condition != nil && condition(err) {
			next := backoffContext.NextBackOff()
			if next == backoff.Stop {
				return WrapError(err, "max retry backoff reached")
			}
			if err = SleepWith(ctx, next); err != nil {
				return WrapError(err, "sleep interrupted")
			}
			continue
		}
		return err
	}
	return WrapError(err, "max retry attempts reached")
}
