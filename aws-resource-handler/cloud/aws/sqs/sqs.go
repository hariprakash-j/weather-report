package sqs

import (
	"cloud-resource-scheduler/aws-resource-handler/cloud/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
)



func getClient() (*sqs.Client, error) {
	cfg, err := aws.Config()
	if err != nil {
		return nil, err
	}
	return sqs.NewFromConfig(*cfg), nil
}
