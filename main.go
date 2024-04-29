package main

import "weather-report/handler"

func main() {
	awsHandler := handler.AwsEventHandler{
		MaxProcessorThreads: 10,
	}
	// register handlers
	awsHandler.Register(
		&handler.Ec2Handler{EventType: "weather-report.test1"},
	)
	awsHandler.Run()
}
