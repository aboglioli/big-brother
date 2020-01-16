package mocks

import (
	"time"

	"github.com/stretchr/testify/mock"
)

type MockCache struct {
	mock.Mock
}

func NewMockCache() *MockCache {
	return &MockCache{}
}

func (c *MockCache) Get(k string) (interface{}, error) {
	args := c.Called(k)
	return args.Get(0), args.Error(1)
}

func (c *MockCache) Set(k string, v interface{}, d time.Duration) error {
	args := c.Called(k, v)
	return args.Error(0)
}

func (c *MockCache) Delete(k string) error {
	args := c.Called(k)
	return args.Error(0)
}
