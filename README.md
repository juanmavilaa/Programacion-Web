# Reps - TP5

Reps es un trackeador de ejercicios de gimnasio que permite registrar clientes y sus entrenamientos.

En este TP se reemplazó el frontend con JavaScript por una aplicación renderizada en el servidor usando templ. Los formularios hacen un POST y el servidor redirige a la página principal (Post/Redirect/Get).

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

También se necesita conexión a internet: Pico.css se carga desde un CDN y Go descarga templ (y su propia versión, si hace falta) la primera vez.

## Ejecutar la aplicación

El TP5 se encuentra en la branch `tp5`.

Después de clonar el repositorio:

```bash
git switch tp5
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

## Tests

Los tests de persistencia se ejecutan con:

```bash
make test
```

## Documentación

La documentación de las vistas y de los handlers está en:

```text
docs/vistas.md
```

La documentación de persistencia está en:

```text
docs/persistencia.md
```
