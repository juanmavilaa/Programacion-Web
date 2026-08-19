package main

import (
	"fmt"
	"net/http"
)

func main() {
	// 1. Define el contenido HTML
	htmlContent := `<!DOCTYPE html><html><head>
     <title>Página</title></head><body>
     <h1>Hola!</h1></body></html>`

	// 2. Registra un manejador (handler) para la ruta raíz "/"
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// 3. Establece la cabecera Content-Type
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		// 4. Escribe el HTML en la respuesta
		fmt.Fprint(w, htmlContent)
	})
	// 5. Define el puerto y muestra un mensaje
	port := ":8080"
	fmt.Printf("Servidor escuchando en http://localhost%s\n", port)

	// 6. Inicia el servidor HTTP
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Printf("Error al iniciar el servidor: %s\n", err)
	}
}
