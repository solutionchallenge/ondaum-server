package utils

import (
	"context"
	"fmt"

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
				return fmt.Errorf("max retry backoff reached: %w", err)
			}
			if err = SleepWith(ctx, next); err != nil {
				return fmt.Errorf("sleep interrupted: %w", err)
			}
			continue
		}
		return err
	}
	return fmt.Errorf("max retry attempts reached: %w", err)
}
