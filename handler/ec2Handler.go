package handler

import (
	"fmt"
)

type Ec2Handler struct {
	EventType string
}

func (e *Ec2Handler) HandleEvent(event *Event) error {
	fmt.Printf("recieved EC2 data: %v\n", event.Detail)
	return nil
}

func (e *Ec2Handler) GetEventType() string {
	return e.EventType
}
