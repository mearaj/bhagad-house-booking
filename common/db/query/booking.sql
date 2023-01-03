-- name: CreateBooking :one
INSERT INTO bookings (start_date,
                      end_date,
                      details)
VALUES ($1, $2, $3)
RETURNING *;

-- name: GetBooking :one
SELECT *
FROM bookings
WHERE id = $1
LIMIT 1;

-- name: ListBookings :many
SELECT *
FROM bookings
WHERE (start_date <= $1 AND end_date >= $2) OR
    (start_date >= $1 AND end_date >= $2) OR
    (start_date <= $1 AND end_date <= $2) OR
    (start_date >= $1 AND end_date <= $2)
ORDER BY start_date;

-- name: UpdateBooking :one
UPDATE bookings
SET start_date     = $2,
    end_date       = $3,
    details    = $4
WHERE id = $1
RETURNING *;

-- name: GetConflictingBookings :many
SELECT *
FROM bookings
WHERE (start_date <= $1 AND end_date >= $1) OR
    (end_date >= $2 AND start_date <= $2);

-- name: DeleteBooking :exec
DELETE
FROM bookings
WHERE id = $1;
