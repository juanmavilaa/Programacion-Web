package ejercicio1

import (
	"database/sql"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func TestUserRepository(t *testing.T) {
	connStr := "host=localhost port=5432 user=postgres password=postgres dbname=tp2_test sslmode=disable"

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		t.Fatalf("Error al abrir la base de datos: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		t.Fatalf("Error al conectar con la base de datos: %v", err)
	}

	db.SetMaxOpenConns(25)

	repository := &UserRepository{
		db: db,
	}

	user := &User{
		Name:  "Juan",
		Email: "juan@example.com",
	}

	t.Run("Crear", func(t *testing.T) {
		err := repository.CreateUser(user)

		if err != nil {
			t.Fatalf("Error al crear el usuario: %v", err)
		}

		if user.ID == 0 {
			t.Errorf("Se esperaba que el usuario tuviera un ID")
		}
	})

	t.Run("Leer", func(t *testing.T) {
		foundUser, err := repository.GetUserByID(user.ID)

		if err != nil {
			t.Fatalf("Error al obtener el usuario: %v", err)
		}

		if foundUser.Name != user.Name {
			t.Errorf(
				"Nombre obtenido = %s; se esperaba %s",
				foundUser.Name,
				user.Name,
			)
		}

		if foundUser.Email != user.Email {
			t.Errorf(
				"Email obtenido = %s; se esperaba %s",
				foundUser.Email,
				user.Email,
			)
		}
	})

	t.Run("Actualizar", func(t *testing.T) {
		user.Name = "Juan Actualizado"

		err := repository.UpdateUser(user)
		if err != nil {
			t.Fatalf("Error al actualizar el usuario: %v", err)
		}

		updatedUser, err := repository.GetUserByID(user.ID)
		if err != nil {
			t.Fatalf("Error al obtener el usuario actualizado: %v", err)
		}

		if updatedUser.Name != "Juan Actualizado" {
			t.Errorf(
				"Nombre obtenido = %s; se esperaba Juan Actualizado",
				updatedUser.Name,
			)
		}
	})

	t.Run("Listar", func(t *testing.T) {
		users, err := repository.ListUsers()

		if err != nil {
			t.Fatalf("Error al listar usuarios: %v", err)
		}

		found := false

		for _, listedUser := range users {
			if listedUser.ID == user.ID {
				found = true
			}
		}

		if !found {
			t.Errorf("El usuario creado no se encontró en la lista")
		}
	})

	t.Run("Eliminar", func(t *testing.T) {
		err := repository.DeleteUser(user.ID)

		if err != nil {
			t.Fatalf("Error al eliminar el usuario: %v", err)
		}

		_, err = repository.GetUserByID(user.ID)

		if err != sql.ErrNoRows {
			t.Errorf(
				"Se esperaba sql.ErrNoRows después de eliminar; se obtuvo %v",
				err,
			)
		}
	})
}
