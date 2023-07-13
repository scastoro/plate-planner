-- name: CreateUser :one
INSERT INTO "Admin"."Users" (first_name, last_name, body_weight, username, email, password, lastLoggedIn)
values($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetUserById :one
SELECT * FROM "Admin"."Users" where id = $1;

-- name: GetUserByEmail :one
SELECT * FROM "Admin"."Users" where email = $1;

-- name: GetUserByIdWithPerms :many
SELECT u.id, u.first_name, u.last_name, u.username, r.name, p.resource, p.action
FROM "Admin"."Users" as u
JOIN roles as r
ON u.role_id = r.id
JOIN rolesPermissions as rp
ON u.role_id = rp.role_id
JOIN permissions as p
ON p.id = rp.permission_id
WHERE u.id = $1;