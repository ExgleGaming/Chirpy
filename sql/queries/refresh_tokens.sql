-- name: CreateRefreshTokens :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at, revoked_at)
VALUES (
         $1,
        Now(),
        Now(),
        $2,
        $3,
        $4
       )
    RETURNING *;

-- name: GetRefreshTokenByToken :one
SELECT * FROM refresh_tokens
WHERE token = $1 AND revoked_at IS NULL AND expires_at > NOW()
LIMIT 1;

-- name: UpdateRefreshToken :exec
UPDATE refresh_tokens
SET revoked_at = NOW(),
    updated_at = NOW()
WHERE token = $1;