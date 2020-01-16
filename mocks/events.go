package mocks

import (
	"github.com/aboglioli/big-brother/pkg/events"
	"github.com/stretchr/testify/mock"
)

type MockEventManager struct {
	mock.Mock
}

func NewMockEventManager() *MockEventManager {
	return &MockEventManager{}
}

func (m *MockEventManager) Publish(body interface{}, opts *events.Options) error {
	args := m.Called(body, opts)
	return args.Error(0)
}

func (m *MockEventManager) Consume(opts *events.Options) (<-chan events.Message, error) {
	args := m.Called(opts)
	msg, ok := args.Get(0).(<-chan events.Message)
	if !ok {
		return nil, args.Error(1)
	}
	return msg, args.Error(1)
}
