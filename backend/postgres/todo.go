package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/cmokbel1/todo-app/backend/todo"
)

var _ todo.ItemListService = (*ItemListService)(nil)

func NewItemListService(db *DB) *ItemListService {
	return &ItemListService{db: db}
}

type ItemListService struct {
	db *DB
}

func (svc *ItemListService) FindListByID(ctx context.Context, id int) (*todo.List, error) {
	tx, err := svc.db.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	list, err := findTodoListByID(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	return list, tx.Commit()
}

func findTodoListByID(ctx context.Context, tx *Tx, id int) (*todo.List, error) {
	lists, err := findTodoLists(ctx, tx, todo.ListFilter{ID: &id})
	if err != nil {
		return nil, err
	} else if len(lists) == 0 {
		return nil, todo.Err(todo.ENOTFOUND, "could not find list with id %d", id)
	}
	return lists[0], nil
}

func (svc *ItemListService) FindLists(ctx context.Context, f todo.ListFilter) ([]*todo.List, error) {
	tx, err := svc.db.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	lists, err := findTodoLists(ctx, tx, f)
	if err != nil {
		return nil, err
	}
	return lists, tx.Commit()
}

func findTodoLists(ctx context.Context, tx *Tx, f todo.ListFilter) ([]*todo.List, error) {
	user, err := todo.ValidUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var args []interface{}
	where := []string{"1 = 1"}
	if v := f.ID; v != nil {
		where, args = append(where, fmt.Sprintf("id = $%d", len(where))), append(args, *v)
	}

	if v := f.UserID; v != nil {
		where, args = append(where, fmt.Sprintf("user_id = $%d", len(where))), append(args, *v)
	}

	if v := f.Name; v != nil {
		where, args = append(where, fmt.Sprintf("name = $%d", len(where))), append(args, *v)
	}

	if v := f.Completed; v != nil {
		where, args = append(where, fmt.Sprintf("completed = $%d", len(where))), append(args, *v)
	}

	query := `
	SELECT 
		id, 
		user_id,
		name, 
		completed, 
		created_at, 
		updated_at 
	FROM lists
	WHERE ` + strings.Join(where, " AND ") + `
	ORDER BY id ASC ` + FormatLimitOffset(f.Limit, f.Offset)
	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	lists := make([]*todo.List, 0)
	for rows.Next() {
		list := todo.List{Items: make([]*todo.Item, 0)}
		if err := rows.Scan(
			&list.ID,
			&list.UserID,
			&list.Name,
			&list.Completed,
			(*Time)(&list.CreatedAt),
			(*Time)(&list.UpdatedAt),
		); err != nil {
			return nil, err
		}

		if list.UserID != user.ID {
			return nil, todo.Unauthorized
		}
		lists = append(lists, &list)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	for _, list := range lists {
		items, err := findTodoItems(ctx, tx, todo.ItemFilter{ListID: &list.ID})
		if err != nil {
			return nil, err
		}
		list.Items = append(list.Items, items...)
	}

	return lists, nil
}

func (svc *ItemListService) UpdateList(ctx context.Context, id int, upd todo.ListUpdate) (*todo.List, error) {
	tx, err := svc.db.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	list, err := updateTodoList(ctx, tx, id, upd)
	if err != nil {
		return nil, err
	}
	return list, tx.Commit()
}

func updateTodoList(ctx context.Context, tx *Tx, id int, upd todo.ListUpdate) (*todo.List, error) {
	list, err := findTodoListByID(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	list.UpdatedAt = tx.now
	if v := upd.Name; v != nil {
		list.Name = *v
	}
	if v := upd.Completed; v != nil {
		list.Completed = *v
	}

	if _, err := tx.ExecContext(ctx, `
	UPDATE lists 
	SET name = $1,
		completed = $2,
		updated_at = $3
	WHERE id = $4 AND user_id = $5`,
		list.Name, list.Completed, (*Time)(&list.UpdatedAt), list.ID, list.UserID); err != nil {
		return list, err
	}

	return list, nil
}

func (svc *ItemListService) CreateList(ctx context.Context, list *todo.List) error {
	tx, err := svc.db.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := createTodoList(ctx, tx, list); err != nil {
		return err
	}

	return tx.Commit()
}

func createTodoList(ctx context.Context, tx *Tx, list *todo.List) error {
	user, err := todo.ValidUserFromContext(ctx)
	if err != nil {
		return err
	}

	list.CreatedAt = tx.now
	list.UpdatedAt = list.CreatedAt
	list.UserID = user.ID
	list.Items = make([]*todo.Item, 0)

	if err := list.Validate(); err != nil {
		return err
	}

	var id int64
	err = tx.QueryRowContext(ctx, `
INSERT INTO lists (user_id, name, completed, created_at, updated_at) 
VALUES ($1, $2, $3, $4, $5)
RETURNING id`,
		list.UserID,
		list.Name,
		list.Completed,
		(*Time)(&list.CreatedAt),
		(*Time)(&list.UpdatedAt)).Scan(&id)
	if err != nil {
		return err
	}
	list.ID = int(id)

	return nil
}

func (svc *ItemListService) DeleteList(ctx context.Context, id int) error {
	tx, err := svc.db.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err = deleteTodoList(ctx, tx, id); err != nil {
		return err
	}
	return tx.Commit()
}

func deleteTodoList(ctx context.Context, tx *Tx, id int) error {
	user := todo.UserFromContext(ctx)
	if user == nil {
		return todo.Unauthorized
	}

	if id <= 0 {
		return todo.Err(todo.EINVALID, "invalid id")
	}

	result, err := tx.ExecContext(ctx, `DELETE FROM lists WHERE id = $1 AND user_id = $2`, id, user.ID)
	if err != nil {
		return err
	} else if n, _ := result.RowsAffected(); n == 0 {
		return todo.Err(todo.ENOTFOUND, "could not delete list with id %v", id)
	}

	return nil
}

func (svc *ItemListService) FindItemByID(ctx context.Context, id int) (*todo.Item, error) {
	tx, err := svc.db.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	item, err := findTodoItem(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	return item, tx.Commit()
}

func findTodoItem(ctx context.Context, tx *Tx, id int) (*todo.Item, error) {
	items, err := findTodoItems(ctx, tx, todo.ItemFilter{ID: &id})
	if err != nil {
		return nil, err
	} else if len(items) == 0 {
		return nil, todo.Err(todo.ENOTFOUND, "could not find item with id %q", id)
	}
	return items[0], nil
}

func (svc *ItemListService) FindItems(ctx context.Context, f todo.ItemFilter) ([]*todo.Item, error) {
	tx, err := svc.db.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	todos, err := findTodoItems(ctx, tx, f)
	if err != nil {
		return nil, err
	}

	return todos, nil
}

func findTodoItems(ctx context.Context, tx *Tx, f todo.ItemFilter) ([]*todo.Item, error) {
	user, err := todo.ValidUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	var args []interface{}
	where := []string{"1 = 1"}

	if v := f.ID; v != nil {
		where, args = append(where, fmt.Sprintf("id = $%d", len(where))), append(args, *v)
	}

	if v := f.ListID; v != nil {
		where, args = append(where, fmt.Sprintf("list_id = $%d", len(where))), append(args, *v)
	}

	if v := f.UserID; v != nil {
		where, args = append(where, fmt.Sprintf("user_id = $%d", len(where))), append(args, *v)
	}

	if v := f.Name; v != nil {
		where, args = append(where, fmt.Sprintf("name = $%d", len(where))), append(args, *v)
	}

	if v := f.Completed; v != nil {
		where, args = append(where, fmt.Sprintf("completed = $%d", len(where))), append(args, *v)
	}
	query := `
	SELECT 
		id, 
		user_id,
		list_id,
		name, 
		completed, 
		created_at, 
		updated_at 
	FROM items
	WHERE ` + strings.Join(where, " AND ") + `
	ORDER BY id ASC;`
	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	items := make([]*todo.Item, 0)
	for rows.Next() {
		var item todo.Item
		if err := rows.Scan(
			&item.ID,
			&item.UserID,
			&item.ListID,
			&item.Name,
			&item.Completed,
			(*Time)(&item.CreatedAt),
			(*Time)(&item.UpdatedAt),
		); err != nil {
			return nil, err
		}

		if item.UserID != user.ID {
			return nil, todo.Err(todo.EUNAUTHORIZED, "user %q cannot read item %q", user.ID, item.ID)
		}
		items = append(items, &item)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

func (svc *ItemListService) UpdateItem(ctx context.Context, id int, upd todo.ItemUpdate) (*todo.Item, error) {
	tx, err := svc.db.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	todo, err := updateTodoItem(ctx, tx, id, upd)
	if err != nil {
		return nil, err
	}
	return todo, tx.Commit()
}

func updateTodoItem(ctx context.Context, tx *Tx, id int, upd todo.ItemUpdate) (*todo.Item, error) {
	item, err := findTodoItem(ctx, tx, id)
	if err != nil {
		return nil, err
	}

	user, err := todo.ValidUserFromContext(ctx)
	if err != nil {
		return nil, err
	}

	item.UpdatedAt = tx.now
	if v := upd.Name; v != nil {
		item.Name = *v
	}
	if v := upd.Completed; v != nil {
		item.Completed = *v
	}

	if err = item.Validate(); err != nil {
		return item, err
	}

	if _, err := tx.ExecContext(ctx, `
	UPDATE items 
	SET name = $1,
		completed = $2,
		updated_at = $3
	WHERE id = $4 AND user_id = $5`,
		item.Name, item.Completed, (*Time)(&item.UpdatedAt), item.ID, user.ID); err != nil {
		return item, err
	}

	return item, nil
}

func (svc *ItemListService) CreateItem(ctx context.Context, item *todo.Item) error {
	tx, err := svc.db.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := createTodoItem(ctx, tx, item); err != nil {
		return err
	}
	return tx.Commit()
}

func createTodoItem(ctx context.Context, tx *Tx, item *todo.Item) error {
	list, err := findTodoListByID(ctx, tx, item.ListID)
	if err != nil {
		return err
	}

	item.UserID = list.UserID
	item.CreatedAt = tx.now
	item.UpdatedAt = item.CreatedAt

	if err := item.Validate(); err != nil {
		return err
	}

	var id int64
	err = tx.QueryRowContext(ctx, `
INSERT INTO items (name, user_id, list_id, completed, created_at, updated_at) 
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id`,
		item.Name,
		item.UserID,
		item.ListID,
		item.Completed,
		(*Time)(&item.CreatedAt),
		(*Time)(&item.UpdatedAt)).Scan(&id)
	if err != nil {
		return err
	}
	item.ID = int(id)

	return nil
}

func (svc *ItemListService) DeleteItem(ctx context.Context, id int) error {
	tx, err := svc.db.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err = deleteTodoItem(ctx, tx, id); err != nil {
		return err
	}
	return tx.Commit()

}

func deleteTodoItem(ctx context.Context, tx *Tx, id int) error {
	user, err := todo.ValidUserFromContext(ctx)
	if err != nil {
		return err
	}

	if id <= 0 {
		return todo.Err(todo.EINVALID, "invalid id")
	}

	result, err := tx.ExecContext(ctx, `DELETE FROM items WHERE id = $1 AND user_id = $2`, id, user.ID)
	if err != nil {
		return err
	} else if n, _ := result.RowsAffected(); n == 0 {
		return todo.Err(todo.ENOTFOUND, "could not delete item with id %v", id)
	}

	return nil
}
