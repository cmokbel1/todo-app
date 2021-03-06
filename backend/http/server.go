package http

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/cmokbel1/todo-app/backend/todo"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/httprate"
)

type Server struct {
	ln     net.Listener
	server *http.Server

	// config values
	Addr               string
	Domain             string
	TLS                bool
	APIKey             string
	CORSAllowedOrigins string
	// AssetsDirectory is the path to the frontend HTML/CSS/JavaScript
	AssetsDirectory string

	Logger todo.Logger
	// LoggerMiddleware is exposed for testing purposes.
	LoggerMiddleware func(http.Handler) http.Handler
	SessionManager   *scs.SessionManager
	ItemListService  todo.ItemListService
	UserService      todo.UserService
}

func NewServer() *Server {
	s := &Server{
		server: &http.Server{
			ReadTimeout:  time.Second * 6,
			WriteTimeout: time.Second * 6,
			IdleTimeout:  time.Second * 6,
		},
		Logger:           todo.NewLogger(),
		LoggerMiddleware: middleware.Logger,
	}
	return s
}

func (s *Server) Listen() (err error) {
	if s.APIKey == "" {
		s.Logger.Warn("API key is empty")
	}

	if s.CORSAllowedOrigins == "*" {
		s.Logger.Warn("CORS allowed origins is '*'")
	}

	r := chi.NewRouter()
	r.Use(middleware.Recoverer)
	// hard coded rate limit of 60 requests/minute/IP
	r.Use(httprate.LimitByIP(60, time.Minute))
	r.Use(s.LoggerMiddleware)
	r.Use(s.cors)
	r.Use(s.sessionMiddleware)
	r.Use(monitorMetrics)
	r.Use(middleware.StripSlashes)

	r.Route("/api", func(r chi.Router) {
		s.registerTodoRoutes(r)
		s.registerUserRoutes(r)
		s.registerBuildRoute(r)
	})

	if s.AssetsDirectory != "" {
		s.Logger.Infof("serving assets out of %q", s.AssetsDirectory)
		r.Get("/*", s.assetsHandler(s.AssetsDirectory, "index.html", "asset-manifest.json", "manifest.json"))
	}

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
		if user := todo.UserFromContext(r.Context()); user == nil {
			next.ServeHTTP(w, r)
			return
		}
		s.Logger.Debug("requireNoAuth user should not be authd to access this route")
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
		s.Logger.Debugf("requireAuth user is not authorized")
		s.error(w, r, todo.Unauthorized)
		return
	})
}

// cors is middleware that enables CORS for localhost domains only. This functionality works only if the configured
// Domain is localhost.
func (s *Server) cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		headers := w.Header()
		if s.Domain == "localhost" {
			headers.Set("Access-Control-Allow-Origin", s.CORSAllowedOrigins)
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

func (s *Server) assetsHandler(dir string, allowed ...string) http.HandlerFunc {
	fs := http.FileServer(http.Dir(dir))
	return func(w http.ResponseWriter, r *http.Request) {
		for _, path := range allowed {
			if r.URL.Path == path {
				fs.ServeHTTP(w, r)
				return
			}
		}
		if r.URL.Path == "/" || r.URL.Path == "/index.html" || r.URL.Path == "/favicon.ico" {
			fs.ServeHTTP(w, r)
			return //
		}

		// Prevent directory browsing
		if strings.HasSuffix(r.URL.Path, "/") {
			http.NotFound(w, r)
			return
		}

		fs.ServeHTTP(w, r)
	}
}
