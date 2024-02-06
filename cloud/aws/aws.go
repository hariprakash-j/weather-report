package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

// generates the aws config, this functions relies on the env variables AWS_REGION for the region info and the aws credentials env vars if they are used or the AWS_PROFILE for the aws credentials
func Config() (*aws.Config, error) {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	return &cfg, err
}
