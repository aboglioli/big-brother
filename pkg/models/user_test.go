package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserSetPassword(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		pwd           string
		expected      string
		shouldBeEqual bool
	}{
		{"123456", "12345", false},
		{"12345", "123456", false},
		{"123456", "123456", true},
		{"123#!ł", "123#~ł€", false},
		{"123#~!ł€", "123#~!ł€", true},
	}

	for _, test := range tests {
		user := NewUser()
		user.SetPassword(test.pwd)
		assert.NotEmpty(user.Password)
		assert.Greater(len(user.Password), 10)
		assert.Equal(user.ComparePassword(test.expected), test.shouldBeEqual)
	}
}
