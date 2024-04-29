package handler

type Publisher interface {
	Register(sub Subscriber) error
	Deregister(sub Subscriber) error
	Notify(messages <-chan interface{})
}

type Subscriber interface {
	HandleEvent(data interface{}) error
}

type Resource interface {
	Subscriber
	GetEventType() string
}

type Handler interface {
	Publisher
	Run()
}
