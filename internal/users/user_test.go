package users

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserSetPassword(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		in            string
		out           string
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
		user.SetPassword(test.in)
		assert.NotEmpty(user.Password)
		assert.Greater(len(user.Password), 10)
		assert.Equal(user.ComparePassword(test.out), test.shouldBeEqual)
	}
}
