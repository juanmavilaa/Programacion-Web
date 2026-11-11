APP_NAME := reps
TEMPL    := go run github.com/a-h/templ/cmd/templ@v0.3.1020

.PHONY: all generate build down db run-app test clean

all: build

# Ejecuta las tareas de generacion configuradas por el proyecto
# (templ va con la misma version que go.mod, el codigo que genera depende de la version)
generate:
	@sqlc generate
	@$(TEMPL) generate

# Construye el binario de la aplicacion
build: generate
	@mkdir -p tmp
	@go build -o tmp/$(APP_NAME) .

# Borra contenedores y volumenes
down:
	@docker compose down -v

# Levanta PostgreSQL y espera a que acepte conexiones
db:
	@docker compose up -d database || { docker compose down -v; exit 1; }
	@echo "Esperando a PostgreSQL..."
	@for i in $$(seq 1 30); do \
		docker compose exec -T database pg_isready -h localhost -U postgres -d reps > /dev/null 2>&1 && exit 0; \
		sleep 1; \
	done; \
	echo "PostgreSQL no levanto"; docker compose logs database; docker compose down -v; exit 1

# Ejecuta la aplicacion en http://localhost:8080
# Con Ctrl+C se detiene y se borran los contenedores y volumenes
run-app: down build db
	@echo "Aplicacion en http://localhost:8080 (Ctrl+C para detener)"
	@trap 'docker compose down -v' EXIT; \
	trap 'exit 0' INT TERM; \
	docker compose up --build api

# Tests de persistencia
test: down build db
	@ok=0; \
	go test -v -count=1 ./... || ok=1; \
	docker compose down -v; \
	exit $$ok

# Limpia los artefactos de construccion
clean: down
	@rm -rf tmp
