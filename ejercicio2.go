package main

import (
	"fmt"
	"net/http"
)

func main() {

	http.HandleFunc("/contacto", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		if r.Method == "GET" {
			fmt.Fprint(w, `
<form method="get" action="/contacto">

<h1>Datos a cargar</h1>

<input type="text" name="nombre" placeholder="ponga su nombre">
<input type="email" name="email" placeholder="ponga su email">
<textarea name="mensaje" placeholder="ponga su mensaje"></textarea>
<input type="submit" value="Enviar">

</form>
`)
		}

		if r.Method == "POST" {
			nombre := r.FormValue("nombre")
			email := r.FormValue("email")
			mensaje := r.FormValue("mensaje")

			if nombre == "" || email == "" || mensaje == "" {
				fmt.Fprint(w, "hay un campo vacio")
				return
			}

			fmt.Fprintf(w, `datos recibidos:
nombre: %s
email: %s
mensaje: %s`, nombre, email, mensaje)
		}
	})

	http.ListenAndServe(":8080", nil)
}
