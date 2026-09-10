#!/usr/bin/env bash
set -e

echo "=========================================================="
echo "  INICIANDO EJECUCIÓN AUTOMATIZADA CON DOCKER (TP2 REPS)  "
echo "=========================================================="

# ------------------------------------------------------------
# 3. TAREAS POSTERIORES (Garantizadas con trap ante éxito o fallo)
# ------------------------------------------------------------
cleanup() {
    echo ""
    echo ">> [3. TAREAS POSTERIORES] Borrando contenedores y volúmenes temporales de Docker..."
    docker compose down -v 2>/dev/null || true
    echo ">> Limpieza finalizada correctamente."
}
trap cleanup EXIT

# ------------------------------------------------------------
# 1. TAREAS PREVIAS
# ------------------------------------------------------------
echo ""
echo ">> [1. TAREAS PREVIAS]"

# Verificar que Docker esté encendido
if ! docker info > /dev/null 2>&1; then
    echo "❌ ERROR: El servicio de Docker no está en ejecución."
    echo "   Inicia Docker antes de continuar (ej. sudo systemctl start docker)."
    exit 1
fi

# A) Borrar contenedores y volúmenes previos
echo "-> Limpiando contenedores y volúmenes anteriores con docker compose..."
docker compose down -v 2>/dev/null || true

# B) Ejecución de sqlc (si está instalado localmente)
if command -v sqlc &> /dev/null; then
    echo "-> Ejecutando 'sqlc generate'..."
    sqlc generate
else
    echo "-> sqlc no encontrado en PATH. Se utiliza el código generado en db/sqlc/."
fi

# C) Levantar contenedor de base de datos con Docker Compose
echo "-> Levantando servicio de base de datos con docker compose..."
docker compose up -d database

# D) Esperar a que la base de datos esté lista
echo "-> Esperando a que PostgreSQL acepte conexiones..."
until docker compose exec -T database pg_isready -U postgres -d reps > /dev/null 2>&1; do
    sleep 1
done
echo "-> Base de datos PostgreSQL lista y con esquema/semillas aplicados."

# ------------------------------------------------------------
# 2. EJECUCIÓN DE TESTS
# ------------------------------------------------------------
echo ""
echo ">> [2. EJECUCIÓN DE TESTS CON DOCKER]"
# Ejecuta los tests dentro del contenedor 'web' usando el Dockerfile
docker compose run --rm web

echo ""
echo ">> ¡TODOS LOS TESTS PASARON EXITOSAMENTE EN DOCKER!"
