-- name: CreateSchool :one
INSERT INTO schools (public_id, name, short_name, school_type, city, phone, email)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetSchool :one
SELECT * FROM schools
WHERE id = $1
LIMIT 1;

-- name: GetSchoolByPublicID :one
SELECT * FROM schools
WHERE public_id = $1
LIMIT 1;

-- name: ListSchools :many
SELECT * FROM schools
ORDER BY created_at DESC
LIMIT $1 OFFSET $2;
