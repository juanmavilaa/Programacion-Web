package views

import (
	db "reps/db/sqlc"
	"strconv"
)

// busca el nombre del cliente de un entrenamiento
func clientName(clients []db.Client, id int32) string {
	for _, client := range clients {
		if client.ID == id {
			return client.Name
		}
	}
	return "Cliente " + strconv.Itoa(int(id))
}
