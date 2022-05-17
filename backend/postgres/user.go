package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/cmokbel1/todo-app/backend/crypto"
	"github.com/cmokbel1/todo-app/backend/todo"
)

var _ todo.UserService = (*UserService)(nil)

func NewUserService(db *DB) *UserService {
	return &UserService{
		db: db,
	}
}

type UserService struct {
	db *DB
}

func (svc *UserService) CreateUser(ctx context.Context, user *todo.User) error {
	tx, err := svc.db.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err = createUser(ctx, tx, user); err != nil {
		return fmt.Errorf("postgres create user: %w", err)
	}

	return tx.Commit()
}

func createUser(ctx context.Context, tx *Tx, user *todo.User) (err error) {
	user.CreatedAt = tx.now
	user.UpdatedAt = user.CreatedAt
	if user.Name == "" {
		return todo.Err(todo.EINVALID, "name is required")
	}
	user.APIKey = crypto.RandomString()

	// create the user if they don't already exist
	if other, err := findUserByName(ctx, tx, user.Name); err != nil && !errors.Is(err, todo.NotFound) {
		return err
	} else if other != nil {
		return todo.Err(todo.ECONFLICT, "name %q is taken", user.Name)
	}

	var id int64
	err = tx.QueryRowContext(ctx, `
INSERT INTO users (name, email, api_key, created_at, updated_at) 
VALUES ($1,$2,$3,$4,$5)
RETURNING id`,
		user.Name,
		user.Email,
		user.APIKey,
		user.CreatedAt,
		user.UpdatedAt).Scan(&id)
	if err != nil {
		return err
	}
	user.ID = int(id)

	return nil
}

func (svc *UserService) DeleteUser(ctx context.Context, id int) error {
	tx, err := svc.db.BeginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err = deleteUser(ctx, tx, id); err != nil {
		return fmt.Errorf("postgres delete user: %w", err)
	}

	return tx.Commit()
}

func (svc *UserService) UpdateUser(ctx context.Context, id int, upd todo.UserUpdate) (*todo.User, error) {
	tx, err := svc.db.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	user, err := updateUser(ctx, tx, id, upd)
	if err != nil {
		return nil, fmt.Errorf("postgres update user: %w", err)
	}

	return user, tx.Commit()
}

func (svc *UserService) FindUsers(ctx context.Context, f todo.UserFilter) ([]*todo.User, error) {
	tx, err := svc.db.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	users, err := findUsers(ctx, tx, f)
	if err != nil {
		return nil, err
	}

	return users, tx.Commit()
}

func (svc *UserService) FindUserByID(ctx context.Context, id int) (*todo.User, error) {
	tx, err := svc.db.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	user, err := findUserByID(ctx, tx, id)
	if err != nil {
		return nil, err
	}
	return user, tx.Commit()
}

func (svc *UserService) FindUserByName(ctx context.Context, name string) (*todo.User, error) {
	tx, err := svc.db.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	user, err := findUserByName(ctx, tx, name)
	if err != nil {
		return nil, err
	}
	return user, tx.Commit()
}

func (svc *UserService) FindUserByAPIKey(ctx context.Context, apiKey string) (*todo.User, error) {
	tx, err := svc.db.BeginTx(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	user, err := findUserByAPIKey(ctx, tx, apiKey)
	if err != nil {
		return nil, err
	}
	return user, tx.Commit()
}

func deleteUser(ctx context.Context, tx *Tx, id int) error {
	result, err := tx.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return err
	} else if n, _ := result.RowsAffected(); n == 0 {
		return todo.Err(todo.ENOTFOUND, "could not delete user with id %v", id)
	}
	return nil
}

func updateUser(ctx context.Context, tx *Tx, id int, upd todo.UserUpdate) (*todo.User, error) {
	wrap := func(err error) error { return fmt.Errorf("updateUser: %w", err) }
	user, err := findUserByID(ctx, tx, id)
	if err != nil {
		return nil, wrap(err)
	}

	user.UpdatedAt = tx.now
	if v := upd.Email; v != nil {
		user.Email = v
	}

	if v := upd.Name; v != nil {
		user.Name = *v
	}

	if _, err := tx.ExecContext(ctx, `
	UPDATE users 
	SET name = $1,
	    email = $2,
	    updated_at = $3
	WHERE id = $4`,
		user.Name, user.Email, user.UpdatedAt, id); err != nil {
		return nil, err
	}

	return user, nil
}

func findUserByName(ctx context.Context, tx *Tx, name string) (*todo.User, error) {
	name = strings.ToLower(name)
	users, err := findUsers(ctx, tx, todo.UserFilter{Name: &name})
	if err != nil {
		return nil, err
	} else if len(users) == 0 {
		return nil, todo.Err(todo.ENOTFOUND, "could not find user with name %q", name)
	}
	return users[0], nil
}

func findUserByID(ctx context.Context, tx *Tx, id int) (*todo.User, error) {
	users, err := findUsers(ctx, tx, todo.UserFilter{ID: &id})
	if err != nil {
		return nil, err
	} else if len(users) == 0 {
		return nil, todo.Err(todo.ENOTFOUND, "could not find user with id %q", id)
	}
	return users[0], nil
}

func findUserByAPIKey(ctx context.Context, tx *Tx, apiKey string) (*todo.User, error) {
	users, err := findUsers(ctx, tx, todo.UserFilter{APIKey: &apiKey})
	if err != nil {
		return nil, err
	} else if len(users) == 0 {
		return nil, todo.Err(todo.ENOTFOUND, "could not find user with api key")
	}
	return users[0], nil
}

func findUsers(ctx context.Context, tx *Tx, f todo.UserFilter) ([]*todo.User, error) {
	var args []interface{}
	where := []string{"1 = 1"}

	if v := f.ID; v != nil {
		where, args = append(where, fmt.Sprintf("id = $%d", len(where))), append(args, *v)
	}

	if v := f.Name; v != nil {
		low := strings.ToLower(*v)
		where, args = append(where, fmt.Sprintf("LOWER(name) = $%d", len(where))), append(args, low)
	}

	if v := f.Email; v != nil {
		low := strings.ToLower(*v)
		where, args = append(where, fmt.Sprintf("LOWER(email) = $%d", len(where))), append(args, low)
	}

	if v := f.APIKey; v != nil {
		where, args = append(where, fmt.Sprintf("api_key = $%d", len(where))), append(args, *v)
	}

	query := `
	SELECT 
		id,
		name, 
		email, 
		api_key,
		created_at, 
		updated_at
	FROM users
	WHERE ` + strings.Join(where, " AND ") + `
	ORDER BY id ASC ` + FormatLimitOffset(f.Limit, f.Offset)
	rows, err := tx.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*todo.User
	for rows.Next() {
		var user todo.User
		if err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.APIKey,
			(*Time)(&user.CreatedAt),
			(*Time)(&user.UpdatedAt),
		); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
