package handlers

import (
	"database/sql"
	"encoding/json"
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
	switch r.URL.Path {
	case "/clients":
		h.clientsHandler(w, r)
	default:
		h.clientHandler(w, r)
	}
}

// Manejador para /clients
func (h *ClientHandler) clientsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getClients(w, r)
	case http.MethodPost:
		h.createClient(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Manejador para /clients/{id}
func (h *ClientHandler) clientHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(parts[2])
	if err != nil {
		http.Error(w, "id de cliente invalido", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getClient(w, r, int32(id))
	case http.MethodPut:
		h.updateClient(w, r, int32(id))
	case http.MethodDelete:
		h.deleteClient(w, r, int32(id))
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GET /clients - Listar todos los clientes
func (h *ClientHandler) getClients(w http.ResponseWriter, r *http.Request) {
	clients, err := h.queries.ListClients(r.Context())
	if err != nil {
		http.Error(w, "error al listar los clientes", http.StatusInternalServerError)
		return
	}

	// si no hay clientes sqlc devuelve nil y el json queda null
	if clients == nil {
		clients = []db.Client{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clients)
}

// POST /clients - Crear nuevo cliente
func (h *ClientHandler) createClient(w http.ResponseWriter, r *http.Request) {
	var params db.CreateClientParams
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validarCliente(params.Name, params.Email); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = h.queries.GetClientByEmail(r.Context(), params.Email)
	if err == nil {
		http.Error(w, "ya hay un cliente con ese email", http.StatusConflict)
		return
	}
	if err != sql.ErrNoRows {
		http.Error(w, "error al crear el cliente", http.StatusInternalServerError)
		return
	}

	client, err := h.queries.CreateClient(r.Context(), params)
	if err != nil {
		http.Error(w, "error al crear el cliente", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(client)
}

// GET /clients/{id} - Obtener cliente especifico
func (h *ClientHandler) getClient(w http.ResponseWriter, r *http.Request, id int32) {
	client, err := h.queries.GetClient(r.Context(), id)
	if err == sql.ErrNoRows {
		http.Error(w, "cliente no encontrado", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "error al buscar el cliente", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(client)
}

// PUT /clients/{id} - Actualizar cliente
func (h *ClientHandler) updateClient(w http.ResponseWriter, r *http.Request, id int32) {
	var params db.UpdateClientParams
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validarCliente(params.Name, params.Email); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = h.queries.GetClient(r.Context(), id)
	if err == sql.ErrNoRows {
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
	if err != nil && err != sql.ErrNoRows {
		http.Error(w, "error al actualizar el cliente", http.StatusInternalServerError)
		return
	}

	// Mantener el ID original
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(client)
}

// DELETE /clients/{id} - Eliminar cliente
func (h *ClientHandler) deleteClient(w http.ResponseWriter, r *http.Request, id int32) {
	_, err := h.queries.GetClient(r.Context(), id)
	if err == sql.ErrNoRows {
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
