APP_NAME := reps
API_URL := http://localhost:8080

.PHONY: all generate build up down run test clean

all: build

generate:
	@sqlc generate

build: generate
	@mkdir -p tmp
	@go build -o tmp/$(APP_NAME) .

up:
	@docker compose up -d --build
	@echo "Esperando a la API..."
	@for i in $$(seq 1 60); do \
		curl -s -o /dev/null $(API_URL)/clients && echo "API lista" && exit 0; \
		sleep 1; \
	done; \
	echo "La API no levanto"; docker compose logs api; exit 1

down:
	@docker compose down -v

run:
	@docker compose up --build

# si falla algun test igual se borran los contenedores al final
test: down build up
	@ok=0; \
	echo "Ejecutando pruebas de la API..."; \
	hurl --test requests.hurl || ok=1; \
	echo "Ejecutando tests de persistencia..."; \
	go test -v -count=1 ./... || ok=1; \
	docker compose down -v; \
	exit $$ok

clean: down
	@rm -rf tmp
