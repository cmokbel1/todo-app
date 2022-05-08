//go:build integration

package postgres_test

import (
	"context"
	"errors"
	"reflect"
	"strings"
	"testing"

	"github.com/cmokbel1/todo-app/backend/postgres"
	"github.com/cmokbel1/todo-app/backend/todo"
)

func newUser() *todo.User {
	v := *randstr(10)
	return &todo.User{Name: v}
}

func TestUserService_CreateUser(t *testing.T) {
	t.Parallel()

	createUser := func(t *testing.T, db *postgres.DB) (context.Context, *todo.User) {
		t.Helper()
		user := newUser()
		s := postgres.NewUserService(db)
		ctx := context.Background()
		if err := s.CreateUser(ctx, user); err != nil {
			t.Fatal(err)
		}
		return todo.NewContextWithUser(ctx, user), user
	}

	db := OpenDB(t)

	s := postgres.NewUserService(db)
	ctx := context.Background()

	t.Run("Success", func(t *testing.T) {
		ctx, user := createUser(t, db)
		if got, err := s.FindUserByID(ctx, user.ID); err != nil {
			t.Fatal(err)
		} else if want := user; !reflect.DeepEqual(got, want) {
			t.Errorf("want user %v got %v", want, got)
		}
	})

	t.Run("ErrConflictUserAlreadyExists", func(t *testing.T) {
		ctx, user := createUser(t, db)
		user.Name = strings.ToUpper(user.Name)

		if got, want := s.CreateUser(ctx, user), todo.Conflict; !errors.Is(got, want) {
			t.Fatalf("want error %v got %v", want, got)
		}
	})

	t.Run("ErrInvalidNameRequired", func(t *testing.T) {
		if got, want := s.CreateUser(ctx, &todo.User{}), todo.Invalid; !errors.Is(got, want) {
			t.Fatalf("want error %v got %v", want, got)
		}
	})
}

func TestUserService_DeleteUser(t *testing.T) {
	t.Parallel()

	db := OpenDB(t)
	s := postgres.NewUserService(db)

	t.Run("Success", func(t *testing.T) {
		user := newUser()

		ctx := context.Background()
		if err := s.CreateUser(ctx, user); err != nil {
			t.Fatal(err)
		}

		if err := s.DeleteUser(ctx, user.ID); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		userID := 5
		want := todo.NotFound
		got := s.DeleteUser(context.Background(), userID)
		if !errors.Is(got, want) {
			t.Fatalf("want error %v got %v", want, got)
		}
	})

	t.Run("ErrNotFoundInvalidUserID", func(t *testing.T) {
		if got, want := s.DeleteUser(context.Background(), -1), todo.NotFound; !errors.Is(got, want) {
			t.Fatalf("want error %v got %v", want, got)
		}
	})
}

func TestUserService_UpdateUser(t *testing.T) {
	t.Parallel()

	db := OpenDB(t)
	s := postgres.NewUserService(db)

	t.Run("Success", func(t *testing.T) {
		ctx := context.Background()
		user := newUser()

		if err := s.CreateUser(ctx, user); err != nil {
			t.Fatal(err)
		}
		user.Name = *randstr(10)
		user.Email = randstr(10)

		if user2, err := s.UpdateUser(ctx, user.ID, todo.UserUpdate{Name: &user.Name}); err != nil {
			t.Fatal(err)
		} else if got, want := user2.Name, user.Name; got != want {
			t.Fatalf("want user name %q got %q", want, got)
		} else if user2, err = s.UpdateUser(ctx, user.ID, todo.UserUpdate{Email: user.Email}); err != nil {
			t.Fatal(err)
		} else if got, want := *user2.Email, *user.Email; got != want {
			t.Fatalf("want user email %q got %q", want, got)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		ctx := context.Background()
		_, got := s.UpdateUser(ctx, 3, todo.UserUpdate{})
		if !errors.Is(got, todo.NotFound) {
			t.Fatalf("want error %v got %v", todo.NotFound, got)
		}
	})
}

func TestUserService_FindUserByID(t *testing.T) {
	t.Parallel()

	db := OpenDB(t)
	s := postgres.NewUserService(db)

	t.Run("Success", func(t *testing.T) {
		ctx := context.Background()
		user := newUser()
		if err := s.CreateUser(ctx, user); err != nil {
			t.Fatal(err)
		}
		if user2, err := s.FindUserByID(ctx, user.ID); err != nil {
			t.Fatal(err)
		} else if !reflect.DeepEqual(user2, user) {
			t.Fatalf("mismatch: %#v != %#v", user2, user)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		_, got := s.FindUserByID(context.Background(), 0)
		if !errors.Is(got, todo.NotFound) {
			t.Fatalf("want error %v got %v", todo.NotFound, got)
		}
	})
}

func TestUserService_FindUsers(t *testing.T) {
	db := OpenDB(t)
	s := postgres.NewUserService(db)

	t.Run("Success", func(t *testing.T) {
		ctx := context.Background()
		user := newUser()

		if err := s.CreateUser(ctx, user); err != nil {
			t.Fatal(err)
		}

		if users, err := s.FindUsers(ctx, todo.UserFilter{ID: &user.ID}); err != nil {
			t.Fatal(err)
		} else if got, want := len(users), 1; got != want {
			t.Fatalf("want %d users got %d", want, got)
		} else if !reflect.DeepEqual(users[0], user) {
			t.Fatalf("mismatch: %#v != %#v", users[0], user)
		}

		if users, err := s.FindUsers(ctx, todo.UserFilter{Name: &user.Name}); err != nil {
			t.Fatal(err)
		} else if got, want := len(users), 1; got != want {
			t.Fatalf("want %d users got %d", want, got)
		} else if !reflect.DeepEqual(users[0], user) {
			t.Fatalf("mismatch: %#v != %#v", users[0], user)
		}

		if users, err := s.FindUsers(ctx, todo.UserFilter{Email: user.Email}); err != nil {
			t.Fatal(err)
		} else if got, want := len(users), 1; got != want {
			t.Fatalf("want %d users got %d", want, got)
		} else if !reflect.DeepEqual(users[0], user) {
			t.Fatalf("mismatch: %#v != %#v", users[0], user)
		}
	})
}
