package sqs

import (
	"weather-report/aws-resource-handler/cloud/aws"
	"errors"
	"os"
	"context"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
)

const AWS_SQS_URL_ENV_NAME = "AWS_SQS_URL"

func getClient() (*sqs.Client, error) {
	cfg, err := aws.Config()
	if err != nil {
		return nil, err
	}
	return sqs.NewFromConfig(*cfg), nil
}

func GetMessages() (*[]types.Message, error) {
	queueUrl, err := getQueueUrl()
	if err != nil {
		return nil, err
	}
	sqsRequest := sqs.ReceiveMessageInput{
		QueueUrl: &queueUrl,
		MaxNumberOfMessages: 10,
	}
	sqsClient, err := getClient()
	if err != nil {
		return nil, err
	}
	sqsResponse, err := sqsClient.ReceiveMessage(context.TODO(),&sqsRequest)
	if err != nil {
		return nil, err
	}
	return &sqsResponse.Messages, nil
}

func DeleteMessage(recieptHandler *string) error {
	sqsClient, err := getClient()
	if err != nil {
		return err
	}
	queueUrl, err := getQueueUrl()
	if err != nil {
		return err
	}
	deleteMessageInput := sqs.DeleteMessageInput{
		QueueUrl: &queueUrl,
		ReceiptHandle: recieptHandler,
	}
	_, err = sqsClient.DeleteMessage(context.TODO(), &deleteMessageInput)
	if err != nil {
		return err
	}
	return nil
}

func getQueueUrl() (string, error){
	queueUrl, ok := os.LookupEnv("AWS_SQS_URL")
	if ok {
		return queueUrl, nil
	} else {
		return "", errors.New("unable to find the sqs url in the env variables")
	}
}
