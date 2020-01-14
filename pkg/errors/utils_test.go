package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCompare(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		err1        error
		err2        error
		shouldEqual bool
	}{{
		errors.New("code1"),
		errors.New("code1"),
		true,
	}, {
		errors.New("code1"),
		errors.New("code2"),
		false,
	}, {
		Internal.New("internal.error_1"),
		Internal.New("internal.error_1"),
		true,
	}, {
		Internal.New("internal.error_1"),
		Internal.New("internal.error_2"),
		false,
	}, {
		Internal.New("internal.error"),
		Status.New("internal.error"),
		false,
	}, {
		Internal.New("internal").C("one", "two").F("one", "two"),
		Internal.New("internal").C("three", "four"),
		true,
	}, {
		Errors{Internal.New("internal_1"), Validation.New("validation_1").F("one", "two")},
		Errors{Internal.New("internal_2"), Validation.New("validation_1").F("one", "two")},
		false,
	}, {
		Errors{Internal.New("internal_1"), Validation.New("validation_1").F("one", "two")},
		Errors{Internal.New("internal_1"), Validation.New("validation_1").F("three", "four", "msg %d", 123)},
		true,
	}, {
		Errors{errors.New("hi"), errors.New("bye"), Status.New("status").S(404)},
		Errors{errors.New("hi"), Status.New("status").S(404)},
		false,
	}, {
		Errors{errors.New("hi"), errors.New("bye")},
		Errors{errors.New("hi"), errors.New("bye"), Status.New("status").S(404)},
		false,
	}, {
		Errors{errors.New("hi"), errors.New("bye"), Status.New("status").S(404)},
		Errors{errors.New("hi"), errors.New("bye"), Status.New("status").S(500)},
		true,
	}, {
		Errors{errors.New("hi"), errors.New("bye")},
		Errors{errors.New("bye"), errors.New("hi")},
		false,
	}, {
		Errors{errors.New("hi"), errors.New("bye"), Status.New("status").S(404)},
		Errors{errors.New("hi"), errors.New("bye"), Internal.New("status1").S(404)},
		false,
	}, {
		nil,
		nil,
		true,
	}, {
		nil,
		Internal.New("nil"),
		false,
	}, {
		Internal.New("nil"),
		nil,
		false,
	}}

	for i, test := range tests {
		assert.Equal(Compare(test.err1, test.err2), test.shouldEqual, i)
	}
}
