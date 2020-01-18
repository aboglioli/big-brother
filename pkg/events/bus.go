package events

type Message interface {
	Body() []byte
	Event() Event
	Ack()
}

type Bus interface {
	Publish(body interface{}, topic, route string) error
	Consume(topic, route, queue string) (<-chan Message, error)
}
