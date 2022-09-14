-- name: CreateOauthUser :one
INSERT INTO oauths (
    id,
    fullname,
    email,
    provider,
    refresh_token
) VALUES (
    $1, $2, $3, $4, $5
) RETURNING *;

-- name: GetOauthUser :one
SELECT * FROM oauths 
WHERE id = $1 LIMIT 1;