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

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
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
	// FindUsers finds one or more Users who match the UserFilter.
	FindUsers(ctx context.Context, f UserFilter) ([]*User, error)
}

type UserFilter struct {
	// Filter fields
	ID    *int    `json:"id"`
	Name  *string `json:"name"`
	Email *string `json:"email"`

	// Range restrictions
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

type UserUpdate struct {
	Name        *string `json:"name"`
	Email       *string `json:"email"`
	DisplayName *string `json:"displayName"`
	Password    *string `json:"password"`
}

func (upd UserUpdate) Validate() error {
	if upd.Name == nil &&
		upd.Email == nil &&
		upd.DisplayName == nil &&
		upd.Password == nil {
		return Err(EINVALID, "one of name, email, password, or display name is required")
	} else if upd.Name != nil && *upd.Name == "" {
		return Err(EINVALID, "name cannot be empty")
	} else if upd.DisplayName != nil && *upd.DisplayName == "" {
		return Err(EINVALID, "display name cannot be empty")
	} else if upd.Password != nil && *upd.Password == "" {
		return Err(EINVALID, "password cannot be empty")
	}

	return nil
}

type UserCredentials struct {
	// ID is the unique identifier for this set of credentials.
	ID int `json:"id"`
	// UserID is the ID of the User who these credentials belong to.
	UserID int `json:"userId"`
	// Name is the login name.
	Name string `json:"name"`
	// Password is a hashed password.
	Password string `json:"password,omitempty"`
	// APIKey used to impersonate this User when using the CLI, generated randomly when the
	// user is first created.
	APIKey string `json:"-"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (uc *UserCredentials) Validate() error {
	if uc.ID <= 0 {
		return Err(EINVALID, "id required")
	} else if uc.UserID <= 0 {
		return Err(EINVALID, "user id required")
	} else if uc.Name == "" {
		return Err(EINVALID, "name required")
	} else if uc.Password == "" {
		return Err(EINVALID, "password required")
	} else if uc.APIKey == "" {
		return Err(EINVALID, "api key required")
	}
	return nil
}

type UserCredentialsService interface {
	// FindByName finds a UserCredentials by name.
	FindByName(ctx context.Context, name string) (*UserCredentials, error)
	// FindByUserID finds a UserCredentials by user ID.
	FindByUserID(ctx context.Context, id string) (*UserCredentials, error)
	// CreateUserCredentials creates a UserCredentials.
	CreateUserCredentials(ctx context.Context, uc *UserCredentials) error
}
