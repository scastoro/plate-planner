-- name: CreateUser :one
INSERT INTO "Admin"."Users" (id, first_name, last_name, body_weight, username, email, password, lastLoggedIn)
values($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetUserById :one
SELECT * FROM "Admin"."Users" where id = $1;