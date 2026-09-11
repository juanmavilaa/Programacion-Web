# Reps - TP2

Reps es un trackeador de ejercicios de gimnasio que permite registrar clientes y sus entrenamientos.

En este TP se agregó la capa de persistencia utilizando PostgreSQL y sqlc.

## Requisitos

Para ejecutar el proyecto y los tests se necesita:

- Go
- Docker
- Docker Compose
- sqlc

Para instalar sqlc:

```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

## Ejecutar los tests

El TP2 se encuentra en la branch `tp2`.

Después de clonar el repositorio:

```bash
git switch tp2
```

Desde la raíz del proyecto:

```bash
./test.sh
```

Si el script no tiene permiso de ejecución:

```bash
chmod +x test.sh
./test.sh
```

El script se encarga de:

1. borrar contenedores y volúmenes anteriores;
2. ejecutar `sqlc generate`;
3. compilar el proyecto;
4. levantar PostgreSQL;
5. esperar hasta que la base de datos esté disponible;
6. crear las tablas definidas en `schema.sql`;
7. ejecutar los tests con `go test -v`;
8. eliminar el contenedor y el volumen al finalizar.

## Servidor web

Para ejecutar la página inicial:

```bash
go run main.go
```

Luego ingresar a:

```text
http://localhost:8080
```

## Persistencia

El esquema de la base de datos está en:

```text
db/schema/schema.sql
```

Las consultas SQL están en:

```text
db/queries/queries.sql
```

La configuración de sqlc está en:

```text
sqlc.yaml
```

El código generado por sqlc queda en:

```text
db/sqlc/
```

Los archivos de `db/sqlc/` son generados automáticamente y no se modifican manualmente.

La documentación de persistencia está en:

```text
docs/persistencia.md
```