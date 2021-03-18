package gelery

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
)

type RedisBroker struct {
	redisClient *redis.Client
}

func NewRedisBroker(brokerInfo *CeleryBrokerInfo) *RedisBroker {
	return &RedisBroker{
		redisClient: redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("%s:%s", brokerInfo.Host, brokerInfo.Port),
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

func (rb *RedisBroker) GetAsyncResult(id uuid.UUID) (*AsyncResult, error) {
	return nil, nil
}
