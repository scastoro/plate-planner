-- name: CreateWorkout :one
INSERT INTO workouts (start_time, duration, total_weight, total_calories, user_id)
values($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetWorkoutById :one
SELECT * FROM workouts where user_id = $1;

-- name: GetWorkoutsByUserIdDesc :many
SELECT * FROM workouts where user_id = sqlc.arg('user_id')::int
ORDER BY sqlc.arg('order_by_col')::text DESC
LIMIT sqlc.arg('limit')::int
OFFSET sqlc.arg('offset')::int;

-- name: GetWorkoutsByUserIdAsc :many
SELECT * FROM workouts where user_id = sqlc.arg('user_id')::int
ORDER BY sqlc.arg('order_by_col')::text ASC
LIMIT sqlc.arg('limit')::int
OFFSET sqlc.arg('offset')::int;

-- name: UpdateWorkoutById :one
UPDATE workouts SET (duration, total_weight, total_calories) 
    = ($1, $2, $3)
WHERE id = $4
RETURNING *;