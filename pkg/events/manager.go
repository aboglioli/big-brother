package events

type Options struct {
	Exchange string
	Route    string
	Queue    string
}

type Message interface {
	Body() []byte
	Event() Event
	Ack()
}

type Manager interface {
	Publish(body interface{}, opts *Options) error
	Consume(opts *Options) (<-chan Message, error)
}
