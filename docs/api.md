# API

La API permite crear, listar, ver, actualizar y borrar clientes y entrenamientos. Las respuestas se devuelven en JSON. Cuando hay un error se devuelve un mensaje de texto con el código correspondiente.

## Clientes

Ejemplo de cliente:

```json
{
  "id": 1,
  "name": "Juan Perez",
  "email": "juan@reps.com"
}
```

Validaciones:

- el nombre no puede estar vacío;
- el email tiene que tener un formato válido;
- no puede haber dos clientes con el mismo email.

| Método | Ruta | Respuesta | Errores |
|--------|------|-----------|---------|
| POST | `/clients` | 201 con el cliente creado | 400, 409 |
| GET | `/clients` | 200 con la lista de clientes | |
| GET | `/clients/{id}` | 200 con el cliente | 400, 404 |
| PUT | `/clients/{id}` | 200 con el cliente actualizado | 400, 404, 409 |
| DELETE | `/clients/{id}` | 204 | 400, 404, 409 |

No se puede borrar un cliente que tenga entrenamientos cargados (devuelve 409). Primero hay que borrar sus entrenamientos.

## Entrenamientos

Ejemplo de entrenamiento:

```json
{
  "id": 1,
  "client_id": 1,
  "exercise": "Press banca",
  "sets": 4,
  "repetitions": 8,
  "weight": "80.00",
  "created_at": "2026-10-05T17:17:31.123456Z"
}
```

El peso se manda como texto (por ejemplo `"80.00"`) porque en la base es `NUMERIC(6,2)` y sqlc lo genera como `string`. La fecha `created_at` la pone la base al crear el registro.

Validaciones:

- el cliente tiene que existir;
- el ejercicio no puede estar vacío;
- las series y las repeticiones tienen que ser mayores a 0;
- el peso tiene que ser un número entre 0 y 9999.99.

| Método | Ruta | Respuesta | Errores |
|--------|------|-----------|---------|
| POST | `/workouts` | 201 con el entrenamiento creado | 400 |
| GET | `/workouts` | 200 con la lista de entrenamientos | |
| GET | `/workouts/{id}` | 200 con el entrenamiento | 400, 404 |
| PUT | `/workouts/{id}` | 200 con el entrenamiento actualizado | 400, 404 |
| DELETE | `/workouts/{id}` | 204 | 400, 404 |

## Códigos de error

- 400: el JSON está mal, falta algún dato, el id no es un número o el cliente no existe.
- 404: no existe el registro.
- 405: el método no está permitido en esa ruta.
- 409: el email ya está registrado o el cliente tiene entrenamientos.
- 500: error de la base de datos.

## Cómo está hecho

En `main.go` se abre la conexión con la base usando la variable de entorno `DATABASE_URL` y se crea el repositorio de sqlc con `db.New`. Ese repositorio se le pasa a `ClientHandler` y a `WorkoutHandler` (en `clients.go` y `workouts.go`).

Cada handler revisa la ruta: si es `/clients` atiende GET y POST, y si es `/clients/{id}` saca el id de la URL y atiende GET, PUT y DELETE. Lo mismo para `/workouts`.

Para los controles que dependen de la base (email repetido, cliente que no existe, cliente con entrenamientos) se usan las consultas `GetClientByEmail`, `GetClient` y `ListWorkoutsByClient` antes de guardar o borrar.

Además, un middleware escribe en el log cada request que llega.
