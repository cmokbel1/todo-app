package http

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/cmokbel1/todo-app/backend/todo"
)

func NewSessionManager() *scs.SessionManager {
	mgr := scs.New()
	mgr.Lifetime = time.Hour * 24
	mgr.IdleTimeout = time.Minute * 20
	mgr.Cookie.Name = "session"
	mgr.Cookie.Path = "/"
	mgr.Cookie.HttpOnly = true
	return mgr
}

func (s *Server) CreateSession(ctx context.Context, user *todo.User) error {
	s.SessionManager.Put(ctx, "user", *user)
	return nil
}

func (s *Server) RenewSession(ctx context.Context) error {
	if err := s.SessionManager.RenewToken(ctx); err != nil {
		return todo.Err(todo.EINTERNAL, "failed to renew session token: %v", err)
	}
	return nil
}

func (s *Server) DestroySession(ctx context.Context) error {
	if err := s.SessionManager.Destroy(ctx); err != nil {
		return todo.Err(todo.EINTERNAL, "failed to destroy session: %v", err)
	}
	return nil
}

// sessionMiddleware populates the context with a user session from either a Cookie or from the owner of the
// api key specified in the Authorization header.
func (s *Server) sessionMiddleware(next http.Handler) http.Handler {
	return s.SessionManager.LoadAndSave(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		if user, ok := s.SessionManager.Get(ctx, "user").(todo.User); ok {
			s.Logger.Debugf("sessionMiddleware found user %q", user.Name)
			ctx = todo.NewContextWithUser(ctx, &user)
		} else if h := r.Header.Get("Authorization"); h != "" {
			if token := strings.TrimPrefix(h, "Bearer "); token != "" {
				s.Logger.Debug("beginning token auth")
				if user, err := s.UserService.FindUserByAPIKey(ctx, token); err != nil {
					s.error(w, r, todo.Err(todo.EUNAUTHORIZED, "invalid credentials"))
					return
				} else {
					ctx = todo.NewContextWithUser(ctx, user)
				}
			}
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	}))
}
