package sqs

import (
	"context"
	"weather-report/cloud/aws"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

type Sqs struct {
	QueueUrl  string
	sqsClient *sqs.Client
}

func (a *Sqs) getClient() (*sqs.Client, error) {
	if a.sqsClient != nil {
		return a.sqsClient, nil
	}
	cfg, err := aws.Config()
	if err != nil {
		return nil, err
	}
	a.sqsClient = sqs.NewFromConfig(*cfg)
	return a.sqsClient, nil
}

func (a *Sqs) GetMessages(count int32) (*[]types.Message, error) {
	sqsRequest := sqs.ReceiveMessageInput{
		QueueUrl:            &a.QueueUrl,
		MaxNumberOfMessages: count,
	}
	sqsClient, err := a.getClient()
	if err != nil {
		return nil, err
	}
	sqsResponse, err := sqsClient.ReceiveMessage(context.TODO(), &sqsRequest)
	if err != nil {
		return nil, err
	}
	return &sqsResponse.Messages, nil
}

func (a *Sqs) DeleteMessage(recieptHandler *string) error {
	sqsClient, err := a.getClient()
	if err != nil {
		return err
	}
	deleteMessageInput := sqs.DeleteMessageInput{
		QueueUrl:      &a.QueueUrl,
		ReceiptHandle: recieptHandler,
	}
	_, err = sqsClient.DeleteMessage(context.TODO(), &deleteMessageInput)
	if err != nil {
		return err
	}
	return nil
}
