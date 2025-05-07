package chat

import (
	"context"
	"encoding/json"
	"time"

	fthandler "github.com/solutionchallenge/ondaum-server/internal/handler/future"
	ftpkg "github.com/solutionchallenge/ondaum-server/pkg/future"
)

func UpsertChattingEndFutureJob(
	ctx context.Context, future *ftpkg.Scheduler, sessionID string, userID int64, triggerAfter time.Duration,
) (*ftpkg.Job, error) {
	job, err := future.FindBy(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if job == nil {
		marshaled, err := json.Marshal(fthandler.ChatFutureHandlerParams{
			UserID:         userID,
			ConversationID: sessionID,
		})
		if err != nil {
			return nil, err
		}
		job, err = future.Create(
			ctx, fthandler.ChatJobType, string(marshaled),
			triggerAfter,
			sessionID,
		)
		if err != nil {
			return nil, err
		}
		return job, nil
	} else {
		err = future.Reschdule(
			ctx, job.ID, triggerAfter, true,
		)
		if err != nil {
			return nil, err
		}
		return job, nil
	}
}
