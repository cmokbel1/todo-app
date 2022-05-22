package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/cmokbel1/todo-app/backend/todo"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

type Server struct {
	ln     net.Listener
	server *http.Server

	// config values
	Addr   string
	Domain string
	TLS    bool
	APIKey string

	Logger todo.Logger
	// LoggerMiddleware is exposed for testing purposes.
	LoggerMiddleware func(http.Handler) http.Handler
	SessionManager   *scs.SessionManager
	ItemListService  todo.ItemListService
	UserService      todo.UserService
}

func NewServer() *Server {
	s := &Server{
		server:           &http.Server{},
		Logger:           todo.NewLogger(),
		LoggerMiddleware: middleware.Logger,
	}
	return s
}

func (s *Server) Listen() (err error) {
	if s.APIKey == "" {
		s.Logger.Warn("API key is empty")
	}

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	r.Use(s.LoggerMiddleware)
	r.Use(s.cors)
	r.Use(s.sessionMiddleware)
	r.Use(monitorMetrics)
	r.Use(middleware.StripSlashes)

	r.Route("/api", func(r chi.Router) {
		// TODO(1gm): add these routes when they are created
		// s.registerTodoRoutes(r)
		s.registerUserRoutes(r)
		s.registerBuildRoute(r)
	})
	r.NotFound(s.notFound)
	s.server.Handler = r

	if s.ln, err = net.Listen("tcp", s.Addr); err != nil {
		return err
	}

	go s.server.Serve(s.ln)
	return nil
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	return s.server.Shutdown(ctx)
}

func (s *Server) Scheme() string {
	if s.TLS {
		return "https"
	}
	return "http"
}

// Port returns the TCP port for the running server.
// This is useful in tests where we allocate a random port by using ":0".
func (s *Server) Port() int {
	if s.ln == nil {
		return 0
	}
	return s.ln.Addr().(*net.TCPAddr).Port
}

// URL returns the local base URL of the running server.
func (s *Server) URL() string {
	scheme, port := s.Scheme(), s.Port()

	// Use localhost unless a domain is specified.
	domain := "localhost"
	if s.Domain != "" {
		domain = s.Domain
	}

	// Return without port if using standard ports.
	if (scheme == "http" && port == 80) || (scheme == "https" && port == 443) {
		return fmt.Sprintf("%s://%s", s.Scheme(), domain)
	}
	return fmt.Sprintf("%s://%s:%d", s.Scheme(), domain, s.Port())
}

func (s *Server) registerBuildRoute(r chi.Router) {
	r.Get("/build", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(todo.BuildJSON())
	})
}

func (s *Server) notFound(w http.ResponseWriter, r *http.Request) {
	s.error(w, r, todo.Err(todo.ENOTFOUND, "not found"))
}

func (s *Server) error(w http.ResponseWriter, r *http.Request, err error) {
	code, msg := StatusCode(err), todo.ErrMessage(err)
	if code == http.StatusInternalServerError {
		s.Logger.E(err)
	} else if code == http.StatusUnauthorized {
		s.Logger.Warn(msg)
		msg = todo.Unauthorized.Message
	} else {
		s.Logger.Info(err.Error())
	}
	s.json(w, r, code, &Error{Message: msg})
}

func (s *Server) json(w http.ResponseWriter, _ *http.Request, code int, body interface{}) {
	w.WriteHeader(code)
	if code == http.StatusNoContent {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(body); err != nil {
		s.Logger.E(err)
	}
}

func (s *Server) requireIntParam(key string) func(next http.Handler) http.Handler {
	requiredErr := &todo.Error{Code: todo.EINVALID, Message: key + " required"}
	parseErr := &todo.Error{Code: todo.EINVALID, Message: key + " must be an integer"}
	fn := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if param := chi.URLParam(r, key); param == "" {
				s.error(w, r, requiredErr)
				return
			} else if val, err := strconv.ParseInt(param, 10, 64); err != nil {
				s.error(w, r, parseErr)
				return
			} else {
				ctx := context.WithValue(r.Context(), key, int(val))
				r = r.WithContext(ctx)
			}
			next.ServeHTTP(w, r)
		})
	}
	return fn
}

func (s *Server) requireAPIKey(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if apiKey := r.Header.Get("Todo-Api-Key"); apiKey != s.APIKey {
			s.Logger.Warnf("invalid API key originating from %v to %v", r.RemoteAddr, r.URL)
			s.error(w, r, todo.NotFound)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func (s *Server) requireNoAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if user := todo.UserFromContext(r.Context()); user != nil {
			next.ServeHTTP(w, r)
			return
		}
		s.error(w, r, todo.Unauthorized)
		return
	})
}

func (s *Server) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if user := todo.UserFromContext(r.Context()); user != nil {
			next.ServeHTTP(w, r)
			return
		}
		s.error(w, r, todo.Unauthorized)
		return
	})
}

func (s *Server) cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headers := w.Header()
		if s.Domain == "localhost" {
			// TODO(1gm): Let's make this a config setting for dev environment.
			headers.Set("Access-Control-Allow-Origin", "http://localhost:5000")
			headers.Set("Access-Control-Allow-Headers", "Content-Type, Origin, Accept")
			headers.Set("Access-Control-Expose-Headers", "Link")
			headers.Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, OPTIONS, DELETE")
			headers.Set("Access-Control-Allow-Credentials", "true")
		}
		headers.Set("Vary", "Origin")
		if r.Method == "OPTIONS" {
			headers.Add("Vary", "Access-Control-Request-Method")
			headers.Add("Vary", "Access-Control-Request-Headers")
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
