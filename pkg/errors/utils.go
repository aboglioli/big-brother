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

		sameFields := true
		if len(err1.Fields) > 0 {
			if len(err1.Fields) != len(err2.Fields) {
				return false
			}
			for i, field1 := range err1.Fields {
				field2 := err2.Fields[i]
				if field1.Field != field2.Field || field1.Code != field2.Code {
					sameFields = false
					break
				}
			}
		}

		sameCause := true
		if err1.Cause != nil {
			if err2.Cause == nil {
				return false
			}
			sameCause = Compare(err1.Cause, err2.Cause)
		}

		return err1.Equals(err2) && sameFields && sameCause
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
