package errors

import (
	"net/http"
)

var (
	ErrRequest        = Status.New("invalid_request").S(http.StatusBadRequest)
	ErrInternalServer = Status.New("internal_server").S(http.StatusInternalServerError)
	ErrNotFound       = Status.New("not_found").S(http.StatusNotFound)
	ErrUnauthorized   = Status.New("unauthorized").S(http.StatusUnauthorized)
)
