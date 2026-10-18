APP_NAME := reps
API_URL  := http://localhost:8080

.PHONY: all generate build down db api run test test-api clean

all: build

# Ejecuta las tareas de generación configuradas por el proyecto
generate:
	@sqlc generate

# Construye el binario de la aplicación
build: generate
	@mkdir -p tmp
	@go build -o tmp/$(APP_NAME) .

# Borra contenedores y volumenes
down:
	@docker compose down -v

# Levanta PostgreSQL y espera a que acepte conexiones
db:
	@docker compose up -d database
	@echo "Esperando a PostgreSQL..."
	@until docker compose exec -T database pg_isready -h localhost -U postgres -d reps > /dev/null 2>&1; do \
		sleep 1; \
	done

# Levanta la API y espera a que responda
api: db
	@docker compose up -d --build api
	@echo "Esperando a la API..."
	@for i in $$(seq 1 30); do \
		curl -s -o /dev/null $(API_URL)/clients && exit 0; \
		sleep 1; \
	done; \
	echo "La API no levanto"; docker compose logs api; exit 1

run:
	@docker compose up --build

# Tests de persistencia (TP2)
test: down build db
	@ok=0; \
	go test -v -count=1 ./... || ok=1; \
	docker compose down -v; \
	exit $$ok

# Prueba todos los endpoints de la API con Hurl
test-api: down build api
	@ok=0; \
	echo "Probando los endpoints..."; \
	hurl --test requests.hurl || ok=1; \
	docker compose down -v; \
	exit $$ok

# Limpia los artefactos de construcción
clean: down
	@rm -rf tmp
