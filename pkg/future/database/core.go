package database

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"time"

	"github.com/benbjohnson/clock"
	"github.com/google/uuid"
	"github.com/solutionchallenge/ondaum-server/pkg/future"
	"github.com/uptrace/bun"
)

var _ future.Core = &Core{}

type Core struct {
	DB    *bun.DB
	Clock clock.Clock
}

func NewCore(db *bun.DB, clk clock.Clock) *Core {
	return &Core{DB: db, Clock: clk}
}

func (c *Core) Create(ctx context.Context, actionType future.JobType, actionParams string, triggerAfter time.Duration, actionIdentifier ...string) (*future.Job, error) {
	actionIdentifierValue := uuid.New().String()
	if len(actionIdentifier) > 0 && actionIdentifier[0] != "" {
		actionIdentifierValue = actionIdentifier[0]
	}

	job := &FutureJob{
		ActionType:       actionType,
		ActionParams:     actionParams,
		TriggeredAt:      c.Clock.Now().UTC().Add(triggerAfter),
		Status:           future.JobStatusPending,
		ActionIdentifier: actionIdentifierValue,
	}

	result, err := c.DB.NewInsert().Model(job).Exec(ctx)
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &future.Job{
		ID:               strconv.FormatInt(id, 10),
		ActionIdentifier: job.ActionIdentifier,
		ActionType:       job.ActionType,
		ActionParams:     job.ActionParams,
		TriggeredAt:      job.TriggeredAt,
		Status:           job.Status,
	}, nil
}

func (c *Core) Update(ctx context.Context, ID string, actionType future.JobType, actionParams string, triggerAfter ...time.Duration) error {
	update := c.DB.NewUpdate().Model((*FutureJob)(nil)).
		Set("action_type = ?", actionType).
		Set("action_params = ?", actionParams).
		Where("id = ?", ID)

	if len(triggerAfter) > 0 {
		update = update.Set("triggered_at = ?", c.Clock.Now().UTC().Add(triggerAfter[0]))
	}

	_, err := update.Exec(ctx)
	return err
}

func (c *Core) Cancel(ctx context.Context, ID string) error {
	_, err := c.DB.NewUpdate().Model((*FutureJob)(nil)).
		Set("status = ?", future.JobStatusCanceled).
		Set("completed_at = ?", c.Clock.Now().UTC()).
		Where("id = ?", ID).
		Exec(ctx)
	return err
}

func (c *Core) Reschdule(ctx context.Context, ID string, triggerAfter time.Duration, onlyPending bool) error {
	query := c.DB.NewUpdate().Model((*FutureJob)(nil)).
		Set("triggered_at = ?", c.Clock.Now().UTC().Add(triggerAfter))

	if onlyPending {
		query = query.Where("status = ?", future.JobStatusPending)
	} else {
		query = query.Set("status = ?", future.JobStatusPending)
	}

	_, err := query.Where("id = ?", ID).Exec(ctx)
	return err
}

func (c *Core) Inspect(ctx context.Context, ID string) (*future.Job, error) {
	var job FutureJob
	err := c.DB.NewSelect().Model(&job).Where("id = ?", ID).Scan(ctx)
	if err != nil {
		return nil, err
	}

	return &future.Job{
		ID:               strconv.FormatInt(job.ID, 10),
		ActionIdentifier: job.ActionIdentifier,
		ActionType:       job.ActionType,
		ActionParams:     job.ActionParams,
		TriggeredAt:      job.TriggeredAt,
		CompletedAt:      job.CompletedAt,
		Status:           job.Status,
		Error:            job.Error,
	}, nil
}

func (c *Core) FindBy(ctx context.Context, actionIdentifier string) (*future.Job, error) {
	var job FutureJob
	err := c.DB.NewSelect().Model(&job).Where("action_identifier = ?", actionIdentifier).Scan(ctx)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &future.Job{
		ID:               strconv.FormatInt(job.ID, 10),
		ActionIdentifier: job.ActionIdentifier,
		ActionType:       job.ActionType,
		ActionParams:     job.ActionParams,
	}, nil
}

func (c *Core) RunNext(ctx context.Context, ignoreTriggerAfter bool) (*future.Job, future.Transaction, error) {
	var job FutureJob

	tx, err := c.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, nil, err
	}

	query := tx.NewSelect().Model(&job).
		Where("status = ?", future.JobStatusPending).
		Order("triggered_at ASC").
		Limit(1)

	if !ignoreTriggerAfter {
		query = query.Where("triggered_at <= ?", c.Clock.Now().UTC())
	}

	err = query.Scan(ctx)
	if err != nil {
		tx.Rollback()
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil, nil
		}
		return nil, nil, err
	}

	_, err = tx.NewUpdate().Model(&job).
		Set("status = ?", future.JobStatusRunning).
		Where("id = ?", job.ID).
		Exec(ctx)
	if err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, nil, err
	}

	jobTx := &Transaction{
		DB:    c.DB,
		Job:   &job,
		Clock: c.Clock,
	}

	return &future.Job{
		ID:               strconv.FormatInt(job.ID, 10),
		ActionIdentifier: job.ActionIdentifier,
		ActionType:       job.ActionType,
		ActionParams:     job.ActionParams,
		TriggeredAt:      job.TriggeredAt,
		Status:           job.Status,
	}, jobTx, nil
}

func (c *Core) DeletePermanently(ctx context.Context, ID string) error {
	_, err := c.DB.NewDelete().Model((*FutureJob)(nil)).Where("id = ?", ID).Exec(ctx)
	return err
}
