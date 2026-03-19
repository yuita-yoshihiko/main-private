-- name: CreateStats :one
INSERT INTO stats (total_items)
VALUES ($1)
RETURNING *;
