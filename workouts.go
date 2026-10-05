package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	db "reps/db/sqlc"
)

type WorkoutHandler struct {
	queries *db.Queries
}

func NewWorkoutHandler(q *db.Queries) *WorkoutHandler {
	return &WorkoutHandler{queries: q}
}

// El peso se guarda como NUMERIC(6,2): máximo 9999.99
const maxWeight = 10000

func validateWorkout(clientID int32, exercise string, sets, repetitions int32, weight string) error {
	if clientID <= 0 {
		return errors.New("el client_id es obligatorio")
	}
	if exercise == "" {
		return errors.New("el ejercicio es obligatorio")
	}
	if sets <= 0 {
		return errors.New("las series deben ser mayores a 0")
	}
	if repetitions <= 0 {
		return errors.New("las repeticiones deben ser mayores a 0")
	}

	value, err := strconv.ParseFloat(weight, 64)
	if err != nil || value < 0 || value >= maxWeight {
		return errors.New("el peso debe ser un número entre 0 y 9999.99")
	}

	return nil
}

// ServeHTTP atiende /workouts y /workouts/{id}
func (h *WorkoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// /workouts
	if r.URL.Path == "/workouts" {
		switch r.Method {
		case http.MethodGet:
			h.listWorkouts(w, r)
		case http.MethodPost:
			h.createWorkout(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		return
	}

	// /workouts/{id}
	id, err := idFromPath(r.URL.Path)
	if err != nil {
		http.Error(w, "ID de entrenamiento inválido", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getWorkout(w, r, id)
	case http.MethodPut:
		h.updateWorkout(w, r, id)
	case http.MethodDelete:
		h.deleteWorkout(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GET /workouts - Listar todos los entrenamientos
func (h *WorkoutHandler) listWorkouts(w http.ResponseWriter, r *http.Request) {
	workouts, err := h.queries.ListWorkouts(r.Context())
	if err != nil {
		http.Error(w, "error al listar entrenamientos", http.StatusInternalServerError)
		return
	}

	// Si no hay registros se devuelve [] en lugar de null
	if workouts == nil {
		workouts = []db.Workout{}
	}

	writeJSON(w, http.StatusOK, workouts)
}

// POST /workouts - Crear un entrenamiento
func (h *WorkoutHandler) createWorkout(w http.ResponseWriter, r *http.Request) {
	var params db.CreateWorkoutParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if err := validateWorkout(params.ClientID, params.Exercise, params.Sets, params.Repetitions, params.Weight); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// El cliente indicado tiene que existir
	_, err := h.queries.GetClient(r.Context(), params.ClientID)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "el cliente indicado no existe", http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, "error al obtener el cliente", http.StatusInternalServerError)
		return
	}

	workout, err := h.queries.CreateWorkout(r.Context(), params)
	if err != nil {
		http.Error(w, "error al crear el entrenamiento", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusCreated, workout)
}

// GET /workouts/{id} - Obtener un entrenamiento
func (h *WorkoutHandler) getWorkout(w http.ResponseWriter, r *http.Request, id int32) {
	workout, err := h.queries.GetWorkout(r.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "entrenamiento no encontrado", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "error al obtener el entrenamiento", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, workout)
}

// PUT /workouts/{id} - Actualizar un entrenamiento
func (h *WorkoutHandler) updateWorkout(w http.ResponseWriter, r *http.Request, id int32) {
	var params db.UpdateWorkoutParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	if err := validateWorkout(params.ClientID, params.Exercise, params.Sets, params.Repetitions, params.Weight); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := h.queries.GetWorkout(r.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "entrenamiento no encontrado", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "error al obtener el entrenamiento", http.StatusInternalServerError)
		return
	}

	// El cliente indicado tiene que existir
	_, err = h.queries.GetClient(r.Context(), params.ClientID)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "el cliente indicado no existe", http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, "error al obtener el cliente", http.StatusInternalServerError)
		return
	}

	// Mantener el ID de la URL
	params.ID = id

	if err := h.queries.UpdateWorkout(r.Context(), params); err != nil {
		http.Error(w, "error al actualizar el entrenamiento", http.StatusInternalServerError)
		return
	}

	workout, err := h.queries.GetWorkout(r.Context(), id)
	if err != nil {
		http.Error(w, "error al obtener el entrenamiento", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, workout)
}

// DELETE /workouts/{id} - Eliminar un entrenamiento
func (h *WorkoutHandler) deleteWorkout(w http.ResponseWriter, r *http.Request, id int32) {
	_, err := h.queries.GetWorkout(r.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "entrenamiento no encontrado", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "error al obtener el entrenamiento", http.StatusInternalServerError)
		return
	}

	if err := h.queries.DeleteWorkout(r.Context(), id); err != nil {
		http.Error(w, "error al eliminar el entrenamiento", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
