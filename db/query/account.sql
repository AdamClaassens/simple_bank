-- name: CreateAccount :one
INSERT INTO accounts (
    owner, balance, currency
) VALUES (
             $1, $2, $3
         )
RETURNING *;

-- name: GetAccount :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1;

/*
Create a get account function we can use concurrently
This is used so that when a transaction is busy on the database, we do not try to read while other transactions are writing
We have to state NO KEY UPDATE so that we do not lock the row for writing when we try and update the balance
*/
-- name: GetAccountForUpdate :one
SELECT * FROM accounts
WHERE id = $1 LIMIT 1
FOR NO KEY UPDATE;

-- name: ListAccounts :many
SELECT * FROM accounts
ORDER BY id
LIMIT $1
OFFSET $2;

-- name: UpdateAccount :one
UPDATE accounts
SET balance = $2
WHERE id = $1
RETURNING *;

-- name: AddAccountBalance :one
UPDATE accounts
SET balance = ROUND((balance + sqlc.arg(amount))::numeric, 2)::float
WHERE id = sqlc.arg(id)
RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE id = $1;