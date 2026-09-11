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

El campo `created_at` guarda automáticamente la fecha y hora en la que se creó el registro.

## Acceso a datos

El esquema está definido en:

```text
db/schema/schema.sql
```

Las consultas SQL están definidas en:

```text
db/queries/queries.sql
```

Se utiliza sqlc para generar el código Go de acceso a datos.

La configuración de sqlc está en:

```text
sqlc.yaml
```

El código generado queda en:

```text
db/sqlc/
```

## Tests

Los tests utilizan PostgreSQL levantado con Docker.

Para ejecutar todas las tareas necesarias:

```bash
./test.sh
```

El script genera el código con sqlc, compila el proyecto, levanta una base PostgreSQL limpia, crea las tablas, ejecuta los tests y elimina el contenedor y el volumen al finalizar.