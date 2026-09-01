# Reps — Persistiendo el Dominio

Trabajo de cursada de **Programación Web**.

En esta etapa se define la estructura de persistencia de **Reps**, una aplicación para registrar entrenamientos de gimnasio.  
La entidad principal representa un registro de entrenamiento con:

- ejercicio;
- series;
- repeticiones;
- peso;
- fecha.

El acceso a datos se prepara utilizando **PostgreSQL como motor SQL** y **sqlc** para generar código Go tipado a partir del esquema y las consultas SQL.

> **Importante:** en esta entrega no es necesario levantar PostgreSQL, Docker ni un servidor web. El objetivo es definir el esquema, las consultas CRUD y generar el paquete Go de acceso a datos para usarlo en el siguiente práctico.

---

## Estructura del proyecto

```text
reps-persistencia/
├── go.mod
├── sqlc.yaml
└── db/
    ├── schema/
    │   └── schema.sql
    ├── queries/
    │   └── queries.sql
    └── sqlc/
        ├── db.go
        ├── models.go
        └── queries.sql.go
```

> Los nombres exactos de los archivos dentro de `db/sqlc/` pueden variar levemente según la versión de `sqlc`.

---

# Instalación desde una computadora sin herramientas

Para trabajar con esta entrega se necesita:

1. **Go**
2. **sqlc**
3. **Git**, únicamente si el proyecto se va a clonar desde un repositorio

No se necesita instalar PostgreSQL ni Docker para esta etapa.

---

## 1. Instalar Go

### Ubuntu / Linux

Abrir una terminal y ejecutar:

```bash
sudo apt update
sudo apt install golang-go
```

Comprobar la instalación:

```bash
go version
```

Debería mostrarse la versión instalada de Go.

### Windows

Abrir PowerShell y ejecutar:

```powershell
winget install GoLang.Go
```

Cerrar y volver a abrir PowerShell. Luego comprobar:

```powershell
go version
```

---

## 2. Obtener el proyecto

Si el proyecto fue entregado como archivo `.zip`, descomprimirlo y abrir una terminal dentro de la carpeta.

Si se obtiene desde GitHub, primero instalar Git.

### Ubuntu / Linux

```bash
sudo apt install git
```

### Windows

```powershell
winget install Git.Git
```

Luego clonar el repositorio correspondiente y entrar a la carpeta de esta entrega:

```bash
git clone <URL_DEL_REPOSITORIO>
cd <CARPETA_DEL_PROYECTO>
```

---

## 3. Instalar sqlc

Con Go ya instalado, ejecutar:

```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

Comprobar:

```bash
sqlc version
```

### Si Linux indica `sqlc: command not found`

Ejecutar:

```bash
export PATH="$PATH:$(go env GOPATH)/bin"
```

y volver a probar:

```bash
sqlc version
```

Para dejar esa configuración disponible en futuras terminales se puede agregar al archivo `~/.bashrc`:

```bash
echo 'export PATH="$PATH:$(go env GOPATH)/bin"' >> ~/.bashrc
source ~/.bashrc
```

### Si Windows no reconoce `sqlc`

Cerrar y volver a abrir PowerShell después de la instalación.

Si sigue sin reconocerlo, agregar temporalmente la carpeta de binarios de Go al `PATH` de la sesión:

```powershell
$env:Path += ";$(go env GOPATH)\bin"
```

Luego comprobar:

```powershell
sqlc version
```

---

# Generar la capa de acceso a datos

Ubicarse en la raíz del proyecto, es decir, en la carpeta donde se encuentran:

```text
go.mod
sqlc.yaml
db/
```

Ejecutar:

```bash
sqlc generate
```

`sqlc` utiliza:

- `db/schema/schema.sql` para conocer la estructura de la tabla;
- `db/queries/queries.sql` para conocer las operaciones SQL;
- `sqlc.yaml` para conocer el motor, las rutas y el paquete Go que debe generar.

Si el comando finaliza sin errores, se genera o actualiza:

```text
db/sqlc/
```

---

# Esquema de datos

La tabla principal es `workouts`.

Cada fila representa un registro de entrenamiento:

```sql
CREATE TABLE workouts (
    id SERIAL PRIMARY KEY,
    exercise VARCHAR(255) NOT NULL,
    sets INTEGER NOT NULL,
    repetitions INTEGER NOT NULL,
    weight NUMERIC(6,2) NOT NULL,
    workout_date DATE NOT NULL
);
```

Los campos representan:

- `id`: identificador único del registro;
- `exercise`: nombre del ejercicio;
- `sets`: cantidad de series;
- `repetitions`: cantidad de repeticiones;
- `weight`: peso utilizado;
- `workout_date`: fecha del entrenamiento.

---

# Consultas disponibles

En `db/queries/queries.sql` se encuentran las operaciones CRUD requeridas:

```text
CreateWorkout
GetWorkout
ListWorkouts
UpdateWorkout
DeleteWorkout
```

Las consultas utilizan las anotaciones de `sqlc`:

```sql
-- name: CreateWorkout :one
-- name: GetWorkout :one
-- name: ListWorkouts :many
-- name: UpdateWorkout :exec
-- name: DeleteWorkout :exec
```

Donde:

- `:one` indica que la consulta devuelve una fila;
- `:many` indica que puede devolver varias filas;
- `:exec` indica que la operación se ejecuta sin devolver filas.

---

# Verificación de la entrega

Desde la raíz del proyecto ejecutar:

```bash
sqlc generate
```

No debe devolver errores.

Luego comprobar que el código fue generado:

### Ubuntu / Linux

```bash
find db/sqlc -type f
```

### Windows PowerShell

```powershell
Get-ChildItem db/sqlc
```

Deberían aparecer archivos Go generados dentro de `db/sqlc/`.

También se puede comprobar que los paquetes Go compilan correctamente:

```bash
go test ./...
```

Y realizar una verificación adicional:

```bash
go vet ./...
```

---

## Comprobación completa desde cero

Para verificar que el proyecto puede regenerar completamente la capa de acceso a datos, se puede eliminar **únicamente el directorio generado** `db/sqlc/` y volver a ejecutar `sqlc generate`.

### Ubuntu / Linux

```bash
rm -rf db/sqlc
sqlc generate
```

### Windows PowerShell

```powershell
Remove-Item -Recurse -Force db/sqlc
sqlc generate
```

Luego comprobar nuevamente que `db/sqlc/` fue creado.

Esto confirma que los archivos fuente necesarios son suficientes:

```text
db/schema/schema.sql
db/queries/queries.sql
sqlc.yaml
```

---

# ¿Qué hace este proyecto?

La aplicación Reps necesita guardar registros de entrenamiento de forma persistente.

En esta entrega todavía no se conecta la aplicación a una base de datos real. Primero se prepara la capa de datos:

```text
schema.sql
    +
queries.sql
    +
sqlc.yaml
    ↓
sqlc generate
    ↓
código Go tipado
    ↓
db/sqlc/
```

El esquema define **qué datos existen**.

Las consultas definen **qué operaciones se podrán realizar**.

`sqlc` toma ambos archivos y genera automáticamente el código Go necesario para ejecutar esas consultas de forma tipada.

De esta forma seguimos escribiendo y controlando el SQL, pero evitamos escribir manualmente código repetitivo de acceso a datos como `QueryRow`, `Query`, `Scan` o `Exec`.

El paquete generado queda preparado para ser utilizado en el siguiente práctico.
