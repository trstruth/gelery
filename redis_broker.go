package gelery

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type RedisBroker struct {
	redisClient        *redis.Client
	redisResultsClient *redis.Client
}

func NewRedisBroker(brokerInfo *CeleryBrokerInfo) *RedisBroker {
	return &RedisBroker{
		redisClient: redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:%s", brokerInfo.Host, brokerInfo.Port),
			DB:   0,
		}),
		redisResultsClient: redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:%s", brokerInfo.Host, brokerInfo.Port),
			DB:   1,
		}),
	}
}

func (rb *RedisBroker) SendCeleryMessage(message *CeleryMessage, queue string) error {
	jsonBytes, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("Failed to marshal message: %s", err)
	}

	ctx := context.Background()
	_, err = rb.redisClient.LPush(ctx, queue, jsonBytes).Result()
	if err != nil {
		return fmt.Errorf("Failed to write to redis: %s", err)
	}

	return nil
}

func (rb *RedisBroker) GetAsyncResult(id *uuid.UUID, queue string) (*AsyncResult, error) {
	ctx := context.Background()
	data, err := rb.redisResultsClient.Get(ctx, fmt.Sprintf("celery-task-meta-%s", id)).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("Task with id %s doesn't exist", id)
	} else if err != nil {
		return nil, fmt.Errorf("Failed to get async result: %s", err)
	}

	asyncRes := &AsyncResult{}
	if err := json.Unmarshal([]byte(data), asyncRes); err != nil {
		return nil, fmt.Errorf("Failed to unmarshal json: %s", err)
	}

	return asyncRes, nil
}
