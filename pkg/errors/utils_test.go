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
		false,
	}, {
		Errors{Internal.New("internal_1"), Validation.New("validation_1").F("one", "two")},
		Errors{Internal.New("internal_2"), Validation.New("validation_1").F("one", "two")},
		false,
	}, {
		Errors{Internal.New("internal_1"), Validation.New("validation_1").F("one", "two")},
		Errors{Internal.New("internal_1"), Validation.New("validation_1").F("three", "four", "msg %d", 123)},
		false,
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
	}, {
		Internal.New("internal_1"),
		Internal.New("internal_1").F("field", "required"),
		true,
	}, {
		Internal.New("internal_1").F("field", "required"),
		Internal.New("internal_1"),
		false,
	}, {
		Internal.New("internal_1").F("field", "required"),
		Internal.New("internal_1").F("field", "required"),
		true,
	}, {
		Internal.New("internal_1").F("field", "required", "msg %s", "random"),
		Internal.New("internal_1").F("field", "required", "hello %s", "world"),
		true,
	}, {
		Internal.New("internal_1").F("field", "required", "msg %s", "random"),
		Internal.New("internal_1").F("field", "not_available", "hello %s", "world"),
		false,
	}, {
		Status.New("status_1").Wrap(Internal.New("internal_1")),
		Status.New("status_1").Wrap(Internal.New("internal_1")),
		true,
	}, {
		Status.New("status_1").Wrap(Internal.New("internal_1")),
		Status.New("status_1").Wrap(Internal.New("internal_2")),
		false,
	}, {
		Status.New("status_1").Wrap(Internal.New("internal_1")),
		Status.New("status_1"),
		false,
	}, {
		Status.New("status_1"),
		Status.New("status_1").Wrap(Internal.New("internal_1")),
		true,
	}, {
		Status.New("status_1").Wrap(Internal.New("internal_1").Wrap(errors.New("err"))),
		Status.New("status_1").Wrap(Internal.New("internal_1").Wrap(errors.New("err"))),
		true,
	}, {
		Status.New("status_1").Wrap(Internal.New("internal_1").Wrap(errors.New("err1"))),
		Status.New("status_1").Wrap(Internal.New("internal_1").Wrap(errors.New("err2"))),
		false,
	}, {
		Status.New("status_1").Wrap(Internal.New("internal_1").Wrap(errors.New("err"))),
		Status.New("status_1").Wrap(Internal.New("internal_1")),
		false,
	}, {
		Status.New("status_1").Wrap(Internal.New("internal_1")),
		Status.New("status_1").Wrap(Internal.New("internal_1").Wrap(errors.New("err"))),
		true,
	}, {
		Status.New("status_1"),
		Status.New("status_1").Wrap(Internal.New("internal_1").Wrap(errors.New("err"))),
		true,
	}, {
		Status.New("status_1").Wrap(Internal.New("internal_1").F("field", "required", "msg1")),
		Status.New("status_1").Wrap(Internal.New("internal_1").F("field", "required", "msg2").Wrap(errors.New("err"))),
		false,
	}, {
		Status.New("status_1").Wrap(Internal.New("internal_1").F("field", "required", "msg1")),
		Status.New("status_1").Wrap(Internal.New("internal_1").F("field", "not_available", "msg2")),
		false,
	}}

	for i, test := range tests {
		if !assert.Equal(Compare(test.err1, test.err2), test.shouldEqual, i) {
			t.Errorf("Error %d:\nexpected: %s\nactual  : %s", i, test.err1, test.err2)
		}
	}
}
