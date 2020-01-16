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
