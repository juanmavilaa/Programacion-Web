# Pruebas de la API de Reps
# Correr con la base vacia (los ids arrancan en 1): hurl --test requests.hurl


# --- Clientes ---

GET http://localhost:8080/clients
HTTP 200
[Asserts]
header "Content-Type" == "application/json"

# crear cliente
POST http://localhost:8080/clients
Content-Type: application/json
{
  "name": "Juan Perez",
  "email": "juan@reps.com"
}
HTTP 201
[Asserts]
jsonpath "$.id" == 1
jsonpath "$.name" == "Juan Perez"
jsonpath "$.email" == "juan@reps.com"

# sin nombre
POST http://localhost:8080/clients
Content-Type: application/json
{
  "name": "",
  "email": "sinnombre@reps.com"
}
HTTP 400

# email mal escrito
POST http://localhost:8080/clients
Content-Type: application/json
{
  "name": "Ana",
  "email": "ana.reps.com"
}
HTTP 400

# email que ya existe
POST http://localhost:8080/clients
Content-Type: application/json
{
  "name": "Otro Juan",
  "email": "juan@reps.com"
}
HTTP 409

GET http://localhost:8080/clients
HTTP 200
[Asserts]
jsonpath "$[0].id" == 1
jsonpath "$[0].name" == "Juan Perez"

GET http://localhost:8080/clients/1
HTTP 200
[Asserts]
jsonpath "$.name" == "Juan Perez"
jsonpath "$.email" == "juan@reps.com"

GET http://localhost:8080/clients/999
HTTP 404

GET http://localhost:8080/clients/abc
HTTP 400

# actualizar cliente
PUT http://localhost:8080/clients/1
Content-Type: application/json
{
  "name": "Juan Carlos Perez",
  "email": "jcperez@reps.com"
}
HTTP 200
[Asserts]
jsonpath "$.id" == 1
jsonpath "$.name" == "Juan Carlos Perez"
jsonpath "$.email" == "jcperez@reps.com"

PUT http://localhost:8080/clients/1
Content-Type: application/json
{
  "name": "",
  "email": "jcperez@reps.com"
}
HTTP 400

PUT http://localhost:8080/clients/999
Content-Type: application/json
{
  "name": "Nadie",
  "email": "nadie@reps.com"
}
HTTP 404


# --- Entrenamientos ---

# crear entrenamiento para el cliente 1
POST http://localhost:8080/workouts
Content-Type: application/json
{
  "client_id": 1,
  "exercise": "Press banca",
  "sets": 4,
  "repetitions": 8,
  "weight": "80.00"
}
HTTP 201
[Asserts]
jsonpath "$.id" == 1
jsonpath "$.client_id" == 1
jsonpath "$.exercise" == "Press banca"
jsonpath "$.sets" == 4
jsonpath "$.repetitions" == 8
jsonpath "$.weight" == "80.00"

# cliente que no existe
POST http://localhost:8080/workouts
Content-Type: application/json
{
  "client_id": 999,
  "exercise": "Sentadilla",
  "sets": 3,
  "repetitions": 10,
  "weight": "100.00"
}
HTTP 400

# series en 0
POST http://localhost:8080/workouts
Content-Type: application/json
{
  "client_id": 1,
  "exercise": "Sentadilla",
  "sets": 0,
  "repetitions": 10,
  "weight": "100.00"
}
HTTP 400

# peso que no es un numero
POST http://localhost:8080/workouts
Content-Type: application/json
{
  "client_id": 1,
  "exercise": "Sentadilla",
  "sets": 3,
  "repetitions": 10,
  "weight": "cien"
}
HTTP 400

GET http://localhost:8080/workouts
HTTP 200
[Asserts]
header "Content-Type" == "application/json"
jsonpath "$[0].id" == 1
jsonpath "$[0].exercise" == "Press banca"

GET http://localhost:8080/workouts/1
HTTP 200
[Asserts]
jsonpath "$.exercise" == "Press banca"
jsonpath "$.client_id" == 1

# actualizar entrenamiento
PUT http://localhost:8080/workouts/1
Content-Type: application/json
{
  "client_id": 1,
  "exercise": "Press inclinado",
  "sets": 4,
  "repetitions": 10,
  "weight": "70.50"
}
HTTP 200
[Asserts]
jsonpath "$.id" == 1
jsonpath "$.exercise" == "Press inclinado"
jsonpath "$.repetitions" == 10
jsonpath "$.weight" == "70.50"

# el cliente tiene un entrenamiento, no se puede borrar
DELETE http://localhost:8080/clients/1
HTTP 409

DELETE http://localhost:8080/workouts/1
HTTP 204

GET http://localhost:8080/workouts/1
HTTP 404


# --- Borrado del cliente ---

DELETE http://localhost:8080/clients/1
HTTP 204

GET http://localhost:8080/clients/1
HTTP 404

DELETE http://localhost:8080/clients/1
HTTP 404

PATCH http://localhost:8080/clients
HTTP 405
