package worker

import (
	"context"

	"github.com/hibiken/asynq"
)

type TaskDistributor interface {
	DistributeTaskSendVerifyEmail(ctx context.Context, payload *PayloadSendVerifyEmail, opts ...asynq.Option) error
}

type TaskDistributorRedis struct {
	client *asynq.Client
}

func NewTaskDistributorRedis(opt asynq.RedisClientOpt) TaskDistributor {
	client := asynq.NewClient(opt)

	return &TaskDistributorRedis{
		client: client,
	}
}
