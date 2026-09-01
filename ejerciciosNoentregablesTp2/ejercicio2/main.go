package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"

	db "ejercicio2/db/sqlc"
)

func main() {
	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=tp2 sslmode=disable"

	database, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("Error al abrir la base de datos: %v", err)
	}
	defer database.Close()

	if err := database.Ping(); err != nil {
		log.Fatalf("Error al conectar con la base de datos: %v", err)
	}

	queries := db.New(database)
	ctx := context.Background()

	// CREATE
	createdUser, err := queries.CreateUser(
		ctx,
		db.CreateUserParams{
			Name:  "Juan SQLC",
			Email: "juan.sqlc@example.com",
		},
	)
	if err != nil {
		log.Fatalf("Error al crear usuario: %v", err)
	}

	fmt.Printf("Usuario creado: %+v\n", createdUser)

	// READ
	user, err := queries.GetUser(ctx, createdUser.ID)
	if err != nil {
		log.Fatalf("Error al obtener usuario: %v", err)
	}

	fmt.Printf("Usuario obtenido: %+v\n", user)

	// UPDATE
	err = queries.UpdateUser(
		ctx,
		db.UpdateUserParams{
			ID:    createdUser.ID,
			Name:  "Juan SQLC Actualizado",
			Email: "juan.sqlc@example.com",
		},
	)
	if err != nil {
		log.Fatalf("Error al actualizar usuario: %v", err)
	}

	updatedUser, err := queries.GetUser(ctx, createdUser.ID)
	if err != nil {
		log.Fatalf("Error al obtener usuario actualizado: %v", err)
	}

	fmt.Printf("Usuario actualizado: %+v\n", updatedUser)

	// LIST
	users, err := queries.ListUsers(ctx)
	if err != nil {
		log.Fatalf("Error al listar usuarios: %v", err)
	}

	fmt.Printf("Usuarios: %+v\n", users)

	// DELETE
	err = queries.DeleteUser(ctx, createdUser.ID)
	if err != nil {
		log.Fatalf("Error al eliminar usuario: %v", err)
	}

	fmt.Println("Usuario eliminado correctamente")

	// Comprobar que ya no existe
	_, err = queries.GetUser(ctx, createdUser.ID)

	if err == sql.ErrNoRows {
		fmt.Println("Usuario no encontrado después de eliminarlo")
	} else if err != nil {
		log.Fatalf("Error al verificar eliminación: %v", err)
	}
}
