package mock

import (
	"github.com/aboglioli/big-brother/pkg/events"
)

// Message
type mockMessage struct{}

func (m *mockMessage) Body() []byte {
	return []byte{}
}

func (m *mockMessage) Event() events.Event {
	return events.Event{}
}

func (m *mockMessage) Ack() {}

// Manager
type EventManager struct {
	Mock Mock
}

func NewMockEventManager() *EventManager {
	return &EventManager{}
}

func (m *EventManager) Publish(body interface{}, opts *events.Options) error {
	m.Mock.Called(Call("Publish", body, opts).Return(nil))
	return nil
}

func (m *EventManager) Consume(opts *events.Options) (<-chan events.Message, error) {
	call := Call("Consume", opts)
	c := make(chan events.Message)
	m.Mock.Called(call.Return(c, nil))
	return c, nil
}
