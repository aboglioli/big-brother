package events

type ManagerOptions struct {
	Exchange string
	Route string
	Queue string
}

type Message interface {
	Body() []byte
	Event() Event
	Ack()
}

type Manager interface {
	Publish(body interface{}, opts *ManagerOptions) error
	Consume(opts *ManagerOptions) (<-chan Message, error)
}
