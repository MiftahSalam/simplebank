package worker

import (
	"context"
	db "simplebank/db/sqlc"

	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
)

type TaskProcessor interface {
	Start() error
	ProcessTaskSendVerifyEmail(ctx context.Context, task *asynq.Task) error
}

type TaskProcessorRedis struct {
	server *asynq.Server
	store  db.Store
}

func NewTaskProcessorRedis(opt *asynq.RedisClientOpt, store db.Store) TaskProcessor {
	server := asynq.NewServer(opt, asynq.Config{
		ErrorHandler: asynq.ErrorHandlerFunc(func(ctx context.Context, task *asynq.Task, err error) {
			log.Error().Err(err).Str("type", task.Type()).Bytes("payload", task.Payload()).Msg("process task failed")
		}),
		Logger: NewLogger(),
	})

	return &TaskProcessorRedis{
		server: server,
		store:  store,
	}
}

// Start implements TaskProcessor.
func (processor *TaskProcessorRedis) Start() error {
	mux := asynq.NewServeMux()

	mux.HandleFunc(TaskSendVerifyEmail, processor.ProcessTaskSendVerifyEmail)

	return processor.server.Start(mux)
}
