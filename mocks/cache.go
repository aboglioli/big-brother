package mocks

import (
	"time"

	"github.com/stretchr/testify/mock"
)

type MockCache struct {
	mock.Mock
}

func NewMockCache(ns string) *MockCache {
	return &MockCache{}
}

func (c *MockCache) Get(k string) interface{} {
	args := c.Called(k)
	return args.Get(0)
}

func (c *MockCache) Set(k string, v interface{}, d time.Duration) {
	c.Called(k, v)
}

func (c *MockCache) Delete(k string) {
	c.Called(k)
}
