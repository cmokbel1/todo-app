package http

import (
	"errors"
	"net/http"

	"github.com/cmokbel1/todo-app/backend/todo"
)

type Error struct {
	Message string `json:"error"`
}

// StatusCode returns the HTTP status code for an error code.
func StatusCode(err error) int {
	var e *todo.Error
	if errors.As(err, &e) {
		switch e.Code {
		case todo.ECONFLICT:
			return http.StatusConflict
		case todo.EINVALID:
			return http.StatusBadRequest
		case todo.ENOTFOUND:
			return http.StatusNotFound
		case todo.EUNAUTHORIZED:
			return http.StatusUnauthorized
		}
	}

	return http.StatusInternalServerError
}

// ErrCodeFromStatus returns the corresponding error code for an HTTP status code.
func ErrCodeFromStatus(code int) string {
	switch code {
	case http.StatusConflict:
		return todo.ECONFLICT
	case http.StatusNotFound:
		return todo.ENOTFOUND
	case http.StatusBadRequest:
		return todo.EINVALID
	case http.StatusUnauthorized:
		return todo.EUNAUTHORIZED
	}

	return todo.EINTERNAL
}
