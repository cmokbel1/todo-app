package todo

import (
	"context"
	"time"
)

type User struct {
	// ID is the unique identifier for this User.
	ID int `json:"id"`
	// Name represents the User's username
	Name string `json:"name"`
	// Email represents the email address associated with this User.
	Email *string `json:"email,omitempty"`
	// APIKey for bypassing normal auth flow access.
	APIKey string `json:"-"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`

	// Auths is a current list of associated auths
	Auths []*Auth `json:"auths"`
}

func (u *User) Validate() error {
	if u.ID <= 0 {
		return Err(EINVALID, "id required")
	} else if u.Name == "" {
		return Err(EINVALID, "name required")
	}
	return nil
}

type UserService interface {
	// CreateUser creates a User and an attached UserLogin with a random password.
	CreateUser(ctx context.Context, user *User) error
	// DeleteUser deletes a User, their UserCredentials and all associated Auths.
	DeleteUser(ctx context.Context, id int) error
	// UpdateUser updates a User.
	UpdateUser(ctx context.Context, id int, upd UserUpdate) (*User, error)
	// FindUserByID finds a User by their User ID.
	FindUserByID(ctx context.Context, id int) (*User, error)
	// FindUserByName finds a User by their Name.
	FindUserByName(ctx context.Context, name string) (*User, error)
	// FindUserByAPIKey finds a User by their API key.
	FindUserByAPIKey(ctx context.Context, apiKey string) (*User, error)
	// FindUsers finds one or more Users who match the UserFilter.
	FindUsers(ctx context.Context, f UserFilter) ([]*User, error)
}

type UserFilter struct {
	// Filter fields
	ID     *int    `json:"id"`
	Name   *string `json:"name"`
	Email  *string `json:"email"`
	APIKey *string `json:"apiKey"`

	// Range restrictions
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type UserUpdate struct {
	Name  *string `json:"name"`
	Email *string `json:"email"`
}

func (upd UserUpdate) Validate() error {
	if upd.Name == nil && upd.Email == nil {
		return Err(EINVALID, "one of name or email is required")
	} else if upd.Name != nil && *upd.Name == "" {
		return Err(EINVALID, "name cannot be empty")
	}

	return nil
}
