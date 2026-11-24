package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	db "reps/db/sqlc"
	"reps/views"
)

type ClientHandler struct {
	queries *db.Queries
}

func NewClientHandler(q *db.Queries) *ClientHandler {
	return &ClientHandler{queries: q}
}

func ValidateEmail(email string) error {
	re := regexp.MustCompile(
		`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !re.MatchString(email) {
		return errors.New("formato de email inválido")
	}
	return nil
}

func validarCliente(name, email string) error {
	if name == "" {
		return errors.New("falta el nombre")
	}
	return ValidateEmail(email)
}

func (h *ClientHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/clients" {
		switch r.Method {
		case http.MethodPost:
			h.Create(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		http.Error(w, "url invalida", http.StatusBadRequest)
		return
	}
	id, err := strconv.ParseInt(parts[2], 10, 32)
	if err != nil || id <= 0 {
		http.Error(w, "id invalido", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodDelete:
		h.Delete(w, r, int32(id))
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ClientHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "formulario invalido", http.StatusBadRequest)
		return
	}

	params := db.CreateClientParams{
		Name:  strings.TrimSpace(r.FormValue("name")),
		Email: strings.TrimSpace(r.FormValue("email")),
	}

	if err := validarCliente(params.Name, params.Email); err != nil {
		mostrarError(w, r, http.StatusUnprocessableEntity, err.Error())
		return
	}

	_, err := h.queries.GetClientByEmail(r.Context(), params.Email)
	if err == nil {
		mostrarError(w, r, http.StatusConflict, "ya hay un cliente con ese email")
		return
	}
	if err != sql.ErrNoRows {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if _, err := h.queries.CreateClient(r.Context(), params); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	clients, err := h.queries.ListClients(r.Context())
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if err := views.ClientCreated(clients).Render(r.Context(), w); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func (h *ClientHandler) Delete(w http.ResponseWriter, r *http.Request, id int32) {
	// no dejamos borrar un cliente que tiene entrenamientos cargados
	workouts, err := h.queries.ListWorkoutsByClient(r.Context(), id)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if len(workouts) > 0 {
		mostrarError(w, r, http.StatusConflict, "el cliente tiene entrenamientos cargados")
		return
	}

	if err := h.queries.DeleteClient(r.Context(), id); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
