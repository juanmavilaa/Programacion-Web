package handlers

import (
	"net/http"

	"reps/views"
)

// responde con el mensaje de error sin tocar la lista que ya esta en pantalla
func mostrarError(w http.ResponseWriter, r *http.Request, status int, mensaje string) {
	w.Header().Set("HX-Reswap", "none")
	w.WriteHeader(status)
	views.ErrorMessage(mensaje).Render(r.Context(), w)
}
