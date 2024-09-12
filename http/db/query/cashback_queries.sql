-- name: CreateCashback :exec
INSERT INTO cashback (id, merchant_id, name, description, percentage, eligible_products, active, expiration, created_at, updated_at)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW());

-- name: GetCashback :one
SELECT * FROM cashback WHERE id = $1;

-- name: UpdateCashback :exec
UPDATE cashback
SET name = $2, description = $3, percentage = $4, eligible_products = $5, active = $6, expiration = $7, updated_at = NOW()
WHERE id = $1;

-- name: DeleteCashback :exec
DELETE FROM cashback WHERE id = $1;

-- name: ListCashbacks :many
SELECT * FROM cashback WHERE merchant_id = $1;

-- name: CreateCashbackUser :exec
INSERT INTO cashback_user (merchant_id,cashback_id, user_id, points)
VALUES ($1, $2, $3,$4);

-- name: GetCashbackUser :one
SELECT * FROM cashback_user WHERE user_id = $1 AND cashback_id = $2;

-- name: UpdateCashbackUser :exec
UPDATE cashback_user
SET points = $2
WHERE user_id = $1 AND cashback_id = $3;

-- name: DeleteCashbackUser :exec
DELETE FROM cashback_user WHERE user_id = $1 AND cashback_id = $2;

-- name: ListCashbackUsers :many
SELECT * FROM cashback_user
WHERE merchant_id = $1;
