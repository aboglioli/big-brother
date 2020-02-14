package mocks

import (
	"github.com/aboglioli/big-brother/pkg/events"
	"github.com/stretchr/testify/mock"
)

type MockEventManager struct {
	mock.Mock
}

func NewMockEventManager() *MockEventManager {
	return new(MockEventManager)
}

func (m *MockEventManager) Publish(body interface{}, topic, route string) error {
	args := m.Called(body, topic, route)
	return args.Error(0)
}

func (m *MockEventManager) Consume(topic, route, queue string) (<-chan events.Message, error) {
	args := m.Called(topic, route, queue)
	msg, ok := args.Get(0).(<-chan events.Message)
	if !ok {
		return nil, args.Error(1)
	}
	return msg, args.Error(1)
}
