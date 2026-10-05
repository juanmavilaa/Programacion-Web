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
		return errors.New("formato de email inválido")
	}
	return nil
}

func validateClient(name, email string) error {
	if name == "" {
		return errors.New("el nombre es obligatorio")
	}
	return ValidateEmail(email)
}

// ServeHTTP atiende /clients y /clients/{id}
func (h *ClientHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// /clients
	if r.URL.Path == "/clients" {
		switch r.Method {
		case http.MethodGet:
			h.listClients(w, r)
		case http.MethodPost:
			h.createClient(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	// /clients/{id}
	id, err := idFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, "ID de cliente inválido", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getClient(w, r, id)
	case http.MethodPut:
		h.updateClient(w, r, id)
	case http.MethodDelete:
		h.deleteClient(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GET /clients - Listar todos los clientes
func (h *ClientHandler) listClients(w http.ResponseWriter, r *http.Request) {
	clients, err := h.queries.ListClients(r.Context())
	if err != nil {
		http.Error(w, "error al listar clientes", http.StatusInternalServerError)
		return
	}

	// Si no hay registros se devuelve [] en lugar de null
	if clients == nil {
		clients = []db.Client{}
	}

	writeJSON(w, http.StatusOK, clients)
}

// POST /clients - Crear un cliente
func (h *ClientHandler) createClient(w http.ResponseWriter, r *http.Request) {
	var params db.CreateClientParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if err := validateClient(params.Name, params.Email); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// El email no puede estar repetido
	_, err := h.queries.GetClientByEmail(r.Context(), params.Email)
	if err == nil {
		http.Error(w, "ya existe un cliente con ese email", http.StatusConflict)
		return
	}
	if !errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "error al buscar el email", http.StatusInternalServerError)
		return
	}

	client, err := h.queries.CreateClient(r.Context(), params)
	if err != nil {
		http.Error(w, "error al crear el cliente", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, client)
}

// GET /clients/{id} - Obtener un cliente
func (h *ClientHandler) getClient(w http.ResponseWriter, r *http.Request, id int32) {
	client, err := h.queries.GetClient(r.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "cliente no encontrado", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "error al obtener el cliente", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, client)
}

// PUT /clients/{id} - Actualizar un cliente
func (h *ClientHandler) updateClient(w http.ResponseWriter, r *http.Request, id int32) {
	var params db.UpdateClientParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if err := validateClient(params.Name, params.Email); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := h.queries.GetClient(r.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "cliente no encontrado", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "error al obtener el cliente", http.StatusInternalServerError)
		return
	}

	// El email no puede pertenecer a otro cliente
	other, err := h.queries.GetClientByEmail(r.Context(), params.Email)
	if err == nil && other.ID != id {
		http.Error(w, "ya existe un cliente con ese email", http.StatusConflict)
		return
	}
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "error al buscar el email", http.StatusInternalServerError)
		return
	}

	// Mantener el ID de la URL
	params.ID = id

	if err := h.queries.UpdateClient(r.Context(), params); err != nil {
		http.Error(w, "error al actualizar el cliente", http.StatusInternalServerError)
		return
	}

	client, err := h.queries.GetClient(r.Context(), id)
	if err != nil {
		http.Error(w, "error al obtener el cliente", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, client)
}

// DELETE /clients/{id} - Eliminar un cliente
func (h *ClientHandler) deleteClient(w http.ResponseWriter, r *http.Request, id int32) {
	_, err := h.queries.GetClient(r.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "cliente no encontrado", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "error al obtener el cliente", http.StatusInternalServerError)
		return
	}

	// No se puede borrar un cliente que tiene entrenamientos
	workouts, err := h.queries.ListWorkoutsByClient(r.Context(), id)
	if err != nil {
		http.Error(w, "error al buscar los entrenamientos", http.StatusInternalServerError)
		return
	}
	if len(workouts) > 0 {
		http.Error(w, "el cliente tiene entrenamientos registrados", http.StatusConflict)
		return
	}

	if err := h.queries.DeleteClient(r.Context(), id); err != nil {
		http.Error(w, "error al eliminar el cliente", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
