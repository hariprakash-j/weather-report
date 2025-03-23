package main

import (
	"runtime"
	"weather-report/cloud/aws/sqs"
	"weather-report/handler"
)

func main() {
	queue := sqs.Sqs{
		QueueUrl: "https://sqs.ap-south-1.amazonaws.com/999999999999/my-test-queue",
	}
	awsHandler := handler.AwsEventHandler{
		Queue:               &queue,
		MaxProcessorThreads: runtime.NumCPU(),
	}
	// register handlers
	awsHandler.Register(
		&handler.Ec2Handler{EventType: "weather-report.test1"},
	)
	awsHandler.Run()
}
