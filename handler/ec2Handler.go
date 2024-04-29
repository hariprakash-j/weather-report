package handler

import (
	"fmt"
	"time"
)

type Ec2Handler struct {
	EventType string
}

type Ec2Event struct {
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

func (e *Ec2Handler) HandleEvent(data interface{}) error {
	fmt.Printf("recieved EC2 data: %v", data)
	return nil
}

func (e *Ec2Handler) GetEventType() string {
	return e.EventType
}
