package handlers

import (
	"database/sql"
	"errors"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	db "reps/db/sqlc"
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
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	switch r.URL.Path {
	case "/clients":
		h.Create(w, r)
	case "/clients/delete":
		h.Delete(w, r)
	default:
		http.NotFound(w, r)
	}
}

// POST /clients - Crear nuevo cliente
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
		redirectConError(w, r, err.Error())
		return
	}

	_, err := h.queries.GetClientByEmail(r.Context(), params.Email)
	if err == nil {
		redirectConError(w, r, "ya hay un cliente con ese email")
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

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// POST /clients/delete - Eliminar cliente
func (h *ClientHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "formulario invalido", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(r.FormValue("id"), 10, 32)
	if err != nil {
		http.Error(w, "id invalido", http.StatusBadRequest)
		return
	}

	_, err = h.queries.GetClient(r.Context(), int32(id))
	if err == sql.ErrNoRows {
		redirectConError(w, r, "cliente no encontrado")
		return
	}
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	// no dejamos borrar un cliente que tiene entrenamientos cargados
	workouts, err := h.queries.ListWorkoutsByClient(r.Context(), int32(id))
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
	if len(workouts) > 0 {
		redirectConError(w, r, "el cliente tiene entrenamientos cargados")
		return
	}

	if err := h.queries.DeleteClient(r.Context(), int32(id)); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
