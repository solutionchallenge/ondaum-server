package future

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/solutionchallenge/ondaum-server/pkg/utils"
)

type Scheduler struct {
	Core
	Config        Config
	Handlers      map[JobType]Handler
	QuitSignal    chan struct{}
	WaitGroup     sync.WaitGroup
	CancelableCtx context.Context
	CancelFunc    context.CancelFunc
}

func NewScheduler(config Config, core Core) *Scheduler {
	ctx, cancel := context.WithCancel(context.Background())
	return &Scheduler{
		Core:          core,
		Config:        config,
		Handlers:      make(map[JobType]Handler),
		QuitSignal:    make(chan struct{}),
		CancelableCtx: ctx,
		CancelFunc:    cancel,
	}
}

type Handler interface {
	Handle(ctx context.Context, job *Job) error
}

func (s *Scheduler) AddHandler(actionType JobType, handler Handler) {
	s.Handlers[actionType] = handler
}

func (s *Scheduler) Start() {
	s.WaitGroup.Add(1)
	go func() {
		defer s.WaitGroup.Done()
		for {
			select {
			case <-s.CancelableCtx.Done():
				utils.Log(utils.InfoLevel).BT().Send("Scheduler is shutting down...")
				return
			case <-time.After(s.Config.ScheduleCycle):
				s.processJobs()
			}
		}
	}()
}

func (s *Scheduler) Stop() {
	utils.Log(utils.InfoLevel).BT().Send("Initiating scheduler shutdown...")
	s.CancelFunc()
	s.WaitGroup.Wait()
	utils.Log(utils.InfoLevel).BT().Send("Scheduler shutdown completed")
}

func (s *Scheduler) processJobs() {
	for {
		select {
		case <-s.CancelableCtx.Done():
			utils.Log(utils.InfoLevel).BT().Send("Scheduler is stopping, exiting burst mode...")
			return
		default:
			job, tx, err := s.Core.RunNext(s.CancelableCtx, false)
			if err != nil {
				utils.Log(utils.WarnLevel).Err(err).BT().Send("Failed to get next job")
				return
			}
			if job == nil {
				return
			}
			if tx == nil {
				utils.Log(utils.WarnLevel).BT().Send("No transaction found for job: %s", job.ID)
				s.Core.Cancel(s.CancelableCtx, job.ID)
				continue
			}

			handler, ok := s.Handlers[job.ActionType]
			if !ok {
				utils.Log(utils.WarnLevel).BT().Send("Failed to get handler for action type: %s", job.ActionType)
				tx.Fail(s.CancelableCtx, "Failed to get handler for action type")
				continue
			}

			func() {
				defer func() {
					if r := recover(); r != nil {
						errMsg := "Handler panicked"
						if err, ok := r.(error); ok {
							errMsg = err.Error()
						}
						utils.Log(utils.ErrorLevel).Err(fmt.Errorf("%v", r)).BT().Send("Handler panicked for job: %s", job.ID)
						tx.Fail(s.CancelableCtx, errMsg)
					}
				}()

				err = handler.Handle(s.CancelableCtx, job)
				if err != nil {
					utils.Log(utils.WarnLevel).Err(err).BT().Send("Failed to handle job")
					tx.Fail(s.CancelableCtx, err.Error())
					return
				}
				tx.Complete(s.CancelableCtx)
			}()

			if s.Config.DeleteAfterCompletion {
				utils.Log(utils.DebugLevel).BT().Send("Job completed: %s", job.ID)
				s.Core.DeletePermanently(s.CancelableCtx, job.ID)
			}
		}
	}
}
