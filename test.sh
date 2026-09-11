#!/usr/bin/env bash

set -e

cleanup() {
    docker compose down -v > /dev/null 2>&1 || true
}

trap cleanup EXIT

echo "Borrando contenedores y volumenes anteriores..."
docker compose down -v > /dev/null 2>&1 || true

echo "Generando codigo con sqlc..."
sqlc generate

echo "Compilando..."
go build ./...

echo "Levantando PostgreSQL..."
docker compose up -d database

echo "Esperando a PostgreSQL..."
until docker compose exec -T database psql -U postgres -d reps -c "SELECT 1;" > /dev/null 2>&1
do
    sleep 1
done

echo "Creando tablas..."
docker compose exec -T database psql -U postgres -d reps < db/schema/schema.sql

echo "Ejecutando tests..."
go test -v ./...