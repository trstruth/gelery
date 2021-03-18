package gelery

import (
	"fmt"

	"github.com/google/uuid"
)

// CeleryClient allows callers to create tasks and get task results
type CeleryClient struct {
	broker CeleryBroker
}

func NewCeleryClient(opts ...Option) (*CeleryClient, error) {
	cc := &CeleryClient{}

	for _, o := range opts {
		err := o(cc)
		if err != nil {
			return nil, fmt.Errorf("Failed to create CeleryClient: %s", err)
		}
	}

	return cc, nil
}

type Option func(cc *CeleryClient) error

func WithBroker(brokerInfo *CeleryBrokerInfo) Option {
	return func(cc *CeleryClient) error {
		switch brokerInfo.Type {
		case "redis":
			cc.broker = NewRedisBroker(brokerInfo)
		default:
			return fmt.Errorf("Unsupported broker type: %s", brokerInfo.Type)
		}

		return nil
	}
}

func (cc *CeleryClient) SendTask(taskName string, args []interface{}, kwargs map[string]interface{}, queue string) (*uuid.UUID, error) {
	sendTaskMessage := NewSendTaskMessage(taskName, args, kwargs)
	err := cc.broker.SendCeleryMessage(sendTaskMessage, queue)
	if err != nil {
		return nil, fmt.Errorf("Failed to SendTask %s: %s", taskName, err)
	}
	return sendTaskMessage.Headers.ID, nil
}

func (cc *CeleryClient) GetResult(id *uuid.UUID, queue string) (*AsyncResult, error) {
	return cc.broker.GetAsyncResult(id, queue)
}
