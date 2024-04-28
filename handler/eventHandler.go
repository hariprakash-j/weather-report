package handler

import (
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

type EventHandler struct {
	subscribers         []Subscriber
	MaxProcessorThreads int
}

func (e *EventHandler) fetchMessages(c <-chan os.Signal) {
	messagesChannel := make(chan types.Message, 20)
	for i := 0; i < e.MaxProcessorThreads; i++ {
		go e.Notify(messagesChannel)
	}

	for {
		select {
		case <-c:
			slog.Info("waiting for messages in flight to process...")
			for len(messagesChannel) > 0 {
				time.Sleep(200 * time.Millisecond)
				slog.Info("waiting..")
			}
			slog.Info("done")
			defer close(messagesChannel)
			slog.Info("message handler process ended")
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

func (e *EventHandler) Register(sub Subscriber) error {
	for _, s := range e.subscribers {
		if reflect.DeepEqual(sub, s) {
			return fmt.Errorf("subscriber already exists")
		}
	}
	e.subscribers = append(e.subscribers, sub)
	return nil
}

func (e *EventHandler) DeRegister(sub Subscriber) error {
	for i, s := range e.subscribers {
		if reflect.DeepEqual(sub, s) {
			e.subscribers[i] = e.subscribers[len(e.subscribers)-1]
			e.subscribers = e.subscribers[:len(e.subscribers)-1]
			return nil
		}
	}
	return fmt.Errorf("subscriber not found")
}

func (e *EventHandler) Notify(messages <-chan types.Message) {
	for {
		select {
		case message := <-messages:
			fmt.Println(*message.Body)
			err := sqs.DeleteMessage(message.ReceiptHandle)
			if err != nil {
				slog.Error("unable to delete the message: ", err)
			}
		default:
		}
	}
}

func (e *EventHandler) Run() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	e.fetchMessages(c)
	defer close(c)
}
