-- name: CreateUser :one
INSERT INTO users (email, password_hash, name, email_verified, image)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: UpdateUserImage :exec
UPDATE users SET image = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1 LIMIT 1;

-- name: GetUserById :one
SELECT * FROM users WHERE id = $1 LIMIT 1;

-- name: UpdateUserEmailVerified :exec
UPDATE users SET email_verified = TRUE, updated_at = CURRENT_TIMESTAMP WHERE id = $1;

-- name: UpdateUserPassword :exec
UPDATE users SET password_hash = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2;

-- name: CreateRefreshToken :one
INSERT INTO refresh_tokens (user_id, token_hash, expires_at)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetRefreshToken :one
SELECT * FROM refresh_tokens WHERE token_hash = $1 AND revoked = FALSE AND expires_at > CURRENT_TIMESTAMP LIMIT 1;

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens SET revoked = TRUE WHERE token_hash = $1;

-- name: RevokeAllUserRefreshTokens :exec
UPDATE refresh_tokens SET revoked = TRUE WHERE user_id = $1;

-- name: CreateVerificationToken :one
INSERT INTO verification_tokens (user_id, token, type, expires_at)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetVerificationToken :one
SELECT * FROM verification_tokens WHERE token = $1 AND expires_at > CURRENT_TIMESTAMP LIMIT 1;

-- name: DeleteVerificationToken :exec
DELETE FROM verification_tokens WHERE token = $1;

-- name: CreateOAuthAccount :one
INSERT INTO oauth_accounts (user_id, provider, provider_user_id, access_token, refresh_token, expires_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetOAuthAccount :one
SELECT * FROM oauth_accounts WHERE provider = $1 AND provider_user_id = $2 LIMIT 1;

-- name: GetOAuthAccountByUserId :one
SELECT * FROM oauth_accounts WHERE user_id = $1 AND provider = $2 LIMIT 1;

-- name: UpdateOAuthTokens :exec
UPDATE oauth_accounts SET access_token = $1, refresh_token = $2, expires_at = $3 WHERE id = $4;
