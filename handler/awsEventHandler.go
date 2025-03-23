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
	"weather-report/abstractions"

	sqsTypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type AwsEventHandler struct {
	Queue               abstractions.Queue[sqsTypes.Message, string, int32]
	Subscribers         []Resource
	MaxProcessorThreads int
}

func (e *AwsEventHandler) handleSucess(sucess <-chan *Message) {
	for {
		select {
		case message := <-sucess:
			slog.Info(
				fmt.Sprintf("deleting the message %v from the queue...", message.Event.Detail),
			)
			err := e.Queue.DeleteMessage(message.QueueMessage.ReceiptHandle)
			if err != nil {
				slog.Error(fmt.Sprintf("unable to delete the message: %s", err))
			}
			slog.Info(fmt.Sprintf("deleted the message %v from the queue", message.Event.Detail))
		default:
		}
	}
}

func (e *AwsEventHandler) processMessages(c <-chan os.Signal) {
	messagesChannel := make(chan *Message, 20)
	sucessChannel := make(chan *Message, 20)
	// failureChannel := make(chan *Message, 20)
	for i := 0; i < e.MaxProcessorThreads; i++ {
		go e.handleSucess(sucessChannel)
		go e.Notify(messagesChannel, sucessChannel)
		// go e.Notify(messagesChannel, sucessChannel, failureChannel)
	}

	slog.Info("ready to recieve messages")

	for {
		select {
		case <-c:
			slog.Info("recieved kill or interupt signal")
			slog.Info("waiting for messages in flight to process...")
			// for len(messagesChannel) > 0 || len(sucessChannel) > 0 || len(failureChannel) > 0 {
			for len(messagesChannel) > 0 || len(sucessChannel) > 0 {
				time.Sleep(200 * time.Millisecond)
				slog.Info("waiting..")
			}
			slog.Info("processed messages in flight")
			defer func() {
				close(messagesChannel)
				close(sucessChannel)
				// close(failureChannel)
			}()
			slog.Info("main process ended")
			return
		default:
			sqsMessages, err := e.Queue.GetMessages(10)
			if err != nil {
				slog.Error("unable to get messages: ", err)
			}
			if len(*sqsMessages) > 0 {
				for _, m := range *sqsMessages {
					event, err := e.unmarshalEvent(m.Body)
					message := Message{
						Event:        event,
						QueueMessage: &m,
					}
					if err != nil {
						continue
					}
					messagesChannel <- &message
				}
			}
		}
	}
}

func (e *AwsEventHandler) Register(sub Resource) error {
	for _, s := range e.Subscribers {
		if reflect.DeepEqual(sub, s) {
			return fmt.Errorf("subscriber already exists")
		}
	}
	e.Subscribers = append(e.Subscribers, sub)
	return nil
}

func (e *AwsEventHandler) DeRegister(sub Subscriber) error {
	for i, s := range e.Subscribers {
		if reflect.DeepEqual(sub, s) {
			e.Subscribers[i] = e.Subscribers[len(e.Subscribers)-1]
			e.Subscribers = e.Subscribers[:len(e.Subscribers)-1]
			return nil
		}
	}
	return fmt.Errorf("subscriber not found")
}

func (e *AwsEventHandler) Notify(
	message <-chan *Message,
	sucess chan<- *Message,
) {
	for {
		select {
		case message := <-message:
			err := e.sendEvents(message.Event)
			if err != nil {
				// failure <- message
				slog.Error("send event failed")
				continue
			}
			sucess <- message
		default:
		}
	}
}

func (e *AwsEventHandler) unmarshalEvent(message *string) (*Event, error) {
	var event Event
	err := json.Unmarshal([]byte(*message), &event)
	if err != nil {
		return &event, err
	}
	return &event, nil
}

func (e *AwsEventHandler) sendEvents(event *Event) error {
	for _, sub := range e.Subscribers {
		if event.Source == sub.GetEventType() {
			sub.HandleEvent(event)
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
	e.processMessages(c)
	defer close(c)
}
