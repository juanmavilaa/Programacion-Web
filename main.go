package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	db "reps/db/sqlc"
	"reps/handlers"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {

	conn, err := ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	queries := db.New(conn)

	clientHandler := handlers.NewClientHandler(queries)
	workoutHandler := handlers.NewWorkoutHandler(queries)

	http.Handle("/clients", loggingMiddleware(clientHandler))
	http.Handle("/clients/", loggingMiddleware(clientHandler))
	http.Handle("/workouts", loggingMiddleware(workoutHandler))
	http.Handle("/workouts/", loggingMiddleware(workoutHandler))

	// pagina del tp1
	archivos := http.FileServer(http.Dir("./static"))
	http.Handle("/", archivos)

	log.Println("Servidor funcionando en http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func ConnectDB() (*sql.DB, error) {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		return nil, fmt.Errorf("DATABASE_URL is not set")
	}

	conn, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	if err := conn.Ping(); err != nil {
		conn.Close()
		return nil, err
	}

	return conn, nil
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Request: %s %s", r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}
