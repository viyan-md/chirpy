-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING id, created_at, updated_at, email, is_chirpy_red;

-- name: GetUserById :one
SELECT * FROM users
WHERE id = $1;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUserFromRefreshToken :one
SELECT u.id,
       u.email,
       u.created_at,
       u.updated_at,
       u.is_chirpy_red
FROM   users           AS u
JOIN   refresh_tokens  AS rt ON rt.user_id = u.id
WHERE  rt.token      = $1
  AND  rt.revoked_at IS NULL
  AND  rt.expires_at > NOW();

-- name: UpdateCredentials :one 
UPDATE users
SET email = $1, hashed_password = $2, updated_at = NOW()
WHERE id = $3
RETURNING id, created_at, updated_at, email, is_chirpy_red;

-- name: UpgradeUserToRed :exec
UPDATE users
SET is_chirpy_red = 'true', updated_at=NOW()
WHERE id = $1;