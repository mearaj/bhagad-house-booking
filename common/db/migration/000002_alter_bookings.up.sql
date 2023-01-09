ALTER TABLE IF EXISTS bookings
ADD COLUMN customer_name varchar NOT NULL DEFAULT (''),
ADD COLUMN total_price double precision NOT NULL DEFAULT (0);
