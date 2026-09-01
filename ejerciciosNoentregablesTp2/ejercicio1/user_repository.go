package ejercicio1

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type User struct {
	ID        int
	Name      string
	Email     string
	CreatedAt time.Time
}

type UserRepository struct {
	db *sql.DB
}

func (r *UserRepository) CreateUser(user *User) error {
	query := `
		INSERT INTO users (name, email)
		VALUES ($1, $2)
		RETURNING id, name, email, created_at
	`

	row := r.db.QueryRow(query, user.Name, user.Email)

	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) GetUserByID(id int) (*User, error) {
	query := `
		SELECT id, name, email, created_at
		FROM users
		WHERE id = $1
	`

	row := r.db.QueryRow(query, id)

	user := &User{}

	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt,
	)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (r *UserRepository) ListUsers() ([]*User, error) {
	query := `
		SELECT id, name, email, created_at
		FROM users
		ORDER BY name
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*User

	for rows.Next() {
		user := &User{}

		err := rows.Scan(
			&user.ID,
			&user.Name,
			&user.Email,
			&user.CreatedAt,
		)

		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) UpdateUser(user *User) error {
	query := `
		UPDATE users
		SET name = $2, email = $3
		WHERE id = $1
	`

	_, err := r.db.Exec(
		query,
		user.ID,
		user.Name,
		user.Email,
	)

	if err != nil {
		return err
	}

	return nil
}

func (r *UserRepository) DeleteUser(id int) error {
	query := `
		DELETE FROM users
		WHERE id = $1
	`

	_, err := r.db.Exec(query, id)

	if err != nil {
		return err
	}

	return nil
}
