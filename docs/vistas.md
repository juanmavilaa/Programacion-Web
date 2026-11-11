# Vistas

La aplicación se renderiza en el servidor con templ. No hay JavaScript ni API JSON: el navegador pide la página y el servidor la devuelve completa con los datos de la base.

## Componentes

Los componentes están en la carpeta `views/`:

- `layout.templ`: estructura base de la página. Incluye Pico.css y `static/styles.css`.
- `entity_form.templ`: formularios para crear un cliente y un entrenamiento.
- `entity_list.templ`: tablas con los clientes y los entrenamientos.
- `index.templ`: página principal (`IndexPage`), que usa el layout y junta los formularios y las listas.

Los archivos `.templ` se compilan a código Go con `templ generate`.

## Handlers

Los handlers están en la carpeta `handlers/`:

- `index.go`: atiende `GET /`. Busca los clientes y los entrenamientos con sqlc y renderiza `IndexPage`.
- `clients.go`: atiende `POST /clients` y `POST /clients/delete`.
- `workouts.go`: atiende `POST /workouts` y `POST /workouts/delete`.

Los formularios HTML solo pueden hacer GET o POST, por eso el borrado también se hace con POST y el id se manda en un campo oculto.

Las validaciones son estas: el nombre no puede estar vacío, el email tiene que ser válido y no repetirse, el cliente de un entrenamiento tiene que existir, las series y repeticiones tienen que ser mayores a 0 y el peso tiene que ser un número. Un cliente que tiene entrenamientos no se puede borrar.

## Post/Redirect/Get

Después de guardar o borrar, el handler redirige a `/` con `http.Redirect`. El navegador pide la página de nuevo y se muestra la lista actualizada. Si se recarga la página no se vuelve a enviar el formulario.

Si algo no es válido se redirige a `/?error=mensaje` y la página muestra el mensaje arriba de todo.

## Generación del código

`make generate` ejecuta `sqlc generate` y `templ generate`.

templ se ejecuta con `go run` en la misma versión que figura en `go.mod` (v0.3.1020). Hace falta porque el código que genera cada versión de templ solo funciona con el runtime de esa misma versión. Si se cambia la versión en `go.mod`, hay que cambiarla también en el `Makefile`.
