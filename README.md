# Reps - TP3

Reps es un trackeador de ejercicios de gimnasio que permite registrar clientes y sus entrenamientos.

## Requisitos

El proyecto está pensado para correr en Linux. Se necesita:

- Go
- Docker
- Docker Compose
- sqlc
- Hurl
- make
- curl

Para instalar sqlc:

```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

Para instalar Hurl en Ubuntu:

```bash
VERSION=6.1.1
curl --location --remote-name https://github.com/Orange-OpenSource/hurl/releases/download/$VERSION/hurl_${VERSION}_amd64.deb
sudo apt update && sudo apt install ./hurl_${VERSION}_amd64.deb
```

## Probar la API

Después de clonar el repositorio:

```bash
git switch tp3
```

Desde la raíz del proyecto:

```bash
make test-api
```

Este comando hace todo lo necesario para probar la API:

1. ejecuta `sqlc generate`;
2. compila el proyecto;
3. borra contenedores y volúmenes anteriores;
4. levanta PostgreSQL y espera a que la base esté lista;
5. levanta la aplicación y espera a que responda;
6. prueba cada uno de los endpoints con Hurl (`requests.hurl`);
7. borra los contenedores y el volumen al finalizar.

Si alguna prueba falla, igual se borran los contenedores y el comando termina con error.

Los tests de persistencia del TP2 se siguen pudiendo correr con:

```bash
make test
```

## Levantar el servidor

```bash
make api
```

Levanta la base y la aplicación en segundo plano. La API queda en `http://localhost:8080` y la página inicial del TP1 sigue estando en la raíz.

Para verlo con los logs en la terminal se puede usar `make run` (se corta con Ctrl+C).

Para bajar todo:

```bash
make down
```

Las tablas se crean solas la primera vez que arranca la base, a partir de `db/schema/schema.sql`.

## Probar a mano

Con el servidor recién levantado:

```bash
hurl --test requests.hurl
```

Las pruebas usan ids fijos (el primer cliente y el primer entrenamiento tienen id 1), por eso hay que correrlas con la base vacía. Si se quieren correr de nuevo, hacer `make down` y `make api` antes.

También se puede probar con curl, por ejemplo:

```bash
curl -i -X POST http://localhost:8080/clients \
  -H "Content-Type: application/json" \
  -d '{"name": "Juan Perez", "email": "juan@reps.com"}'

curl http://localhost:8080/clients

curl -i -X POST http://localhost:8080/workouts \
  -H "Content-Type: application/json" \
  -d '{"client_id": 1, "exercise": "Press banca", "sets": 4, "repetitions": 8, "weight": "80.00"}'

curl -i -X DELETE http://localhost:8080/workouts/1
```

## Endpoints

| Método | Ruta | Descripción |
|--------|------|-------------|
| POST | `/clients` | Crear un cliente |
| GET | `/clients` | Listar clientes |
| GET | `/clients/{id}` | Ver un cliente |
| PUT | `/clients/{id}` | Actualizar un cliente |
| DELETE | `/clients/{id}` | Borrar un cliente |
| POST | `/workouts` | Crear un entrenamiento |
| GET | `/workouts` | Listar entrenamientos |
| GET | `/workouts/{id}` | Ver un entrenamiento |
| PUT | `/workouts/{id}` | Actualizar un entrenamiento |
| DELETE | `/workouts/{id}` | Borrar un entrenamiento |

## Documentación

La documentación de la API está en:

```text
docs/api.md
docs/openapi.yaml
```

La documentación de persistencia está en:

```text
docs/persistencia.md
```
