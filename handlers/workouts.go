package handlers

import (
	"database/sql"
	"errors"
	"math"
	"net/http"
	"strconv"
	"strings"

	db "reps/db/sqlc"
	"reps/views"
)

type WorkoutHandler struct {
	queries *db.Queries
}

func NewWorkoutHandler(q *db.Queries) *WorkoutHandler {
	return &WorkoutHandler{queries: q}
}

// devuelve 0 si el texto no es un numero valido
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
	if r.URL.Path == "/workouts" {
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
		mostrarError(w, r, http.StatusUnprocessableEntity, err.Error())
		return
	}

	_, err := h.queries.GetClient(r.Context(), params.ClientID)
	if err == sql.ErrNoRows {
		mostrarError(w, r, http.StatusUnprocessableEntity, "no existe el cliente")
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

	if err := views.WorkoutCreated(workouts, clients).Render(r.Context(), w); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
	}
}

func (h *WorkoutHandler) Delete(w http.ResponseWriter, r *http.Request, id int32) {
	if err := h.queries.DeleteWorkout(r.Context(), id); err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
