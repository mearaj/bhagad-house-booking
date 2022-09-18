-- name: CreateBooking :one
INSERT INTO bookings (start_date,
                      end_date,
                      customer_id,
                      rate,
                      rate_time_unit)
VALUES ($1, $2, $3,$4,$5)
RETURNING *;

-- name: GetBooking :one
SELECT * FROM bookings
WHERE id = $1 LIMIT 1;

-- name: ListBookings :many
SELECT * FROM bookings
ORDER BY start_date desc
LIMIT $1
OFFSET $2;

-- name: UpdateBooking :one
UPDATE bookings SET start_date = $2, end_date = $3,customer_id = $4, rate = $5, rate_time_unit = $6
WHERE id = $1
RETURNING *;

-- name: DeleteBooking :exec
DELETE FROM bookings WHERE id = $1;