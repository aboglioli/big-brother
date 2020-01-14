package errors

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateError(t *testing.T) {
	assert := assert.New(t)

	tests := []struct {
		in  error
		out Error
	}{{
		Internal.New("INTERNAL"),
		Error{
			Type: Internal,
			Code: "INTERNAL",
		},
	}, {
		Status.New("STATUS"),
		Error{
			Type: Status,
			Code: "STATUS",
		},
	}, {
		Validation.New("VALIDATION"),
		Error{
			Type: Validation,
			Code: "VALIDATION",
		},
	}, {
		Validation.New("VALIDATION").P("path").M("new %s", "msg"),
		Error{
			Type:    Validation,
			Code:    "VALIDATION",
			Path:    "path",
			Message: "new msg",
		},
	}, {
		Validation.New("V").F("field1", "invalid").F("field2", "invalid", "id %d", 123),
		Error{
			Type:   Validation,
			Code:   "V",
			Fields: []Field{{"field1", "invalid", ""}, {"field2", "invalid", "id 123"}},
		},
	}, {
		Status.New("S").S(404).C("id", "123").C("session", "S-%d", 456),
		Error{
			Type:   Status,
			Code:   "S",
			Status: 404,
			Context: Context{
				"id":      "123",
				"session": "S-456",
			},
		},
	}}

	for i, test := range tests {
		in, ok := test.in.(Error)
		assert.True(ok, i)
		assert.Equal(test.out, in, i)
	}
}

func TestWrapError(t *testing.T) {
	assert := assert.New(t)

	err := errors.New("err")

	tests := []struct {
		in  error
		out error
	}{{
		Internal.New("I").Wrap(err),
		err,
	}, {
		Internal.New("I1").Wrap(Internal.New("I2").Wrap(err)),
		Internal.New("I2").Wrap(err),
	}, {
		Status.New("S").Wrap(Validation.New("V").Wrap(err)),
		Validation.New("V").Wrap(err),
	}}

	for i, test := range tests {
		in, ok := test.in.(Error)
		assert.True(ok, i)
		assert.Equal(test.out, in.Cause, i)
	}
}

func TestReuseError(t *testing.T) {
	assert := assert.New(t)

	rawErr1 := errors.New("err1")
	rawErr2 := errors.New("err2")

	err1 := Internal.New("I").P("err1").M("err%d", 1).S(1)
	err2 := err1.P("err2").M("err%d", 2).S(2)
	err3 := err2.M("err%d", 3).Wrap(rawErr1)
	err4 := err3.C("id", "user%d", 123).F("field1", "invalid", "s-%d", 456).Wrap(rawErr2)
	err5 := err4.C("one", "two").F("one", "two").S(5)
	err6 := err4.C("k", "v").F("f", "c")

	tests := []struct {
		in  Error
		out Error
	}{{
		err1,
		Error{
			Type:    Internal,
			Code:    "I",
			Path:    "err1",
			Message: "err1",
			Status:  1,
		},
	}, {
		err2,
		Error{
			Type:    Internal,
			Code:    "I",
			Path:    "err2",
			Message: "err2",
			Status:  2,
		},
	}, {
		err3,
		Error{
			Type:    Internal,
			Code:    "I",
			Path:    "err2",
			Message: "err3",
			Status:  2,
			Cause:   rawErr1,
		},
	}, {
		err4,
		Error{
			Type:    Internal,
			Code:    "I",
			Path:    "err2",
			Message: "err3",
			Status:  2,
			Context: Context{
				"id": "user123",
			},
			Fields: []Field{{
				Field:   "field1",
				Code:    "invalid",
				Message: "s-456",
			}},
			Cause: rawErr2,
		},
	}, {
		err5,
		Error{
			Type:    Internal,
			Code:    "I",
			Path:    "err2",
			Message: "err3",
			Status:  5,
			Context: Context{
				"id":  "user123",
				"one": "two",
			},
			Fields: []Field{{
				Field:   "field1",
				Code:    "invalid",
				Message: "s-456",
			}, {
				Field: "one",
				Code:  "two",
			}},
			Cause: rawErr2,
		},
	}, {
		err6,
		Error{
			Type:    Internal,
			Code:    "I",
			Path:    "err2",
			Message: "err3",
			Status:  2,
			Context: Context{
				"id": "user123",
				"k":  "v",
			},
			Fields: []Field{{
				Field:   "field1",
				Code:    "invalid",
				Message: "s-456",
			}, {
				Field: "f",
				Code:  "c",
			}},
			Cause: rawErr2,
		},
	}}

	for i, test := range tests {
		assert.Equal(test.out, test.in, i)
	}
}
