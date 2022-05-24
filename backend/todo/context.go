package todo

import "context"

// contextKey represents an internal key for adding context fields.
// This is considered best practice as it prevents other packages from
// interfering with our context keys.
type contextKey int

// userContextKey stores the current logged-in user in the context.
const userContextKey = contextKey(iota + 1)

// NewContextWithUser returns a new context with the given user.
func NewContextWithUser(ctx context.Context, user *User) context.Context {
	return context.WithValue(ctx, userContextKey, user)
}

// UserFromContext returns the current logged-in user.
func UserFromContext(ctx context.Context) *User {
	user, _ := ctx.Value(userContextKey).(*User)
	return user
}

// ValidUserFromContext returns the current logged-in user if they pass basic validation,
// otherwise, an ErrUnauthorized is returned.
func ValidUserFromContext(ctx context.Context) (*User, error) {
	if ctx == nil {
		return nil, Unauthorized
	}

	user := UserFromContext(ctx)
	if user == nil {
		return nil, Unauthorized
	} else if user.ID <= 0 {
		return nil, Unauthorized
	} else if err := user.Validate(); err != nil {
		return nil, Unauthorized
	}

	return user, nil
}
