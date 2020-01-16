package errors

import "testing"

func Compare(expected, actual error) bool {
	if expected == nil && actual == nil {
		return true
	}

	switch err1 := expected.(type) {
	case Errors:
		err2, ok := actual.(Errors)
		if !ok {
			return false
		}

		if len(err1) != len(err2) {
			return false
		}

		for i, err1 := range err1 {
			err2 := err2[i]
			if !Compare(err1, err2) {
				return false
			}
		}

		return true
	case Error:
		err2, ok := actual.(Error)
		if !ok {
			return false
		}

		if len(err1.Fields) > 0 {
			if len(err1.Fields) != len(err2.Fields) {
				return false
			}
			for i, field1 := range err1.Fields {
				field2 := err2.Fields[i]
				if field1.Field != field2.Field || field1.Code != field2.Code {
					return false
				}
			}
		}

		if err1.Cause != nil {
			if err2.Cause == nil {
				return false
			}
			if !Compare(err1.Cause, err2.Cause) {
				return false
			}
		}

		if len(err1.Context) > 0 {
			if len(err1.Context) != len(err2.Context) {
				return false
			}
			for k, v1 := range err1.Context {
				v2, ok := err2.Context[k]
				if !ok {
					return false
				}
				if v1 != v2 {
					return false
				}
			}
		}

		return err1.Equals(err2)
	case error:
		return err1.Error() == actual.Error()
	}

	return false
}

func Assert(t *testing.T, expectedErr, actualErr error) {
	if !Compare(expectedErr, actualErr) {
		t.Errorf("Error:\nexpected: %s\nactual  : %s", expectedErr, actualErr)
	}
}
