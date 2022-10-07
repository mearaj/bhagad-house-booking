CREATE TYPE "rate_time_units" AS ENUM (
    'Day',
    'Hour',
    'Week',
    'Month'
    );

CREATE TABLE "bookings"
(
    "id"             bigserial PRIMARY KEY,
    "created_at"     timestamptz DEFAULT (now()),
    "updated_at"     timestamptz DEFAULT (now()),
    "start_date"     timestamptz,
    "end_date"       timestamptz,
    "customer_id"    bigint,
    "rate"           double precision,
    "rate_time_unit" rate_time_units NOT NULL DEFAULT 'Day'
);

CREATE TABLE "customers"
(
    "id"         bigserial PRIMARY KEY,
    "created_at" timestamptz DEFAULT (now()),
    "updated_at" timestamptz DEFAULT (now()),
    "name"       varchar NOT NULL,
    "address"    varchar NOT NULL,
    "phone"      varchar NOT NULL,
    "email"      varchar NOT NULL
);

CREATE INDEX ON "bookings" ("start_date");

CREATE INDEX ON "bookings" ("end_date");

CREATE INDEX ON "bookings" ("start_date", "end_date");

CREATE INDEX ON "bookings" ("customer_id");

CREATE INDEX ON "customers" ("name");

CREATE INDEX ON "customers" ("address");

CREATE INDEX ON "customers" ("phone");

CREATE INDEX ON "customers" ("email");

ALTER TABLE "bookings"
    ADD FOREIGN KEY ("customer_id") REFERENCES "customers" ("id");
