package database

import (
	"context"

	"github.com/benbjohnson/clock"
	"github.com/solutionchallenge/ondaum-server/pkg/future"
	"github.com/solutionchallenge/ondaum-server/pkg/utils"
	"github.com/uptrace/bun"
)

type Transaction struct {
	DB    *bun.DB
	Job   *FutureJob
	Clock clock.Clock
}

func (t *Transaction) Complete(ctx context.Context) error {
	_, err := t.DB.NewUpdate().Model(t.Job).
		Set("status = ?", future.JobStatusCompleted).
		Set("completed_at = ?", t.Clock.Now().UTC()).
		Where("id = ?", t.Job.ID).
		Exec(ctx)
	if err != nil {
		return utils.WrapError(err, "failed to complete future job")
	}
	return nil
}

func (t *Transaction) Fail(ctx context.Context, errorMessage string) error {
	_, err := t.DB.NewUpdate().Model(t.Job).
		Set("status = ?", future.JobStatusFailed).
		Set("error = ?", errorMessage).
		Set("completed_at = ?", t.Clock.Now().UTC()).
		Where("id = ?", t.Job.ID).
		Exec(ctx)
	if err != nil {
		return utils.WrapError(err, "failed to fail future job")
	}
	return nil
}
