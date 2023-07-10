-- name: CreateWorkout :one
INSERT INTO workouts (start_time, duration, total_weight, total_calories, user_id)
values($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetWorkoutById :one
SELECT * 
FROM workouts 
WHERE user_id = $1;

-- name: GetWorkoutsByUserIdDesc :many
SELECT count(*) OVER(), * 
FROM workouts 
WHERE user_id = sqlc.arg('user_id')::int
ORDER BY sqlc.arg('order_by_col')::text DESC
LIMIT sqlc.arg('limit')::int
OFFSET sqlc.arg('offset')::int;

-- name: GetWorkoutsByUserIdAsc :many
SELECT count(*) OVER(), * 
FROM workouts 
WHERE user_id = sqlc.arg('user_id')::int
ORDER BY sqlc.arg('order_by_col')::text ASC
LIMIT sqlc.arg('limit')::int
OFFSET sqlc.arg('offset')::int;

-- name: UpdateWorkoutById :one
UPDATE workouts SET (duration, total_weight, total_calories) 
    = ($1, $2, $3)
WHERE id = $4
RETURNING *;

-- name: GetWorkoutsByIdWithSets :many
SELECT wo.*, s.*
FROM workouts AS wo
JOIN sets AS s 
ON s.workout_id = wo.id
WHERE wo.id = $1;