-- name: GetAccount :one
SELECT * FROM aws_account
WHERE account_id = $1;

-- name: CreateAccount :one
INSERT INTO aws_account (account_id, account_type)
VALUES ($1, $2) RETURNING *;

-- name: DeleteAccount :exec
DELETE FROM aws_account
where account_id = $1;
