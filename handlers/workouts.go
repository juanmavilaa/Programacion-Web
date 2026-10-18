package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	db "reps/db/sqlc"
)

type WorkoutHandler struct {
	queries *db.Queries
}

func NewWorkoutHandler(q *db.Queries) *WorkoutHandler {
	return &WorkoutHandler{queries: q}
}

func validarWorkout(clientID int32, exercise string, sets, repetitions int32, weight string) error {
	if clientID <= 0 {
		return errors.New("falta el client_id")
	}
	if exercise == "" {
		return errors.New("falta el ejercicio")
	}
	if sets <= 0 || repetitions <= 0 {
		return errors.New("las series y repeticiones tienen que ser mayores a 0")
	}

	// weight viene como string porque en la base es NUMERIC(6,2)
	peso, err := strconv.ParseFloat(weight, 64)
	if err != nil || peso < 0 || peso > 9999.99 {
		return errors.New("el peso no es valido")
	}

	return nil
}

func (h *WorkoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/workouts":
		h.workoutsHandler(w, r)
	default:
		h.workoutHandler(w, r)
	}
}

// Manejador para /workouts
func (h *WorkoutHandler) workoutsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.getWorkouts(w, r)
	case http.MethodPost:
		h.createWorkout(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Manejador para /workouts/{id}
func (h *WorkoutHandler) workoutHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(parts[2])
	if err != nil {
		http.Error(w, "id de entrenamiento invalido", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		h.getWorkout(w, r, int32(id))
	case http.MethodPut:
		h.updateWorkout(w, r, int32(id))
	case http.MethodDelete:
		h.deleteWorkout(w, r, int32(id))
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GET /workouts - Listar todos los entrenamientos
func (h *WorkoutHandler) getWorkouts(w http.ResponseWriter, r *http.Request) {
	workouts, err := h.queries.ListWorkouts(r.Context())
	if err != nil {
		http.Error(w, "error al listar los entrenamientos", http.StatusInternalServerError)
		return
	}

	if workouts == nil {
		workouts = []db.Workout{}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workouts)
}

// POST /workouts - Crear nuevo entrenamiento
func (h *WorkoutHandler) createWorkout(w http.ResponseWriter, r *http.Request) {
	var params db.CreateWorkoutParams
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validarWorkout(params.ClientID, params.Exercise, params.Sets, params.Repetitions, params.Weight); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = h.queries.GetClient(r.Context(), params.ClientID)
	if err == sql.ErrNoRows {
		http.Error(w, "no existe el cliente", http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, "error al crear el entrenamiento", http.StatusInternalServerError)
		return
	}

	workout, err := h.queries.CreateWorkout(r.Context(), params)
	if err != nil {
		http.Error(w, "error al crear el entrenamiento", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(workout)
}

// GET /workouts/{id} - Obtener entrenamiento especifico
func (h *WorkoutHandler) getWorkout(w http.ResponseWriter, r *http.Request, id int32) {
	workout, err := h.queries.GetWorkout(r.Context(), id)
	if err == sql.ErrNoRows {
		http.Error(w, "entrenamiento no encontrado", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "error al buscar el entrenamiento", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workout)
}

// PUT /workouts/{id} - Actualizar entrenamiento
func (h *WorkoutHandler) updateWorkout(w http.ResponseWriter, r *http.Request, id int32) {
	var params db.UpdateWorkoutParams
	err := json.NewDecoder(r.Body).Decode(&params)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := validarWorkout(params.ClientID, params.Exercise, params.Sets, params.Repetitions, params.Weight); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = h.queries.GetWorkout(r.Context(), id)
	if err == sql.ErrNoRows {
		http.Error(w, "entrenamiento no encontrado", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "error al buscar el entrenamiento", http.StatusInternalServerError)
		return
	}

	_, err = h.queries.GetClient(r.Context(), params.ClientID)
	if err == sql.ErrNoRows {
		http.Error(w, "no existe el cliente", http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, "error al actualizar el entrenamiento", http.StatusInternalServerError)
		return
	}

	// Mantener el ID original
	params.ID = id
	if err := h.queries.UpdateWorkout(r.Context(), params); err != nil {
		http.Error(w, "error al actualizar el entrenamiento", http.StatusInternalServerError)
		return
	}

	workout, err := h.queries.GetWorkout(r.Context(), id)
	if err != nil {
		http.Error(w, "error al buscar el entrenamiento", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(workout)
}

// DELETE /workouts/{id} - Eliminar entrenamiento
func (h *WorkoutHandler) deleteWorkout(w http.ResponseWriter, r *http.Request, id int32) {
	_, err := h.queries.GetWorkout(r.Context(), id)
	if err == sql.ErrNoRows {
		http.Error(w, "entrenamiento no encontrado", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "error al buscar el entrenamiento", http.StatusInternalServerError)
		return
	}

	if err := h.queries.DeleteWorkout(r.Context(), id); err != nil {
		http.Error(w, "error al borrar el entrenamiento", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
