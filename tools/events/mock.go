package events

import (
	"github.com/aboglioli/big-brother/events"
	"github.com/aboglioli/big-brother/tools/mock"
)

// MEssage
type mockMessage struct{}

func (m *mockMessage) Body() []byte {
	return []byte{}
}

func (m *mockMessage) Event() events.Event {
	return events.Event{}
}

func (m *mockMessage) Ack() {}

// Manager
type MockEventManager struct {
	Mock mock.Mock
}

func NewMockManager() *MockEventManager {
	return &MockEventManager{}
}

func (m *MockEventManager) Publish(body interface{}, opts *events.Options) error {
	m.Mock.Called(mock.Call("Publish", body, opts).Return(nil))
	return nil
}

func (m *MockEventManager) Consume(opts *events.Options) (<-chan events.Message, error) {
	call := mock.Call("Consume", opts)
	c := make(chan events.Message)
	m.Mock.Called(call.Return(c, nil))
	return c, nil
}
