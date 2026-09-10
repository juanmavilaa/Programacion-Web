package main

import (
	"context"
	"database/sql"
	"os"
	"testing"

	db "reps/db/sqlc"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func getTestDB(t *testing.T) *db.Queries {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "host=localhost port=5432 user=postgres password=postgres dbname=reps sslmode=disable"
	}

	conn, err := sql.Open("pgx", connStr)
	if err != nil {
		t.Fatalf("Error al conectar con la base de datos de test: %v", err)
	}

	if err := conn.Ping(); err != nil {
		t.Fatalf("No se pudo hacer ping a la base de datos de test: %v", err)
	}

	return db.New(conn)
}

func TestClientsAndWorkouts(t *testing.T) {
	queries := getTestDB(t)
	ctx := context.Background()

	var clientID int32

	t.Run("CreateClient", func(t *testing.T) {
		client, err := queries.CreateClient(ctx, db.CreateClientParams{
			Name:  "Test Client",
			Email: "test.client@example.com",
		})
		if err != nil {
			t.Fatalf("Error al crear cliente: %v", err)
		}
		if client.ID == 0 {
			t.Errorf("Se esperaba que el cliente tuviera un ID asignado")
		}
		clientID = client.ID
		t.Logf("✔ Cliente creado con ID: %d", clientID)
	})

	t.Run("GetClient", func(t *testing.T) {
		client, err := queries.GetClient(ctx, clientID)
		if err != nil {
			t.Fatalf("Error al obtener cliente: %v", err)
		}
		if client.Name != "Test Client" {
			t.Errorf("Nombre esperado 'Test Client', obtenido '%s'", client.Name)
		}
		t.Logf("✔ Cliente obtenido: %s (%s)", client.Name, client.Email)
	})

	t.Run("ListClients", func(t *testing.T) {
		clients, err := queries.ListClients(ctx)
		if err != nil {
			t.Fatalf("Error al listar clientes: %v", err)
		}
		if len(clients) == 0 {
			t.Errorf("Se esperaba al menos un cliente en la lista")
		}
		t.Logf("✔ Lista de clientes obtenida: %d clientes en total", len(clients))
	})

	t.Run("UpdateClient", func(t *testing.T) {
		err := queries.UpdateClient(ctx, db.UpdateClientParams{
			ID:    clientID,
			Name:  "Test Client Modificado",
			Email: "test.client.modificado@example.com",
		})
		if err != nil {
			t.Fatalf("Error al actualizar cliente: %v", err)
		}
		t.Logf("✔ Cliente ID %d actualizado", clientID)
	})

	t.Run("CreateWorkout", func(t *testing.T) {
		w, err := queries.CreateWorkout(ctx, db.CreateWorkoutParams{
			ClientID:    clientID,
			Exercise:    "Sentadilla con barra",
			Sets:        4,
			Repetitions: 8,
			Weight:      "100.00",
		})
		if err != nil {
			t.Fatalf("Error al crear workout: %v", err)
		}
		if w.ID == 0 {
			t.Errorf("Se esperaba que el workout tuviera un ID asignado")
		}
		t.Logf("✔ Workout creado: %s con %s kg para cliente %d", w.Exercise, w.Weight, clientID)
	})

	t.Run("ListWorkoutsByClient", func(t *testing.T) {
		workouts, err := queries.ListWorkoutsByClient(ctx, clientID)
		if err != nil {
			t.Fatalf("Error al listar workouts: %v", err)
		}
		if len(workouts) == 0 {
			t.Errorf("Se esperaba al menos un workout para el cliente %d", clientID)
		}
		t.Logf("✔ Workouts obtenidos para cliente %d: %d entrenamientos", clientID, len(workouts))
	})

	t.Run("DeleteClient_Cascada", func(t *testing.T) {
		err := queries.DeleteClient(ctx, clientID)
		if err != nil {
			t.Fatalf("Error al eliminar cliente: %v", err)
		}

		// Comprobar que no existe el cliente
		_, err = queries.GetClient(ctx, clientID)
		if err == nil {
			t.Errorf("El cliente debería haber sido eliminado")
		}

		// Comprobar que sus entrenamientos fueron eliminados en cascada
		workouts, err := queries.ListWorkoutsByClient(ctx, clientID)
		if err != nil {
			t.Fatalf("Error al consultar workouts: %v", err)
		}
		if len(workouts) != 0 {
			t.Errorf("Los workouts del cliente deberían haberse borrado en cascada")
		}
		t.Logf("✔ Cliente ID %d y sus entrenamientos borrados en cascada", clientID)
	})
}
