-- name: CreateAccounts :one
INSERT INTO accounts (
  owner,
  balance,
  currency
) VALUES (
  $1, $2, $3
)
RETURNING *;

-- name: GetAccounts :one
SELECT * FROM accounts
WHERE id = $1 
LIMIT 1;

-- name: GetAccountsForUpdate :one
SELECT * FROM accounts
WHERE id = $1 
LIMIT 1
FOR NO KEY UPDATE;

-- name: ListAccounts :many
SELECT * FROM accounts
WHERE owner = $1
ORDER BY id
LIMIT $2
OFFSET $3;

-- name: UpdateAccounts :one
UPDATE accounts
set balance = $2
WHERE id = $1
RETURNING *;

-- name: AddBalance :one
UPDATE accounts
set balance = balance + $2
WHERE id = $1
RETURNING *;

-- name: DeleteAccounts :exec
DELETE FROM accounts
WHERE id = $1;