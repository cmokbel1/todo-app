package todo

import (
	"errors"
	"fmt"
	"strings"
)

const (
	ENOTFOUND     = "NOT_FOUND"
	EINVALID      = "INVALID"
	EUNAUTHORIZED = "UNAUTHORIZED"
	ECONFLICT     = "CONFLICT"
	EINTERNAL     = "INTERNAL"
)

// the following errors are intended to be used as sentinel values to determine error likeness
// using errors.Is
var (
	Unauthorized   = &Error{EUNAUTHORIZED, "unauthorized"}
	Internal       = &Error{EINTERNAL, "internal error"}
	NotFound       = &Error{ENOTFOUND, "not found"}
	Invalid        = &Error{EINVALID, "invalid"}
	Conflict       = &Error{ECONFLICT, "conflict"}
	NotImplemented = &Error{EINTERNAL, "not implemented"}
)

type Error struct {
	// Code is the application level error code
	Code string
	// Message is either a loggable and or human-readable message.
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s (%s)", e.Message, e.Code)
}

func (e Error) Is(target error) bool {
	// when checking against sentinel values we only use the error code to determine likeness.
	switch target {
	case NotFound, Unauthorized, Internal, Invalid, Conflict:
		// check for an exact code match or a prefix match,
		// e.g. for Invalid (who's code is `invalid`) it will match a code with 'invalid_name_required'
		code := ErrCode(target)
		return e.Code == code || strings.HasPrefix(e.Code, code)
	}
	return e.Code == ErrCode(target) && e.Message == ErrMessage(target)
}

// Err is a utility method to create an Error with a Code.
func Err(code string, tmpl string, args ...interface{}) error {
	if code == "" {
		panic("code cannot be empty")
	}

	return &Error{
		Code:    code,
		Message: fmt.Sprintf(tmpl, args...),
	}
}

// ErrCode extracts the ErrCode from an Error. If err is not an Error then EINTERNAL is returned.
func ErrCode(err error) string {
	var e *Error
	if err == nil {
		return ""
	} else if errors.As(err, &e) {
		return e.Code
	}
	return EINTERNAL
}

// ErrMessage extracts the Message from an error. If err is not an Error then "internal error" is returned.
func ErrMessage(err error) string {
	var e *Error
	if err == nil {
		return ""
	} else if errors.As(err, &e) {
		return e.Message
	}
	return "internal error"
}
