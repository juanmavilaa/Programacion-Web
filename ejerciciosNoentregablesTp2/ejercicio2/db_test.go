package main

import (
	"context"
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"

	db "ejercicio2/db/sqlc"
)

func TestQueries_CRUD(t *testing.T) {
	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=tp2_test sslmode=disable"

	database, err := sql.Open("pgx", connStr)
	if err != nil {
		t.Fatalf("Error al abrir la base de datos: %v", err)
	}
	defer database.Close()

	if err := database.Ping(); err != nil {
		t.Fatalf("Error al conectar con la base de datos: %v", err)
	}

	queries := db.New(database)
	ctx := context.Background()

	var createdUser db.User

	t.Run("CreateUser", func(t *testing.T) {
		createdUser, err = queries.CreateUser(
			ctx,
			db.CreateUserParams{
				Name:  "Usuario Test",
				Email: "usuario.sqlc.test@example.com",
			},
		)

		if err != nil {
			t.Fatalf("Error al crear usuario: %v", err)
		}

		if createdUser.ID == 0 {
			t.Errorf("Se esperaba que el usuario tuviera un ID")
		}
	})

	t.Run("GetUser", func(t *testing.T) {
		user, err := queries.GetUser(ctx, createdUser.ID)

		if err != nil {
			t.Fatalf("Error al obtener usuario: %v", err)
		}

		if user.Name != createdUser.Name {
			t.Errorf(
				"Nombre obtenido = %s; esperado = %s",
				user.Name,
				createdUser.Name,
			)
		}

		if user.Email != createdUser.Email {
			t.Errorf(
				"Email obtenido = %s; esperado = %s",
				user.Email,
				createdUser.Email,
			)
		}
	})

	t.Run("UpdateUser", func(t *testing.T) {
		err := queries.UpdateUser(
			ctx,
			db.UpdateUserParams{
				ID:    createdUser.ID,
				Name:  "Usuario Test Actualizado",
				Email: createdUser.Email,
			},
		)

		if err != nil {
			t.Fatalf("Error al actualizar usuario: %v", err)
		}
	})

	t.Run("ListUsers", func(t *testing.T) {
		users, err := queries.ListUsers(ctx)

		if err != nil {
			t.Fatalf("Error al listar usuarios: %v", err)
		}

		found := false

		for _, user := range users {
			if user.ID == createdUser.ID {
				found = true

				if user.Name != "Usuario Test Actualizado" {
					t.Errorf(
						"Nombre obtenido = %s; esperado = Usuario Test Actualizado",
						user.Name,
					)
				}
			}
		}

		if !found {
			t.Errorf("El usuario creado no se encontró en la lista")
		}
	})

	t.Run("DeleteUser", func(t *testing.T) {
		err := queries.DeleteUser(ctx, createdUser.ID)

		if err != nil {
			t.Fatalf("Error al eliminar usuario: %v", err)
		}

		_, err = queries.GetUser(ctx, createdUser.ID)

		if err != sql.ErrNoRows {
			t.Errorf(
				"Se esperaba sql.ErrNoRows; se obtuvo %v",
				err,
			)
		}
	})
}
