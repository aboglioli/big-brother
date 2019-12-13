package mock

import (
	"fmt"
	"testing"
)

const (
	Any    = "mock.Anything"
	Nil    = "mock.Nil"
	NotNil = "mock.NotNil"
)

type call struct {
	Func string
	Args []interface{}
	Ret  []interface{}
}

func Call(f string, args ...interface{}) call {
	return call{
		Func: f,
		Args: args,
	}
}

func (c call) Return(ret ...interface{}) call {
	c.Ret = ret
	return c
}

type Mock struct {
	Calls []call
}

func (m *Mock) Called(c call) {
	m.Calls = append(m.Calls, c)
}

func (m *Mock) Assert(t *testing.T, calls ...call) {
	if len(m.Calls) != len(calls) {
		callsStr := "expected:\n"
		for _, call := range calls {
			callsStr += fmt.Sprintf("- %s %v -> %v\n", call.Func, call.Args, call.Ret)
		}

		callsStr += "actual:\n"
		for _, call := range m.Calls {
			callsStr += fmt.Sprintf("- %s %v -> %v\n", call.Func, call.Args, call.Ret)
		}
		t.Fatalf("MOCK: Different number of calls\n%s%s\n", callsStr, stackInfo())
	}

	for i, call1 := range m.Calls {
		call2 := calls[i]

		if call1.Func != call2.Func || !compareArgs(call1.Args, call2.Args) || (len(call2.Ret) > 0 && !compareArgs(call1.Ret, call2.Ret)) {
			t.Fatalf("MOCK:\nexpected: %s %v -> %v\nactual: %s %v -> %v\n%s\n", call2.Func, call2.Args, call2.Ret, call1.Func, call1.Args, call1.Ret, stackInfo())
		}
	}
}

func (m *Mock) Reset() {
	m.Calls = []call{}
}

func compareArgs(args1 []interface{}, args2 []interface{}) bool {
	if len(args1) != len(args2) {
		return false
	}

	for i, arg1 := range args1 {
		arg2 := args2[i]
		if arg1 == Any || arg2 == Any {
			continue
		}

		if arg2 == Nil && arg1 == nil {
			continue
		}

		if arg2 == NotNil && arg1 != nil {
			continue
		}

		if arg1 != arg2 {
			return false
		}
	}

	return true
}
