package http

import (
	"encoding/json"
	"net/http"

	"github.com/cmokbel1/todo-app/backend/todo"
	"github.com/go-chi/chi"
)

func (s *Server) registerTodoRoutes(r chi.Router) {
	r.Route("/todos", func(r chi.Router) {
		r.Use(s.requireAuth)
		r.Get("/", s.handleTodoListIndex)
		r.Post("/", s.handleTodoListCreate)
		r.Route("/{id}", func(r chi.Router) {
			r.Use(s.requireIntParam("id"))
			r.Get("/", s.handleTodoListGet)
			r.Patch("/", s.handleTodoListEdit)
			r.Delete("/", s.handleTodoListDelete)
			r.Post("/", s.handleTodoItemCreate)
			r.Route("/{itemID}", func(r chi.Router) {
				r.Use(s.requireIntParam("itemID"))
				r.Get("/", s.handleTodoItemGet)
				r.Patch("/", s.handleTodoItemEdit)
				r.Delete("/", s.handleTodoItemDelete)
			})
		})
	})
}

func (s *Server) handleTodoListIndex(w http.ResponseWriter, r *http.Request) {
	user, err := todo.ValidUserFromContext(r.Context())
	if err != nil {
		s.error(w, r, err)
		return
	}

	lists, err := s.ItemListService.FindLists(r.Context(), todo.ListFilter{UserID: &user.ID})
	if err != nil {
		s.error(w, r, err)
		return
	}
	s.json(w, r, http.StatusOK, lists)
}

func (s *Server) handleTodoListCreate(w http.ResponseWriter, r *http.Request) {
	list := &todo.List{}
	if err := json.NewDecoder(r.Body).Decode(list); err != nil {
		s.error(w, r, err)
		return
	}

	if err := s.ItemListService.CreateList(r.Context(), list); err != nil {
		s.error(w, r, err)
		return
	}
	s.json(w, r, http.StatusCreated, list)
}

func (s *Server) handleTodoListGet(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	list, err := s.ItemListService.FindListByID(r.Context(), id)
	if err != nil {
		s.error(w, r, err)
		return
	}
	s.json(w, r, http.StatusOK, list)
}

func (s *Server) handleTodoListDelete(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	if err := s.ItemListService.DeleteList(r.Context(), id); err != nil {
		s.error(w, r, err)
		return
	}
	s.json(w, r, http.StatusNoContent, nil)
}

func (s *Server) handleTodoListEdit(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("id").(int)
	var req todo.ListUpdate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.error(w, r, err)
		return
	}

	list, err := s.ItemListService.UpdateList(r.Context(), id, req)
	if err != nil {
		s.error(w, r, err)
		return
	}
	s.json(w, r, http.StatusOK, list)
}

func (s *Server) handleTodoItemCreate(w http.ResponseWriter, r *http.Request) {
	item := &todo.Item{}
	if err := json.NewDecoder(r.Body).Decode(item); err != nil {
		s.error(w, r, err)
		return
	}

	item.ListID = r.Context().Value("id").(int)

	if err := s.ItemListService.CreateItem(r.Context(), item); err != nil {
		s.error(w, r, err)
		return
	}

	s.json(w, r, http.StatusCreated, item)
}

func (s *Server) handleTodoItemEdit(w http.ResponseWriter, r *http.Request) {
	var req todo.ItemUpdate
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		s.error(w, r, err)
		return
	}

	itemID := r.Context().Value("itemID").(int)

	if item, err := s.ItemListService.UpdateItem(r.Context(), itemID, req); err != nil {
		s.error(w, r, err)
	} else {
		s.json(w, r, http.StatusOK, item)
	}
}

func (s *Server) handleTodoItemGet(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("itemID").(int)
	item, err := s.ItemListService.FindItemByID(r.Context(), id)
	if err != nil {
		s.error(w, r, err)
		return
	}
	s.json(w, r, http.StatusOK, item)
}

func (s *Server) handleTodoItemDelete(w http.ResponseWriter, r *http.Request) {
	id := r.Context().Value("itemID").(int)
	if err := s.ItemListService.DeleteItem(r.Context(), id); err != nil {
		s.error(w, r, err)
		return
	}
	s.json(w, r, http.StatusNoContent, nil)
}
