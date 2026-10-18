# Persistencia

Reps utiliza PostgreSQL para almacenar clientes y sus entrenamientos.

## Tablas

La base de datos tiene dos tablas: `clients` y `workouts`.

### clients

Guarda los datos de cada cliente:

- `id`
- `name`
- `email`

### workouts

Guarda los entrenamientos:

- `id`
- `client_id`
- `exercise`
- `sets`
- `repetitions`
- `weight`
- `created_at`

Cada entrenamiento pertenece a un cliente mediante `client_id`.

El campo `created_at` guarda automáticamente la fecha y hora en la que se creó el registro. En este TP se le agregó `NOT NULL` para que sqlc lo genere como `time.Time` y en el JSON se vea la fecha directamente.

## Acceso a datos

El esquema está definido en:

```text
db/schema/schema.sql
```

Las consultas SQL están definidas en:

```text
db/queries/queries.sql
```

Además del CRUD de cada tabla se agregó `GetClientByEmail`, que la API usa para no permitir emails repetidos.

Se utiliza sqlc para generar el código Go de acceso a datos.

La configuración de sqlc está en:

```text
sqlc.yaml
```

El código generado queda en:

```text
db/sqlc/
```

## Creación de las tablas

En `docker-compose.yml` la carpeta `db/schema/` se monta en `/docker-entrypoint-initdb.d` del contenedor de PostgreSQL. Así las tablas se crean solas la primera vez que arranca la base.

## Tests

Los tests de persistencia (`persistence_test.go`) utilizan PostgreSQL levantado con Docker.

Para ejecutar todas las tareas necesarias:

```bash
make test
```

El comando genera el código con sqlc, compila, levanta la base, corre los tests de persistencia y al final borra los contenedores y el volumen.
