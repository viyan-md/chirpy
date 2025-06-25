-- name: AddRefreshToken :one
INSERT INTO refresh_tokens (token, user_id) 
VALUES (
    $1,
    $2
) RETURNING token;

-- name: RevokeToken :exec
UPDATE refresh_tokens
SET    revoked_at = NOW(),
       updated_at = NOW()
WHERE  token = $1; 