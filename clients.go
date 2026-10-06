package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"regexp"

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
		return errors.New("el email no es valido")
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
		case http.MethodGet:
			h.list(w, r)
		case http.MethodPost:
			h.create(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	id, err := idFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, "id invalido", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.get(w, r, id)
	case http.MethodPut:
		h.update(w, r, id)
	case http.MethodDelete:
		h.delete(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (h *ClientHandler) list(w http.ResponseWriter, r *http.Request) {
	clients, err := h.queries.ListClients(r.Context())
	if err != nil {
		http.Error(w, "error al listar los clientes", http.StatusInternalServerError)
		return
	}

	// sin esto si no hay clientes devuelve null
	if clients == nil {
		clients = []db.Client{}
	}

	writeJSON(w, http.StatusOK, clients)
}

func (h *ClientHandler) create(w http.ResponseWriter, r *http.Request) {
	var params db.CreateClientParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "JSON invalido", http.StatusBadRequest)
		return
	}

	if err := validarCliente(params.Name, params.Email); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := h.queries.GetClientByEmail(r.Context(), params.Email)
	if err == nil {
		http.Error(w, "ya hay un cliente con ese email", http.StatusConflict)
		return
	}
	if !errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "error al crear el cliente", http.StatusInternalServerError)
		return
	}

	client, err := h.queries.CreateClient(r.Context(), params)
	if err != nil {
		http.Error(w, "error al crear el cliente", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, client)
}

func (h *ClientHandler) get(w http.ResponseWriter, r *http.Request, id int32) {
	client, err := h.queries.GetClient(r.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "cliente no encontrado", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "error al buscar el cliente", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, client)
}

func (h *ClientHandler) update(w http.ResponseWriter, r *http.Request, id int32) {
	var params db.UpdateClientParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "JSON invalido", http.StatusBadRequest)
		return
	}

	if err := validarCliente(params.Name, params.Email); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := h.queries.GetClient(r.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "cliente no encontrado", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "error al buscar el cliente", http.StatusInternalServerError)
		return
	}

	// el email puede ser el mismo que ya tenia, pero no el de otro cliente
	otro, err := h.queries.GetClientByEmail(r.Context(), params.Email)
	if err == nil && otro.ID != id {
		http.Error(w, "ya hay un cliente con ese email", http.StatusConflict)
		return
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "error al actualizar el cliente", http.StatusInternalServerError)
		return
	}

	params.ID = id
	if err := h.queries.UpdateClient(r.Context(), params); err != nil {
		http.Error(w, "error al actualizar el cliente", http.StatusInternalServerError)
		return
	}

	client, err := h.queries.GetClient(r.Context(), id)
	if err != nil {
		http.Error(w, "error al buscar el cliente", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, client)
}

func (h *ClientHandler) delete(w http.ResponseWriter, r *http.Request, id int32) {
	_, err := h.queries.GetClient(r.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "cliente no encontrado", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "error al buscar el cliente", http.StatusInternalServerError)
		return
	}

	// no dejamos borrar un cliente que tiene entrenamientos cargados
	workouts, err := h.queries.ListWorkoutsByClient(r.Context(), id)
	if err != nil {
		http.Error(w, "error al borrar el cliente", http.StatusInternalServerError)
		return
	}
	if len(workouts) > 0 {
		http.Error(w, "el cliente tiene entrenamientos cargados", http.StatusConflict)
		return
	}

	if err := h.queries.DeleteClient(r.Context(), id); err != nil {
		http.Error(w, "error al borrar el cliente", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
