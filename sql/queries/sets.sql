-- name: CreateSet :one
INSERT INTO sets (exercise, count, intensity, type, weight, workout_id, created_at, updated_at)
values($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetSetById :one
SELECT * FROM sets where id = $1;

-- name: GetSetsByWorkoutIdDesc :many
SELECT * FROM sets where workout_id = sqlc.arg('workout_id')::int
ORDER BY sqlc.arg('order_by_col')::text DESC
LIMIT sqlc.arg('limit')::int
OFFSET sqlc.arg('offset')::int;

-- name: GetSetsByWorkoutIdAsc :many
SELECT * FROM sets where workout_id = sqlc.arg('workout_id')::int
ORDER BY sqlc.arg('order_by_col')::text ASC
LIMIT sqlc.arg('limit')::int
OFFSET sqlc.arg('offset')::int;