package users

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
)

func TestPasswordBcrypt(t *testing.T) {
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
		crypt := &bcryptCrypt{bcrypt.MinCost}

		hash, err := crypt.Hash(test.pwd)
		assert.Nil(err)
		assert.Greater(len(hash), 10)

		assert.Equal(crypt.Compare(hash, test.expected), test.shouldBeEqual)
	}
}
