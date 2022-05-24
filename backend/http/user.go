package http

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/cmokbel1/todo-app/backend/todo"
	"github.com/go-chi/chi"
)

func (s *Server) registerUserRoutes(r chi.Router) {
	r.With(s.requireAuth).Get("/user", s.handleMe)
	r.With(s.requireAuth).Get("/user/key", s.handleApiKey)
	r.With(s.requireNoAuth).Post("/user/login", s.handleLogin)
	r.With(s.requireAuth).Delete("/user/logout", s.handleLogout)

	r.With(s.requireAPIKey).Get("/users", s.handleUsersIndex)
	r.With(s.requireAPIKey).Post("/users", s.handleUserCreate)
	r.Route("/users/{id}", func(r chi.Router) {
		r.Use(s.requireIntParam("id"))
		r.With(s.requireAPIKey).Delete("/", s.handleUserDelete)
		r.With(s.requireAuth).Patch("/", s.handleUserUpdate)
	})
}

func (s *Server) handleMe(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user := todo.UserFromContext(ctx)
	latest, err := s.UserService.FindUserByID(ctx, user.ID)
	if err != nil {
		s.error(w, r, err)
		return
	}
	latest.Password = ""
	s.json(w, r, http.StatusOK, latest)
}

func (s *Server) handleApiKey(w http.ResponseWriter, r *http.Request) {
	s.Logger.Debug("in handleApiKey")
	ctx := r.Context()
	user := todo.UserFromContext(ctx)
	latest, err := s.UserService.FindUserByID(ctx, user.ID)
	if err != nil {
		s.error(w, r, err)
		return
	}
	s.json(w, r, http.StatusOK, latest.APIKey)
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var user *todo.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		s.Logger.E(err)
		s.error(w, r, err)
		return
	}

	if err := s.UserService.LoginUser(r.Context(), user); err != nil {
		s.Logger.E(err)
		s.error(w, r, fmt.Errorf("invalid credentials: %w", err))
		return
	}

	if err := s.CreateSession(r.Context(), user); err != nil {
		s.error(w, r, err)
		return
	}

	// do not render the password
	user.Password = ""
	s.json(w, r, http.StatusOK, *user)
	return
}

func (s *Server) handleLogout(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	user, err := todo.ValidUserFromContext(ctx)
	if err != nil {
		s.error(w, r, err)
		return
	}

	if err := s.DestroySession(ctx); err != nil {
		s.Logger.Errorf("failed to destroy user session for user %d: %v", user.ID, err)
		s.error(w, r, err)
		return
	}
	s.json(w, r, http.StatusNoContent, nil)
}

func (s *Server) handleUsersIndex(w http.ResponseWriter, r *http.Request) {
	users, err := s.UserService.FindUsers(r.Context(), todo.UserFilter{})
	if err != nil {
		s.error(w, r, err)
	}
	s.json(w, r, http.StatusCreated, users)
}

func (s *Server) handleUserCreate(w http.ResponseWriter, r *http.Request) {
	var user *todo.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		s.error(w, r, err)
		return
	}

	s.Logger.Infof("received request to create user %q", user.Name)
	if err := s.UserService.CreateUser(r.Context(), user); err != nil {
		s.error(w, r, err)
		return
	}
	s.Logger.Infof("created user %q (id = %q)", user.Name, user.ID)
	s.json(w, r, http.StatusCreated, user)
}

func (s *Server) handleUserDelete(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	if err := s.UserService.DeleteUser(r.Context(), id); err != nil {
		s.error(w, r, err)
		return
	}

	s.json(w, r, http.StatusNoContent, nil)
}

func (s *Server) handleUserUpdate(w http.ResponseWriter, r *http.Request) {
	var upd todo.UserUpdate
	if err := json.NewDecoder(r.Body).Decode(&upd); err != nil {
		s.error(w, r, err)
		return
	}

	ctx := r.Context()
	id := ctx.Value("id").(int)
	user, err := todo.ValidUserFromContext(ctx)
	if err != nil {
		s.error(w, r, err)
		return
	}

	if id != user.ID {
		s.error(w, r, todo.Unauthorized)
		return
	}

	user, err = s.UserService.UpdateUser(ctx, id, upd)
	if err != nil {
		s.error(w, r, err)
		return
	}

	if err = s.CreateSession(ctx, user); err != nil {
		s.error(w, r, err)
		return
	}
	s.json(w, r, http.StatusOK, user)
}
