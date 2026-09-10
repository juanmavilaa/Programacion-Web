# Reps — Persistiendo el Dominio (TP2)

Trabajo Práctico 2 de Programación Web.

## Descripción del proyecto

Reps es una aplicación pensada para registrar entrenamientos de gimnasio. En este trabajo práctico se implementó la capa de persistencia utilizando PostgreSQL como base de datos relacional y sqlc para generar el código Go tipado a partir de las consultas SQL.

### Base de datos y persistencia

El modelo relacional se compone de dos tablas:

- **clients**: Guarda los datos del cliente (`id`, `name` y `email`).
- **workouts**: Registra los entrenamientos realizados (`id`, `client_id`, `exercise`, `sets`, `repetitions`, `weight` y `created_at`).

La tabla `workouts` está vinculada a `clients` mediante una clave foránea con `ON DELETE CASCADE`, de modo que al eliminar un cliente se eliminan automáticamente todos sus entrenamientos asociados.

## Requisitos

- Docker y Docker Compose
- make (opcional)

## Ejecución de tests

Para correr las pruebas del proyecto se puede clonar el repositorio, posicionarse en la rama `tp2` y ejecutar:

```bash
make test
```

O directamente ejecutando el script bash:

```bash
./test.sh
```

El script se encarga automáticamente de:
1. **Tareas previas:** Limpiar contenedores y volúmenes previos con `docker compose down -v`, verificar `sqlc`, levantar el servicio de PostgreSQL en Docker y esperar a que esté listo para recibir consultas, aplicando el esquema (`db/schema/schema.sql`) y las semillas (`db/seeds/seeds.sql`).
2. **Ejecución de tests:** Ejecutar la suite de pruebas completamente dentro de Docker mediante el servicio `web` definido en el `Dockerfile` (`docker compose run --rm web`).
3. **Tareas posteriores:** Destruir y borrar los contenedores temporales y sus volúmenes con `docker compose down -v` al finalizar (incluso si algún test falla).
