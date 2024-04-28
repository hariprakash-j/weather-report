package handler

type Publisher interface {
	Register(sub Subscriber) error
	Deregister(sub Subscriber) error
	Notify()
}

type Subscriber interface {
	HandleEvent() error
}

type Resource interface {
	Subscriber
	GetEventType() (string, error)
}

type Handler interface {
	Publisher
	Run()
}
