package handlers

import (
	"net/http"
	"net/url"

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

// GET / - Pagina principal
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

	errorMsg := r.URL.Query().Get("error")

	if err := views.IndexPage(clients, workouts, errorMsg).Render(r.Context(), w); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

// vuelve a / con el mensaje de error
func redirectConError(w http.ResponseWriter, r *http.Request, mensaje string) {
	http.Redirect(w, r, "/?error="+url.QueryEscape(mensaje), http.StatusSeeOther)
}
