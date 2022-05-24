package todo_test

import (
	"errors"
	"testing"

	"github.com/cmokbel1/todo-app/backend/todo"
)

func TestError(t *testing.T) {
	tt := []struct {
		Error       error
		InTarget    error
		Want        bool
		WantMessage string
		WantCode    string
	}{
		{
			Error:       todo.Err(todo.EINVALID, ""),
			InTarget:    todo.Invalid,
			Want:        true,
			WantMessage: "",
			WantCode:    todo.EINVALID,
		},
		{
			Error:       todo.Err(todo.EINTERNAL, "test message"),
			InTarget:    todo.Internal,
			Want:        true,
			WantMessage: "test message",
			WantCode:    todo.EINTERNAL,
		},
		{
			Error:       todo.Err(todo.ENOTFOUND, ""),
			InTarget:    todo.NotFound,
			Want:        true,
			WantMessage: "",
			WantCode:    todo.ENOTFOUND,
		},
		{
			Error:       todo.Err(todo.EUNAUTHORIZED, ""),
			InTarget:    todo.Unauthorized,
			Want:        true,
			WantMessage: "",
			WantCode:    todo.EUNAUTHORIZED,
		},
		{
			Error:       todo.Err(todo.ECONFLICT, ""),
			InTarget:    todo.Conflict,
			Want:        true,
			WantMessage: "",
			WantCode:    todo.ECONFLICT,
		},
		{
			Error:       todo.Err(todo.EINVALID, ""),
			InTarget:    todo.Err(todo.EINVALID+"_different_reason", "different message"),
			Want:        false,
			WantMessage: "",
			WantCode:    todo.EINVALID,
		},
		{
			Error:       todo.Err(todo.EINVALID, ""),
			InTarget:    nil,
			Want:        false,
			WantMessage: "",
			WantCode:    todo.EINVALID,
		},
	}

	for i, tc := range tt {
		if got, want := errors.Is(tc.Error, tc.InTarget), tc.Want; want != got {
			t.Errorf("%d want %v got %v", i, want, got)
		}

		if got, want := todo.ErrCode(tc.Error), tc.WantCode; want != got {
			t.Errorf("%d want code %q got %q", i, want, got)
		} else if got, want = todo.ErrMessage(tc.Error), tc.WantMessage; want != got {
			t.Errorf("%d want message %q got %q", i, want, got)
		}
	}
}
