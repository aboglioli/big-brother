package errors

import "fmt"

func Compare(err1, err2 error) bool {
	switch err1 := err1.(type) {
	case Errors:
		err2, ok := err2.(Errors)
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
		err2, ok := err2.(Error)
		if !ok {
			return false
		}

		return err1.Equals(err2)
	case error:
		return err1.Error() == err2.Error()
	}

	return false
}

func Print(err error) string {
	switch err := err.(type) {
	case Errors:
		str := "Errors:"
		for i, err := range err {
			str = fmt.Sprintf("%s\n\t%d: %s", str, i, Print(err))
		}
		return str
	case Error:

	}
	return ""
}
