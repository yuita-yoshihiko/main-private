-- name: CreateItem :one
INSERT INTO items (name, description)
VALUES ($1, $2)
RETURNING *;

-- name: GetItem :one
SELECT * FROM items
WHERE id = $1;

-- name: ListItems :many
SELECT * FROM items
ORDER BY created_at DESC;

-- name: UpdateItem :one
UPDATE items
SET name = $1, description = $2, updated_at = NOW()
WHERE id = $3
RETURNING *;

-- name: DeleteItem :exec
DELETE FROM items
WHERE id = $1;
