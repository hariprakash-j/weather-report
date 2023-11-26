package queue

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

type sqsQueue struct {
	awsClient aws.Config
}

func (s *sqsQueue)configureClient() {
	var err error;
	s.awsClient, err = config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic(err)
	}
}
