package queue

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"
	"weather-report/aws-resource-handler/cloud/aws/sqs"

	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)
func fetchMessages(c <-chan os.Signal) {
	messagesChannel := make(chan types.Message, 20)
	for i := 0; i <5; i++ {
		go processMessage(messagesChannel)
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
				for _,message := range *messages {
					messagesChannel <- message
				}
			}
		}
	}

}

func processMessage(messages <-chan types.Message) {
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

func syncProcessor(c <-chan os.Signal) {
	for {
		select {
		case <-c:
			return
		default:
			messages, err := sqs.GetMessages()
			if err != nil {
				slog.Error("unable to get messages: ", err)
			}
			if len(*messages) > 0 {
				for _,message := range *messages {
					fmt.Println(*message.Body)
					err := sqs.DeleteMessage(message.ReceiptHandle)
					if err != nil {
						slog.Error("unable to process the message: ", err)
					}
				}
			}
		}
	}
}

			
func Run() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	fetchMessages(c)
	defer close(c)
}
