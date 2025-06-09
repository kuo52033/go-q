package queue

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kuo-52033/go-q/internal/db"
)

type Producer struct {
	rdb *db.RedisClient
}

func NewProducer(rdb *db.RedisClient) *Producer {
	return &Producer{rdb: rdb}
}

func (p *Producer) Enqueue(
	jobType string,
	payload JobPayload,
	queueName string,
	maxAttempts int,
) (*Job, error) {

	job := Job{
		ID: JobId(uuid.New().String()),
		Type: jobType,
		Payload: payload,
		Status: StatusQueued,
		Queue: queueName,
		CreatedAt: time.Now(),
		AttemptCount: 0,
		MaxAttempts: maxAttempts,
	}

	key := fmt.Sprintf("job:%s", job.ID)

	err := p.rdb.Client.HSet(p.rdb.Ctx, key, job).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to set job in hash: %w", err)
	}

	listKey := fmt.Sprintf("queue:%s", queueName)	
	
	err = p.rdb.Client.LPush(p.rdb.Ctx, listKey, JobId(job.ID)).Err()
	if err != nil {
		return nil, fmt.Errorf("failed to push job to list: %w", err)
	}

	return &job, nil
}
