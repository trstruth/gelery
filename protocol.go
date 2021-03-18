package gelery

import (
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// https://docs.celeryproject.org/en/stable/internals/protocol.html#version-2
type CeleryMessage struct {
	Body            string            `json:"body"`
	ContentType     string            `json:"content-type"`
	ContentEncoding string            `json:"content-encoding"`
	Headers         *CeleryHeaders    `json:"headers"`
	Properties      *CeleryProperties `json:"properties"`
}

func NewSendTaskMessage(taskName string, args []interface{}, kwargs map[string]interface{}) *CeleryMessage {
	id := uuid.New()
	replyTo := uuid.New()
	deliveryTag := uuid.New()

	body := []interface{}{args, kwargs, nil}
	jsonBodyData, _ := json.Marshal(body)
	b64Data := base64.StdEncoding.EncodeToString(jsonBodyData)

	return &CeleryMessage{
		Body:            string(b64Data),
		ContentType:     "application/json",
		ContentEncoding: "utf-8",
		Headers: &CeleryHeaders{
			Task:   taskName,
			ID:     &id,
			Lang:   "py",
			RootID: &id,
			TimingInfo: CeleryTimingInfo{
				BeforePublish: time.Now(),
			},
		},
		Properties: &CeleryProperties{
			ContentEncoding: "utf-8",
			CorrelationID:   &id,
			BodyEncoding:    "base64",
			ReplyTo:         &replyTo,
			DeliveryMode:    2,
			DeliveryTag:     &deliveryTag,
			DeliveryInfo: CeleryDeliveryInfo{
				Exchange:   "",
				RoutingKey: "heartbeat",
			},
		},
	}
}

type CeleryBody struct {
	Args   []interface{}
	Kwargs map[string]interface{}
	Embed  *CeleryBodyEmbed
}

func (cb *CeleryBody) MarshalJson() ([]byte, error) {
	return json.Marshal([]interface{}{cb.Args, cb.Kwargs, cb.Embed})
}

type CeleryBodyEmbed struct {
	Callbacks *[]string `json:"callbacks"`
	Errbacks  *[]string `json:"errbacks"`
	Chain     *[]string `json:"chain"`
	Chord     *[]string `json:"chord"`
}

type CeleryHeaders struct {
	Lang       string           `json:"lang"`
	Task       string           `json:"task"`
	ID         *uuid.UUID       `json:"id"`
	RootID     *uuid.UUID       `json:"root_id"`
	ParentID   *uuid.UUID       `json:"parent_id"`
	Group      *uuid.UUID       `json:"group"`
	TimingInfo CeleryTimingInfo `json:"timing_info"`
}

type CeleryTimingInfo struct {
	BeforePublish time.Time `json:"before_publish"`
}

type CeleryProperties struct {
	CorrelationID   *uuid.UUID         `json:"correlation_id"`
	ContentEncoding string             `json:"content_encoding"`
	ReplyTo         *uuid.UUID         `json:"reply_to"`
	DeliveryTag     *uuid.UUID         `json:"delivery_tag"`
	DeliveryMode    int                `json:"delivery_mode"`
	DeliveryInfo    CeleryDeliveryInfo `json:"delivery_info"`
	Priority        int                `json:"priority"`
	BodyEncoding    string             `json:"body_encoding"`
}

type CeleryDeliveryInfo struct {
	Exchange   string `json:"exchange"`
	RoutingKey string `json:"routing_key"`
}

type AsyncResult struct {
	TaskID    *uuid.UUID             `json: "task_id"`
	Status    string                 `json:"status"`
	Result    interface{}            `json:"result"`
	Traceback *string                `json:"traceback"`
	Name      string                 `json:"name"`
	Args      []interface{}          `json:"args"`
	Kwargs    map[string]interface{} `json:"kwargs"`
	Worker    string                 `json:"worker"`
	Retries   int                    `json:"retries"`
	Queue     string                 `json:"queue"`
}
