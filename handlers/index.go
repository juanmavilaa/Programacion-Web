package handlers

import (
	"net/http"

	db "reps/db/sqlc"
	"reps/views"
)

type IndexHandler struct {
	queries *db.Queries
}

func NewIndexHandler(q *db.Queries) *IndexHandler {
	return &IndexHandler{queries: q}
}

func (h *IndexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/":
		h.Index(w, r)
	default:
		http.NotFound(w, r)
	}
}

func (h *IndexHandler) Index(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	clients, err := h.queries.ListClients(r.Context())
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	workouts, err := h.queries.ListWorkouts(r.Context())
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if err := views.IndexPage(clients, workouts).Render(r.Context(), w); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}
