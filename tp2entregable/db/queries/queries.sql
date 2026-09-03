-- name: CreateClient :one
INSERT INTO clients (name, email)
VALUES ($1, $2)
RETURNING id, name, email;

-- name: GetClient :one
SELECT id, name, email
FROM clients
WHERE id = $1;

-- name: ListClients :many
SELECT id, name, email
FROM clients
ORDER BY name;

-- name: UpdateClient :exec
UPDATE clients
SET name = $2, email = $3
WHERE id = $1;

-- name: DeleteClient :exec
DELETE FROM clients
WHERE id = $1;


-- name: CreateWorkout :one
INSERT INTO workouts (client_id, exercise, sets, repetitions, weight)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, client_id, exercise, sets, repetitions, weight, created_at;

-- name: GetWorkout :one
SELECT id, client_id, exercise, sets, repetitions, weight, created_at
FROM workouts
WHERE id = $1;

-- name: ListWorkouts :many
SELECT id, client_id, exercise, sets, repetitions, weight, created_at
FROM workouts
ORDER BY created_at DESC;

-- name: ListWorkoutsByClient :many
SELECT id, client_id, exercise, sets, repetitions, weight, created_at
FROM workouts
WHERE client_id = $1
ORDER BY created_at DESC;

-- name: UpdateWorkout :exec
UPDATE workouts
SET client_id = $2,
    exercise = $3,
    sets = $4,
    repetitions = $5,
    weight = $6
WHERE id = $1;

-- name: DeleteWorkout :exec
DELETE FROM workouts
WHERE id = $1;