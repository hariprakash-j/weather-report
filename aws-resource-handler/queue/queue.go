package queue

import (
	"cloud-resource-scheduler/aws-resource-handler/cloud/aws/sqs"
	"fmt"
	"log/slog"
)
func PrintMessages() {
	for {
		messages, err := sqs.GetMessages()
		if err != nil {
			slog.Error("unable to get messages: ", err)
		}
		if len(*messages) > 0 {
			for _, message := range *messages {
				fmt.Println(*message.Body)
				err := sqs.DeleteMessage(message.ReceiptHandle)
				if err != nil {
					slog.Error("unable to delete the message: ", err)
				}
			}
		}
	}
}
