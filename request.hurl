# requests.hurl - Pruebas de la API REST de Reps
# Requiere la base vacía (make test la recrea) y la API en http://localhost:8080
# Ejecutar: hurl --test requests.hurl

# ======================= CLIENTES =======================

# 1. Listar clientes con la base vacía
GET http://localhost:8080/clients
HTTP 200
[Asserts]
header "Content-Type" == "application/json"

# 2. Crear un cliente
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

# 3. Crear un cliente sin nombre
POST http://localhost:8080/clients
Content-Type: application/json
{
  "name": "",
  "email": "otro@reps.com"
}
HTTP 400

# 4. Crear un cliente con email inválido
POST http://localhost:8080/clients
Content-Type: application/json
{
  "name": "Ana",
  "email": "ana-sin-arroba"
}
HTTP 400

# 5. Crear un cliente con email repetido
POST http://localhost:8080/clients
Content-Type: application/json
{
  "name": "Juan Duplicado",
  "email": "juan@reps.com"
}
HTTP 409

# 6. Listar clientes
GET http://localhost:8080/clients
HTTP 200
[Asserts]
jsonpath "$[0].id" == 1
jsonpath "$[0].name" == "Juan Perez"

# 7. Obtener el cliente 1
GET http://localhost:8080/clients/1
HTTP 200
[Asserts]
jsonpath "$.name" == "Juan Perez"
jsonpath "$.email" == "juan@reps.com"

# 8. Obtener un cliente inexistente
GET http://localhost:8080/clients/999
HTTP 404

# 9. Obtener un cliente con ID inválido
GET http://localhost:8080/clients/abc
HTTP 400

# 10. Actualizar el cliente 1
PUT http://localhost:8080/clients/1
Content-Type: application/json
{
  "name": "Juan Actualizado",
  "email": "juan.actualizado@reps.com"
}
HTTP 200
[Asserts]
jsonpath "$.id" == 1
jsonpath "$.name" == "Juan Actualizado"
jsonpath "$.email" == "juan.actualizado@reps.com"

# 11. Actualizar el cliente 1 con datos inválidos
PUT http://localhost:8080/clients/1
Content-Type: application/json
{
  "name": "",
  "email": "juan.actualizado@reps.com"
}
HTTP 400

# 12. Actualizar un cliente inexistente
PUT http://localhost:8080/clients/999
Content-Type: application/json
{
  "name": "Nadie",
  "email": "nadie@reps.com"
}
HTTP 404

# ===================== ENTRENAMIENTOS =====================

# 13. Crear un entrenamiento para el cliente 1
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

# 14. Crear un entrenamiento para un cliente inexistente
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

# 15. Crear un entrenamiento con series inválidas
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

# 16. Crear un entrenamiento con peso inválido
POST http://localhost:8080/workouts
Content-Type: application/json
{
  "client_id": 1,
  "exercise": "Sentadilla",
  "sets": 3,
  "repetitions": 10,
  "weight": "mucho"
}
HTTP 400

# 17. Listar entrenamientos
GET http://localhost:8080/workouts
HTTP 200
[Asserts]
header "Content-Type" == "application/json"
jsonpath "$[0].id" == 1
jsonpath "$[0].exercise" == "Press banca"

# 18. Obtener el entrenamiento 1
GET http://localhost:8080/workouts/1
HTTP 200
[Asserts]
jsonpath "$.exercise" == "Press banca"
jsonpath "$.client_id" == 1

# 19. Actualizar el entrenamiento 1
PUT http://localhost:8080/workouts/1
Content-Type: application/json
{
  "client_id": 1,
  "exercise": "Press banca inclinado",
  "sets": 4,
  "repetitions": 10,
  "weight": "75.50"
}
HTTP 200
[Asserts]
jsonpath "$.id" == 1
jsonpath "$.exercise" == "Press banca inclinado"
jsonpath "$.repetitions" == 10
jsonpath "$.weight" == "75.50"

# 20. No se puede borrar un cliente con entrenamientos
DELETE http://localhost:8080/clients/1
HTTP 409

# 21. Eliminar el entrenamiento 1
DELETE http://localhost:8080/workouts/1
HTTP 204

# 22. El entrenamiento 1 ya no existe
GET http://localhost:8080/workouts/1
HTTP 404

# ======================= LIMPIEZA =======================

# 23. Eliminar el cliente 1
DELETE http://localhost:8080/clients/1
HTTP 204

# 24. El cliente 1 ya no existe
GET http://localhost:8080/clients/1
HTTP 404

# 25. Eliminar un cliente inexistente
DELETE http://localhost:8080/clients/1
HTTP 404

# 26. Método no permitido
PATCH http://localhost:8080/clients
HTTP 405