package todo

import (
	"context"
	"time"
)

// Auth is an auth context associated with a User. A User may have more than one
// Auth associated with them but never more than one Auth per Source.
// An Auth is represented as either being a set of UserCredentials or OAuth credentials.
//
// In the event that the Auth Source is 'app' then AccessToken and RefreshToken will
// both be empty.
type Auth struct {
	// ID represents the ID of the Auth context.
	ID string `json:"id"`

	// UserID is the User ID associated with this Auth.
	UserID int `json:"userId"`
	// User is the User associated with this Auth.
	User *User `json:"user"`

	// Source is the origin of this Auth. It can be one of: "app", "github"
	Source string `json:"source"`
	// SourceID is the ID of the User associated with the Auth at the Source, e.g. the GitHub account's user ID.
	SourceID string `json:"sourceId"`

	// AccessToken is the OAuth2 access token
	AccessToken string `json:"accessToken"`
	// RefreshToken is the OAuth2 refresh token
	RefreshToken string `json:"refreshToken"`

	// Expiry is the expiry time for this Auth.
	Expiry time.Time `json:"expiry"`

	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

func (a *Auth) Validate() error {
	if a.UserID == 0 {
		return Err(EINVALID, "user id required.")
	} else if a.Source == "" {
		return Err(EINVALID, "source required.")
	} else if a.SourceID == "" {
		return Err(EINVALID, "source id required.")
	} else if a.AccessToken == "" {
		return Err(EINVALID, "access token required.")
	}
	return nil
}

type AuthService interface {
	// FindAuths finds one or more Auth objects according to the AuthFilter provided.
	FindAuths(ctx context.Context, f AuthFilter) ([]*Auth, error)
	// CreateAuth creates an Auth. If a User is attached to the Auth, then
	// the Auth will be linked to the User, otherwise, a new User will be created.
	CreateAuth(ctx context.Context, auth *Auth) error
	// DeleteAuth removes the Auth but does not remove the associated User.
	DeleteAuth(ctx context.Context, id int) error
}

type AuthFilter struct {
	// Filter fields
	ID     *int `json:"id"`
	UserID *int `json:"userId"`

	// Range restrictions
	Offset int `json:"offset"`
	Limit  int `json:"limit"`
}

func (f AuthFilter) Validate() error {
	if f.ID == nil && f.UserID == nil {
		return Err(EINVALID, "one of id or user id is required")
	} else if f.ID != nil && *f.ID <= 0 {
		return Err(EINVALID, "id is invalid")
	} else if f.UserID != nil && *f.UserID <= 0 {
		return Err(EINVALID, "user id is invalid")
	}
	return nil
}
