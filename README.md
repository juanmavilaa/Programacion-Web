# Reps - TP6

Reps es un trackeador de ejercicios de gimnasio que permite registrar clientes y sus entrenamientos.

La aplicación se renderiza en el servidor con templ. En este TP se agregó HTMX para que crear y borrar clientes y entrenamientos actualice solo la parte de la página que cambia, sin recargarla y sin escribir JavaScript.

## Requisitos

Para ejecutar el proyecto se necesita:

- Go 1.25 o superior
- Docker
- Docker Compose
- sqlc
- make

Para instalar sqlc:

```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

No hace falta instalar templ, el Makefile lo ejecuta con `go run`. La primera vez tarda un poco más porque lo descarga y lo compila.

También se necesita conexión a internet: Pico.css y HTMX se cargan desde un CDN y Go descarga templ (y su propia versión, si hace falta) la primera vez.

## Ejecutar la aplicación

El TP6 se encuentra en la branch `tp6`.

Después de clonar el repositorio:

```bash
git switch tp6
```

Desde la raíz del proyecto:

```bash
make run-app
```

El comando se encarga de:

1. borrar contenedores y volúmenes anteriores;
2. ejecutar `sqlc generate` y `templ generate`;
3. compilar el proyecto;
4. levantar PostgreSQL;
5. esperar hasta que la base de datos esté disponible;
6. ejecutar la aplicación.

Luego ingresar a:

```text
http://localhost:8080
```

Para detener la aplicación se usa Ctrl+C. Al cortarla se borran los contenedores y el volumen.

Las tablas se crean solas la primera vez que arranca la base, a partir de `db/schema/schema.sql`.

## Uso

En la página principal hay un formulario para agregar clientes y otro para agregar entrenamientos. Al enviarlos se actualiza la tabla correspondiente. Cada fila tiene un botón Borrar que pide confirmación y saca la fila. Si algo no es válido (email repetido, borrar un cliente que tiene entrenamientos, etc.) aparece un mensaje arriba de todo.

## Tests

Los tests de persistencia se ejecutan con:

```bash
make test
```

## Documentación

La documentación de las vistas, los handlers y el uso de HTMX está en:

```text
docs/vistas.md
```

La documentación de persistencia está en:

```text
docs/persistencia.md
```
