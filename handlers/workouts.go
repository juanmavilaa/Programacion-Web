package handlers

import (
	"database/sql"
	"errors"
	"math"
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

// devuelve 0 si el texto no es un numero o no entra en un int32 (lo rechaza validarWorkout)
func numero(s string) int32 {
	n, err := strconv.ParseInt(s, 10, 32)
	if err != nil {
		return 0
	}
	return int32(n)
}

func validarWorkout(clientID int32, exercise string, sets, repetitions int32, weight string) error {
	if clientID <= 0 {
		return errors.New("falta el cliente")
	}
	if exercise == "" {
		return errors.New("falta el ejercicio")
	}
	if sets <= 0 || repetitions <= 0 {
		return errors.New("las series y repeticiones tienen que ser mayores a 0")
	}

	// weight es un string porque en la base es NUMERIC(6,2)
	peso, err := strconv.ParseFloat(weight, 64)
	if err != nil || math.IsNaN(peso) || peso < 0 || peso > 9999.99 {
		return errors.New("el peso no es valido")
	}

	return nil
}

func (h *WorkoutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	switch r.URL.Path {
	case "/workouts":
		h.Create(w, r)
	case "/workouts/delete":
		h.Delete(w, r)
	default:
		http.NotFound(w, r)
	}
}

// POST /workouts - Crear nuevo entrenamiento
func (h *WorkoutHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "formulario invalido", http.StatusBadRequest)
		return
	}

	params := db.CreateWorkoutParams{
		ClientID:    numero(r.FormValue("client_id")),
		Exercise:    strings.TrimSpace(r.FormValue("exercise")),
		Sets:        numero(r.FormValue("sets")),
		Repetitions: numero(r.FormValue("repetitions")),
		Weight:      strings.TrimSpace(r.FormValue("weight")),
	}

	if err := validarWorkout(params.ClientID, params.Exercise, params.Sets, params.Repetitions, params.Weight); err != nil {
		redirectConError(w, r, err.Error())
		return
	}

	_, err := h.queries.GetClient(r.Context(), params.ClientID)
	if err == sql.ErrNoRows {
		redirectConError(w, r, "no existe el cliente")
		return
	}
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if _, err := h.queries.CreateWorkout(r.Context(), params); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// POST /workouts/delete - Eliminar entrenamiento
func (h *WorkoutHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "formulario invalido", http.StatusBadRequest)
		return
	}

	id, err := strconv.ParseInt(r.FormValue("id"), 10, 32)
	if err != nil {
		http.Error(w, "id invalido", http.StatusBadRequest)
		return
	}

	_, err = h.queries.GetWorkout(r.Context(), int32(id))
	if err == sql.ErrNoRows {
		redirectConError(w, r, "entrenamiento no encontrado")
		return
	}
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if err := h.queries.DeleteWorkout(r.Context(), int32(id)); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
