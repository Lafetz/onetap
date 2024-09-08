-- name: CreateTierLevel :exec
INSERT INTO tier_level (tier_id, merchant_id, name, min_points, created_at, updated_at)
VALUES ($1, $2, $3, $4, NOW(), NOW());


-- name: GetTierLevel :one
SELECT merchant_id, name, min_points, created_at, updated_at
FROM tier_level
WHERE merchant_id = $1 AND name = $2;


-- name: UpdateTierLevel :exec
UPDATE tier_level
SET name = $2, min_points = $3, updated_at = NOW()
WHERE tier_id = $1;

-- name: DeleteTierLevel :exec
DELETE FROM tier_level
WHERE merchant_id = $1 AND name = $2;


-- name: ListTierLevels :many
SELECT merchant_id, name, min_points, created_at, updated_at
FROM tier_level
WHERE merchant_id = $1
ORDER BY min_points ASC;


-- name: CreateCustomerTier :exec
INSERT INTO customer_tier (merchant_id, customer_id, tier_name, points, created_at, updated_at)
VALUES ($1, $2, $3, $4, NOW(), NOW());


-- name: GetCustomerTier :one
SELECT merchant_id, customer_id, tier_name, points, created_at, updated_at
FROM customer_tier
WHERE merchant_id = $1 AND customer_id = $2;


-- name: UpdateCustomerTier :exec
UPDATE customer_tier
SET tier_name = $3, points = $4, updated_at = NOW()
WHERE merchant_id = $1 AND customer_id = $2;


-- name: DeleteCustomerTier :exec
DELETE FROM customer_tier
WHERE merchant_id = $1 AND customer_id = $2;
