package gelery

import (
	"github.com/google/uuid"
)

type CeleryBroker interface {
	SendCeleryMessage(*CeleryMessage, string) error
	GetAsyncResult(*uuid.UUID, string) (*AsyncResult, error)
}

type CeleryBrokerInfo struct {
	Host string
	Port string
	Type string
}
