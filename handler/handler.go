package handler

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type Publisher interface {
	Register(sub Subscriber) error
	Deregister(sub Subscriber) error
	Notify(messages <-chan *Message, sucess <-chan *Message)
}

type Subscriber interface {
	HandleEvent(data *Event) error
}

type Resource interface {
	Subscriber
	GetEventType() string
}

type Handler interface {
	Publisher
	Run()
}

type Event struct {
	Version    string      `json:"version"`
	ID         string      `json:"id"`
	DetailType string      `json:"detail-type"`
	Source     string      `json:"source"`
	Account    string      `json:"account"`
	Time       time.Time   `json:"time"`
	Region     string      `json:"region"`
	Resources  []any       `json:"resources"`
	Detail     interface{} `json:"detail"`
}

type Message struct {
	Event        *Event
	QueueMessage *types.Message
}
