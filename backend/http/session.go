package http

import (
	"bytes"
	"context"
	"encoding/json"
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
	mgr.Codec = jsonSessionCodec{}
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

// jsonSessionCodec is used to serialize the Session object as a JSON.
type jsonSessionCodec struct{}

func (jsonSessionCodec) Encode(deadline time.Time, values map[string]interface{}) ([]byte, error) {
	aux := &struct {
		Values   map[string]interface{} `json:"values"`
		Deadline time.Time              `json:"deadline"`
	}{
		Deadline: deadline,
		Values:   values,
	}

	var b bytes.Buffer
	if err := json.NewEncoder(&b).Encode(&aux); err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

// Decode converts a byte slice into a session deadline and values.
func (jsonSessionCodec) Decode(b []byte) (time.Time, map[string]interface{}, error) {
	aux := &struct {
		Values   map[string]interface{} `json:"values"`
		Deadline time.Time              `json:"deadline"`
	}{}

	r := bytes.NewReader(b)
	if err := json.NewDecoder(r).Decode(&aux); err != nil {
		return time.Time{}, nil, err
	}

	return aux.Deadline, aux.Values, nil
}
