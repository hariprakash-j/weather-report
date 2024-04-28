package main

import "weather-report/handler"

func main() {
	handler := handler.EventHandler{
		MaxProcessorThreads: 10,
	}
	handler.Run()
}
