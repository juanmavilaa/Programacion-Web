package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	db "reps/db/sqlc"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {

	conn, err := ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	queries := db.New(conn)

	clientHandler := NewClientHandler(queries)
	workoutHandler := NewWorkoutHandler(queries)

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
		return nil, fmt.Errorf("falta la variable DATABASE_URL")
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

// saca el id de urls como /clients/5
func idFromPath(path string) (int32, error) {
	parts := strings.Split(path, "/")
	if len(parts) != 3 {
		return 0, fmt.Errorf("url invalida")
	}

	id, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0, err
	}

	return int32(id), nil
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}
