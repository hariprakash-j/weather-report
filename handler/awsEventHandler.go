package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"reflect"
	"syscall"
	"time"
	"weather-report/cloud/aws/sqs"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type AwsEventHandler struct {
	subscribers         []Resource
	MaxProcessorThreads int
}

func (e *AwsEventHandler) fetchMessages(c <-chan os.Signal) {
	messagesChannel := make(chan types.Message, 20)
	for i := 0; i < e.MaxProcessorThreads; i++ {
		go e.Notify(messagesChannel)
	}

	slog.Info("ready to recieve messages")

	for {
		select {
		case <-c:
			slog.Info("recieved kill or interupt signal")
			slog.Info("waiting for messages in flight to process...")
			for len(messagesChannel) > 0 {
				time.Sleep(200 * time.Millisecond)
				slog.Info("waiting..")
			}
			slog.Info("processed messages in flight")
			defer close(messagesChannel)
			slog.Info("main process ended")
			return
		default:
			messages, err := sqs.GetMessages()
			if err != nil {
				slog.Error("unable to get messages: ", err)
			}
			if len(*messages) > 0 {
				for _, message := range *messages {
					messagesChannel <- message
				}
			}
		}
	}
}

func (e *AwsEventHandler) Register(sub Resource) error {
	for _, s := range e.subscribers {
		if reflect.DeepEqual(sub, s) {
			return fmt.Errorf("subscriber already exists")
		}
	}
	e.subscribers = append(e.subscribers, sub)
	return nil
}

func (e *AwsEventHandler) DeRegister(sub Subscriber) error {
	for i, s := range e.subscribers {
		if reflect.DeepEqual(sub, s) {
			e.subscribers[i] = e.subscribers[len(e.subscribers)-1]
			e.subscribers = e.subscribers[:len(e.subscribers)-1]
			return nil
		}
	}
	return fmt.Errorf("subscriber not found")
}

func (e *AwsEventHandler) Notify(messages <-chan types.Message) {
	for {
		select {
		case message := <-messages:
			e.filterEvents(message.Body)
			err := sqs.DeleteMessage(message.ReceiptHandle)
			if err != nil {
				slog.Error("unable to delete the message: ", err)
			}
		default:
		}
	}
}

func (e *AwsEventHandler) filterEvents(message *string) error {
	var event Ec2Event
	err := json.Unmarshal([]byte(*message), &event)
	if err != nil {
		return err
	}
	for _, sub := range e.subscribers {
		if event.Source == sub.GetEventType() {
			sub.HandleEvent(&event)
			return nil
		}
	}
	return fmt.Errorf("unable to match the event")
}

func (e *AwsEventHandler) Run() {
	slog.Info("setting up kill and interupt signal handlers...")
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	slog.Info("created kill and interupt signal handlers")
	slog.Info("starting the main process...")
	e.fetchMessages(c)
	defer close(c)
}
