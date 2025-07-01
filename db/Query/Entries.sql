-- name: CreateEntries :one
INSERT INTO entries (
  account_id,
  amount
) VALUES (
  $1, $2
)
RETURNING *;

-- name: GetEntries :one
SELECT * FROM entries
WHERE id = $1 
LIMIT 1;

-- name: ListEntries :many
SELECT * FROM entries
ORDER BY id
OFFSET $1;

-- name: UpdateEntries :exec
UPDATE entries
set amount = $2
WHERE id = $1;

-- name: DeleteEntries :exec
DELETE FROM entries
WHERE account_id = $1;