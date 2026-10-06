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
	if r.URL.Path == "/workouts" {
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

func (h *WorkoutHandler) list(w http.ResponseWriter, r *http.Request) {
	workouts, err := h.queries.ListWorkouts(r.Context())
	if err != nil {
		http.Error(w, "error al listar los entrenamientos", http.StatusInternalServerError)
		return
	}

	if workouts == nil {
		workouts = []db.Workout{}
	}

	writeJSON(w, http.StatusOK, workouts)
}

func (h *WorkoutHandler) create(w http.ResponseWriter, r *http.Request) {
	var params db.CreateWorkoutParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "JSON invalido", http.StatusBadRequest)
		return
	}

	if err := validarWorkout(params.ClientID, params.Exercise, params.Sets, params.Repetitions, params.Weight); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := h.queries.GetClient(r.Context(), params.ClientID)
	if errors.Is(err, sql.ErrNoRows) {
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

	writeJSON(w, http.StatusCreated, workout)
}

func (h *WorkoutHandler) get(w http.ResponseWriter, r *http.Request, id int32) {
	workout, err := h.queries.GetWorkout(r.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "entrenamiento no encontrado", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "error al buscar el entrenamiento", http.StatusInternalServerError)
		return
	}

	writeJSON(w, http.StatusOK, workout)
}

func (h *WorkoutHandler) update(w http.ResponseWriter, r *http.Request, id int32) {
	var params db.UpdateWorkoutParams
	if err := json.NewDecoder(r.Body).Decode(&params); err != nil {
		http.Error(w, "JSON invalido", http.StatusBadRequest)
		return
	}

	if err := validarWorkout(params.ClientID, params.Exercise, params.Sets, params.Repetitions, params.Weight); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err := h.queries.GetWorkout(r.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "entrenamiento no encontrado", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, "error al buscar el entrenamiento", http.StatusInternalServerError)
		return
	}

	_, err = h.queries.GetClient(r.Context(), params.ClientID)
	if errors.Is(err, sql.ErrNoRows) {
		http.Error(w, "no existe el cliente", http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, "error al actualizar el entrenamiento", http.StatusInternalServerError)
		return
	}

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

	writeJSON(w, http.StatusOK, workout)
}

func (h *WorkoutHandler) delete(w http.ResponseWriter, r *http.Request, id int32) {
	_, err := h.queries.GetWorkout(r.Context(), id)
	if errors.Is(err, sql.ErrNoRows) {
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
