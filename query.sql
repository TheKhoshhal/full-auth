-- name: GetAccountUsername :one
SELECT username FROM accounts
WHERE username = ?;

-- name: GetAccountPassword :one
SELECT password FROM accounts
WHERE username = ?;

-- name: CraeteAccount :one
INSERT INTO accounts (
  username, password, email
) VALUES (
  ?, ?, ?
)
RETURNING username;

-- name: DeleteAccount :exec
DELETE FROM accounts
WHERE username = ?;
