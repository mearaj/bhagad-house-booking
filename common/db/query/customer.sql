-- name: CreateCustomer :one
INSERT INTO customers (name,
                       address,
                       phone,
                       email)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetCustomer :one
SELECT *
FROM customers
WHERE id = $1
LIMIT 1;

-- name: ListCustomers :many
SELECT *
FROM customers
ORDER BY name
LIMIT $1 OFFSET $2;

-- name: UpdateCustomer :one
UPDATE customers
SET name    = $2,
    address    = $3,
    phone    = $4,
    email    = $5
WHERE id = $1
RETURNING *;

-- name: DeleteCustomer :exec
DELETE
FROM customers
WHERE id = $1;