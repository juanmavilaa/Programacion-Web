-- name: CreateWorkout :one
INSERT INTO workouts (
    exercise,
    sets,
    repetitions,
    weight,
    workout_date
)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, exercise, sets, repetitions, weight, workout_date;

-- name: GetWorkout :one
SELECT id, exercise, sets, repetitions, weight, workout_date
FROM workouts
WHERE id = $1;

-- name: ListWorkouts :many
SELECT id, exercise, sets, repetitions, weight, workout_date
FROM workouts
ORDER BY workout_date DESC;

-- name: UpdateWorkout :exec
UPDATE workouts
SET
    exercise = $2,
    sets = $3,
    repetitions = $4,
    weight = $5,
    workout_date = $6
WHERE id = $1;

-- name: DeleteWorkout :exec
DELETE FROM workouts
WHERE id = $1;