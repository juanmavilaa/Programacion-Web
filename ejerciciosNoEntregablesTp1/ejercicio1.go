package main

import (
	"fmt"
	"net/http"
)

func main() {
	//defino el contenido html

	htmlContent := `<!DOCTYPE html>
	<html>
	<head><title>Hola Mundo</title></head>
	<body><h1>¡Servidor Funcionando!</h1></body>
	</html>`

	//Registra un manejador para la ruta raiz
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		//establece la cabecera content-type
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		//escribe el html en la respuesta
		fmt.Fprint(w, htmlContent)
	})

	//Registra un manejador para la ruta /about
	http.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		fmt.Fprint(w, "<h1>Acerca del servidor</h1><p>Servidor HTTP hecho en Go.</p>")
	})

	//Inicia el servidor en el puerto 8080
	port := ":8080"
	fmt.Printf("Servidor escuchando en http://localhost%s\n", port)
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Printf("Error al iniciar el servidor: %s\n", err)
	}
}
