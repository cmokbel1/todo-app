//go:build integration

package postgres_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/cmokbel1/todo-app/backend/postgres"
	"github.com/cmokbel1/todo-app/backend/todo"
)

func TestItemListService(t *testing.T) {
	t.Parallel()

	createUser := func(t *testing.T, db *postgres.DB) (context.Context, *todo.User) {
		t.Helper()
		user := &todo.User{Name: *randstr(10), Password: *randstr(10)}
		s := postgres.NewUserService(db)
		ctx := context.Background()
		if err := s.CreateUser(ctx, user); err != nil {
			t.Fatal(err)
		}
		return todo.NewContextWithUser(ctx, user), user
	}

	createUserAndList := func(t *testing.T, db *postgres.DB) (context.Context, *todo.User, *todo.List) {
		t.Helper()
		ctx, user := createUser(t, db)
		s := postgres.NewItemListService(db)
		list := &todo.List{UserID: user.ID, Name: *randstr(10)}
		if err := s.CreateList(ctx, list); err != nil {
			t.Fatal(err)
		}

		return ctx, user, list
	}

	createUserAndListWithItems := func(t *testing.T, db *postgres.DB) (context.Context, *todo.User, *todo.List) {
		t.Helper()
		ctx, user := createUser(t, db)
		s := postgres.NewItemListService(db)
		list := &todo.List{UserID: user.ID, Name: *randstr(10)}
		if err := s.CreateList(ctx, list); err != nil {
			t.Fatal(err)
		}

		item1 := &todo.Item{ListID: list.ID, UserID: user.ID, Name: *randstr(10)}
		item2 := &todo.Item{ListID: list.ID, UserID: user.ID, Name: *randstr(10)}
		var err error
		if err = s.CreateItem(ctx, item1); err != nil {
			t.Fatal(err)
		}
		if err = s.CreateItem(ctx, item2); err != nil {
			t.Fatal(err)
		}
		list.Items = append(list.Items, item1, item2)
		return ctx, user, list
	}

	t.Run("CreateList", func(t *testing.T) {
		db := OpenDB(t)

		t.Run("Success", func(t *testing.T) {
			ctx, user := createUser(t, db)
			s := postgres.NewItemListService(db)
			list := &todo.List{UserID: user.ID, Name: "Name"}

			if err := s.CreateList(ctx, list); err != nil {
				t.Fatal(err)
			}

			if got, err := s.FindListByID(ctx, list.ID); err != nil {
				t.Fatal(err)
			} else if !reflect.DeepEqual(got, list) {
				t.Fatalf("want list %v got %v", list, got)
			}
		})
		t.Run("Unauthorized", func(t *testing.T) {
			s := postgres.NewItemListService(db)
			list := &todo.List{Name: "Name"}
			if got, want := s.CreateList(context.Background(), list), todo.Unauthorized; !errors.Is(got, want) {
				t.Fatalf("want error %v got %v", want, got)
			}
		})
	})

	t.Run("ReadListWithItems", func(t *testing.T) {
		db := OpenDB(t)

		ctx, _, list := createUserAndListWithItems(t, db)
		s := postgres.NewItemListService(db)

		if got, err := s.FindListByID(ctx, list.ID); err != nil {
			t.Fatal(err)
		} else if want := list; !reflect.DeepEqual(got, want) {
			t.Fatalf("want list %v got %v", want, got)
		}
	})

	t.Run("DeleteList", func(t *testing.T) {
		db := OpenDB(t)

		t.Run("Success", func(t *testing.T) {
			ctx, _, list := createUserAndList(t, db)
			s := postgres.NewItemListService(db)

			if err := s.DeleteList(ctx, list.ID); err != nil {
				t.Fatal(err)
			}

			if _, got := s.FindListByID(ctx, list.ID); got == nil {
				t.Fatal("want error got none")
			} else if !errors.Is(got, todo.NotFound) {
				t.Fatalf("want error %v got %v", todo.NotFound, got)
			}
		})

		t.Run("ErrNotFoundOtherUsersList", func(t *testing.T) {
			_, _, list := createUserAndList(t, db)
			ctx, _, _ := createUserAndList(t, db)

			s := postgres.NewItemListService(db)
			if err := s.DeleteList(ctx, list.ID); err == nil {
				t.Fatal("want err got none")
			} else if got, want := err, todo.NotFound; !errors.Is(got, want) {
				t.Fatalf("want error %v got %v", got, want)
			}
		})
	})

	t.Run("UpdateList", func(t *testing.T) {
		db := OpenDB(t)

		t.Run("Success", func(t *testing.T) {
			ctx, _, list := createUserAndList(t, db)
			s := postgres.NewItemListService(db)
			var upd todo.ListUpdate
			{
				upd.Name = randstr(10)
			}

			if got, err := s.UpdateList(ctx, list.ID, upd); err != nil {
				t.Fatal(err)
			} else if want, err := s.FindListByID(ctx, got.ID); err != nil {
				t.Fatal(err)
			} else if !reflect.DeepEqual(got, want) {
				t.Fatalf("want list %v got %v", want, got)
			}

			{
				completed := true
				upd.Completed = &completed
			}

			if got, err := s.UpdateList(ctx, list.ID, upd); err != nil {
				t.Fatal(err)
			} else if want, err := s.FindListByID(ctx, got.ID); err != nil {
				t.Fatal(err)
			} else if !reflect.DeepEqual(got, want) {
				t.Fatalf("want list %v got %v", want, got)
			}
		})

		t.Run("ErrUnauthorizedNoUser", func(t *testing.T) {
			_, _, list := createUserAndList(t, db)
			s := postgres.NewItemListService(db)

			_, got := s.UpdateList(context.Background(), list.ID, todo.ListUpdate{})
			if !errors.Is(got, todo.Unauthorized) {
				t.Errorf("want error %v got %v", todo.Unauthorized, got)
			}
		})

		t.Run("ErrUnauthorizedWrongUser", func(t *testing.T) {
			ctx, _, _ := createUserAndList(t, db)
			_, _, list2 := createUserAndList(t, db)

			s := postgres.NewItemListService(db)
			if _, got := s.UpdateList(ctx, list2.ID, todo.ListUpdate{}); !errors.Is(got, todo.Unauthorized) {
				t.Errorf("want error %v got %v", todo.Unauthorized, got)
			}
		})
	})

	t.Run("CreateItem", func(t *testing.T) {
		db := OpenDB(t)
		t.Run("Success", func(t *testing.T) {
			ctx, user, list := createUserAndList(t, db)
			s := postgres.NewItemListService(db)
			item := &todo.Item{ListID: list.ID, UserID: user.ID, Name: *randstr(10)}

			if err := s.CreateItem(ctx, item); err != nil {
				t.Fatal(err)
			}
		})

		t.Run("NotFound", func(t *testing.T) {
			ctx, user, _ := createUserAndList(t, db)
			s := postgres.NewItemListService(db)
			item := &todo.Item{ListID: 0, UserID: user.ID, Name: *randstr(10)}
			if err := s.CreateItem(ctx, item); err == nil {
				t.Fatal("want error but got none")
			} else if got, want := err, todo.NotFound; !errors.Is(got, want) {
				t.Fatalf("want error %v got %v", want, got)
			}
		})

		t.Run("ErrUnauthorizedNoUser", func(t *testing.T) {
			_, user, list := createUserAndList(t, db)
			s := postgres.NewItemListService(db)
			item := &todo.Item{ListID: list.ID, UserID: user.ID, Name: *randstr(10)}

			if got := s.CreateItem(context.Background(), item); got == nil {
				t.Fatal("want error but got none")
			} else if !errors.Is(got, todo.Unauthorized) {
				t.Fatalf("want error %v got %v", todo.Unauthorized, got)
			}
		})

		t.Run("ErrUnauthorizedDifferentListOwner", func(t *testing.T) {
			_, _, list := createUserAndList(t, db)
			ctx2, user, _ := createUserAndList(t, db)
			s := postgres.NewItemListService(db)
			item := &todo.Item{ListID: list.ID, UserID: user.ID, Name: *randstr(10)}

			if got := s.CreateItem(ctx2, item); got == nil {
				t.Fatal("want error but got none")
			} else if !errors.Is(got, todo.Unauthorized) {
				t.Fatalf("want error %v got %v", todo.Unauthorized, got)
			}
		})
	})

	t.Run("DeleteItem", func(t *testing.T) {
		db := OpenDB(t)
		t.Run("Success", func(t *testing.T) {
			ctx, user, list := createUserAndList(t, db)
			s := postgres.NewItemListService(db)
			item := &todo.Item{ListID: list.ID, UserID: user.ID, Name: *randstr(10)}

			if err := s.CreateItem(ctx, item); err != nil {
				t.Fatal(err)
			}

			if err := s.DeleteItem(ctx, item.ID); err != nil {
				t.Fatal(err)
			}
		})

		t.Run("NotFound", func(t *testing.T) {
			ctx, user, _ := createUserAndList(t, db)
			s := postgres.NewItemListService(db)
			item := &todo.Item{ID: 999, UserID: user.ID, Name: *randstr(10)}
			if err := s.DeleteItem(ctx, item.ID); err == nil {
				t.Fatal("want error but got none")
			} else if got, want := err, todo.NotFound; !errors.Is(got, want) {
				t.Fatalf("want error %v got %v", want, got)
			}
		})

		t.Run("Unauthorized", func(t *testing.T) {
			s := postgres.NewItemListService(db)
			if err := s.DeleteItem(context.Background(), 5); err == nil {
				t.Fatal("want error but got none")
			} else if got, want := err, todo.Unauthorized; !errors.Is(got, want) {
				t.Fatalf("want error %v got %v", want, got)
			}
		})

		t.Run("ErrInvalidItemID", func(t *testing.T) {
			ctx, user, _ := createUserAndList(t, db)
			s := postgres.NewItemListService(db)
			item := &todo.Item{ID: 0, UserID: user.ID, Name: *randstr(10)}
			if err := s.DeleteItem(ctx, item.ID); err == nil {
				t.Fatal("want error but got none")
			} else if got, want := err, todo.Invalid; !errors.Is(got, want) {
				t.Fatalf("want error %v got %v", want, got)
			}
		})
	})

	t.Run("UpdateItem", func(t *testing.T) {
		db := OpenDB(t)

		t.Run("Success", func(t *testing.T) {
			ctx, _, list := createUserAndListWithItems(t, db)
			s := postgres.NewItemListService(db)

			if got, err := s.UpdateItem(ctx, list.Items[0].ID, todo.ItemUpdate{Name: randstr(15)}); err != nil {
				t.Fatal(err)
			} else if want, err := s.FindItemByID(ctx, list.Items[0].ID); err != nil {
				t.Fatal(err)
			} else if !reflect.DeepEqual(want, got) {
				t.Fatalf("want item %v got %v", want, got)
			}
		})

		t.Run("NotFound", func(t *testing.T) {
			ctx, _, _ := createUserAndList(t, db)
			s := postgres.NewItemListService(db)
			completed := true
			upd := todo.ItemUpdate{Name: randstr(10), Completed: &completed}

			_, got := s.UpdateItem(ctx, -10000, upd)
			if want := todo.NotFound; !errors.Is(got, want) {
				t.Fatalf("want error %v got %v", want, got)
			}
		})

		t.Run("Unauthorized", func(t *testing.T) {
			_, _, list := createUserAndListWithItems(t, db)
			upd := todo.ItemUpdate{Name: randstr(10)}
			s := postgres.NewItemListService(db)
			_, got := s.UpdateItem(context.Background(), list.Items[0].ID, upd)
			if want := todo.Unauthorized; !errors.Is(got, want) {
				t.Fatalf("want error %v got %v", want, got)
			}
		})
	})
}
