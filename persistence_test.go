package main

import (
	"context"
	"database/sql"
	"testing"

	db "reps/db/sqlc"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestPersistenceCRUD(t *testing.T) {
	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=reps sslmode=disable"

	database, err := sql.Open("pgx", connStr)
	if err != nil {
		t.Fatalf("error al abrir la base: %v", err)
	}
	defer database.Close()

	if err := database.Ping(); err != nil {
		t.Fatalf("error al conectar con la base: %v", err)
	}

	queries := db.New(database)
	ctx := context.Background()

	client, err := queries.CreateClient(ctx, db.CreateClientParams{
		Name:  "Juan Test",
		Email: "juan@test.com",
	})
	if err != nil {
		t.Fatalf("error creando cliente: %v", err)
	}

	if client.ID == 0 {
		t.Errorf("el cliente deberia tener un id")
	}

	savedClient, err := queries.GetClient(ctx, client.ID)
	if err != nil {
		t.Fatalf("error buscando cliente: %v", err)
	}

	if savedClient.Name != "Juan Test" {
		t.Errorf("el nombre del cliente no coincide")
	}

	clients, err := queries.ListClients(ctx)
	if err != nil {
		t.Fatalf("error listando clientes: %v", err)
	}

	if len(clients) == 0 {
		t.Errorf("la lista de clientes esta vacia")
	}

	err = queries.UpdateClient(ctx, db.UpdateClientParams{
		ID:    client.ID,
		Name:  "Juan Actualizado",
		Email: "juan@test.com",
	})
	if err != nil {
		t.Fatalf("error actualizando cliente: %v", err)
	}

	updatedClient, err := queries.GetClient(ctx, client.ID)
	if err != nil {
		t.Fatalf("error buscando cliente actualizado: %v", err)
	}

	if updatedClient.Name != "Juan Actualizado" {
		t.Errorf("el cliente no se actualizo")
	}

	workout, err := queries.CreateWorkout(ctx, db.CreateWorkoutParams{
		ClientID:    client.ID,
		Exercise:    "Press banca",
		Sets:        4,
		Repetitions: 8,
		Weight:      "80.00",
	})
	if err != nil {
		t.Fatalf("error creando entrenamiento: %v", err)
	}

	if workout.ID == 0 {
		t.Errorf("el entrenamiento deberia tener un id")
	}

	if !workout.CreatedAt.Valid {
		t.Errorf("created_at deberia generarse automaticamente")
	}

	savedWorkout, err := queries.GetWorkout(ctx, workout.ID)
	if err != nil {
		t.Fatalf("error buscando entrenamiento: %v", err)
	}

	if savedWorkout.Exercise != "Press banca" {
		t.Errorf("el ejercicio no coincide")
	}

	workouts, err := queries.ListWorkouts(ctx)
	if err != nil {
		t.Fatalf("error listando entrenamientos: %v", err)
	}

	if len(workouts) == 0 {
		t.Errorf("la lista de entrenamientos esta vacia")
	}

	clientWorkouts, err := queries.ListWorkoutsByClient(ctx, client.ID)
	if err != nil {
		t.Fatalf("error listando entrenamientos del cliente: %v", err)
	}

	if len(clientWorkouts) != 1 {
		t.Errorf("se esperaba un entrenamiento para el cliente")
	}

	err = queries.UpdateWorkout(ctx, db.UpdateWorkoutParams{
		ID:          workout.ID,
		ClientID:    client.ID,
		Exercise:    "Press banca inclinado",
		Sets:        4,
		Repetitions: 10,
		Weight:      "75.00",
	})
	if err != nil {
		t.Fatalf("error actualizando entrenamiento: %v", err)
	}

	updatedWorkout, err := queries.GetWorkout(ctx, workout.ID)
	if err != nil {
		t.Fatalf("error buscando entrenamiento actualizado: %v", err)
	}

	if updatedWorkout.Exercise != "Press banca inclinado" {
		t.Errorf("el entrenamiento no se actualizo")
	}

	err = queries.DeleteWorkout(ctx, workout.ID)
	if err != nil {
		t.Fatalf("error borrando entrenamiento: %v", err)
	}

	_, err = queries.GetWorkout(ctx, workout.ID)
	if err != sql.ErrNoRows {
		t.Errorf("el entrenamiento deberia estar borrado")
	}

	err = queries.DeleteClient(ctx, client.ID)
	if err != nil {
		t.Fatalf("error borrando cliente: %v", err)
	}

	_, err = queries.GetClient(ctx, client.ID)
	if err != sql.ErrNoRows {
		t.Errorf("el cliente deberia estar borrado")
	}
}
