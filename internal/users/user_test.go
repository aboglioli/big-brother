package users

import (
	"testing"
)

func TestUserSetPassword(t *testing.T) {
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

	for i, test := range tests {
		user := NewUser()
		user.SetPassword(test.in)

		if len(user.Password) < 10 {
			t.Errorf("weak hashing")
		}

		if user.ComparePassword(test.out) != test.shouldBeEqual {
			t.Errorf("test %d: password %s comparison with hash %s should be %v", i, test.out, user.Password, test.shouldBeEqual)
		}
	}
}
