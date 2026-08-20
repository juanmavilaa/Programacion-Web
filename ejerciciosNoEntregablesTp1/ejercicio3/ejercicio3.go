package main

import (
	"fmt"
	"net/http"
)

func main() {

	archivos := http.FileServer(http.Dir("./static"))

	http.Handle("/", archivos)

	fmt.Println("Servidor funcionando en http://localhost:8080")

	http.ListenAndServe(":8080", nil)
}
