package mock

import (
	"errors"
	"testing"
)

type mockStruct struct {
	Mock
}

func TestMock(t *testing.T) {
	err := errors.New("err")

	tests := []struct {
		in  []call
		out []call
	}{{
		in:  []call{Call("simple"), Call("simple")},
		out: []call{Call("simple"), Call("simple")},
	}, {
		in: []call{
			Call("FirstMethod", "one", "two"),
			Call("SecondMethod", "one", 2),
			Call("ThirdMethod", 1, "two"),
		},
		out: []call{
			call{"FirstMethod", []interface{}{"one", "two"}, []interface{}{}},
			call{"SecondMethod", []interface{}{"one", 2}, []interface{}{}},
			call{"ThirdMethod", []interface{}{1, "two"}, []interface{}{}},
		},
	}, {
		in:  []call{Call("WithAny", 1, Any, 3)},
		out: []call{Call("WithAny", 1, 45, 3)},
	}, {
		in: []call{
			Call("One", 1, 2),
			Call("Two", 1, 2),
		},
		out: []call{
			Call("One", 1, Any),
			call{"Two", []interface{}{Any, 2}, []interface{}{}},
		},
	}, {
		in: []call{
			Call("Method1", "data", nil),
			Call("Method1", nil, err),
			Call("Method2", 1, 2, nil),
		},
		out: []call{
			Call("Method1", "data", Nil),
			Call("Method1", nil, NotNil),
			Call("Method2", 1, NotNil, Nil),
		},
	}, {
		in: []call{
			Call("Method1", "data").Return("return", nil),
			Call("Method1", "data").Return(nil, err),
		},
		out: []call{
			Call("Method1", "data"),
			Call("Method1", "data"),
		},
	}, {
		in: []call{
			Call("Method1", "data").Return("return", nil),
			Call("Method1", "data").Return(nil, err),
		},
		out: []call{
			Call("Method1", "data").Return("return", nil),
			Call("Method1", "data").Return(nil, err),
		},
	}, {
		in: []call{
			Call("Method1", "data").Return("return", nil),
			Call("Method1", "data").Return(nil, err),
		},
		out: []call{
			Call("Method1", "data").Return("return", Nil),
			Call("Method1", NotNil).Return(Nil, NotNil),
		},
	}}

	for _, test := range tests {
		m := &mockStruct{}
		for _, call := range test.in {
			m.Called(call)
		}

		m.Assert(t, test.out...)
	}
}
