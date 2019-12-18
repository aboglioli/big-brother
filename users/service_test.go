package users

import (
	"reflect"
	"testing"

	"github.com/aboglioli/big-brother/errors"
	"github.com/aboglioli/big-brother/tools/events"
)

func TestGetByID(t *testing.T) {
	repo, events := newMockRepository(), events.NewMockManager()
	user1 := newMockUser()
	user2 := newMockUser()
	user2.Validated = true

	tests := []struct {
		in   string
		out  error
		user *User
	}{{
		"123",
		ErrNotFound,
		user1,
	}, {
		"abc1235",
		ErrNotFound,
		user1,
	}, {
		"xyz123",
		ErrNotFound,
		user1,
	}, {
		user1.ID.Hex(),
		ErrNotValidated,
		user1,
	}, {
		user2.ID.Hex(),
		nil,
		user2,
	}}

	for i, test := range tests {
		repo.Collection = []*User{test.user}
		serv := NewService(repo, events)
		u, err := serv.GetByID(test.in)

		if test.out != nil && err != nil {
			expectedErr := test.out.(errors.Error)
			err := err.(errors.Error)
			if expectedErr.Code != err.Code {
				t.Errorf("test %d:\n-expected:%#v\n-actual:  %#v", i, test.out, err)
			}
		}

		if test.out == nil && err != nil {
			t.Errorf("test %d:\n-expected: nil error\n-actual:%#v", i, err)
			break
		}

		if err == nil && !reflect.DeepEqual(u, test.user) {
			t.Errorf("test %d:\n-expected:%#v\n-actual:  %#v", i, test.user, u)
		}
	}
}
