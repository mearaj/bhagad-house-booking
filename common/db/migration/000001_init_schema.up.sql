CREATE TYPE "rate_time_units" AS ENUM (
    'Day',
    'Hour',
    'Week',
    'Month'
    );

CREATE TYPE "user_roles" AS ENUM (
    'User',
    'Admin'
    );

CREATE TABLE "bookings"
(
    "id"             bigserial PRIMARY KEY,
    "created_at"     timestamptz      NOT NULL DEFAULT (now()),
    "updated_at"     timestamptz      NOT NULL DEFAULT (now()),
    "start_date"     timestamptz      NOT NULL,
    "end_date"       timestamptz      NOT NULL,
    "customer_id"    bigint           NOT NULL,
    "rate"           double precision NOT NULL,
    "rate_time_unit" rate_time_units  NOT NULL DEFAULT 'Day'
);

CREATE TABLE "customers"
(
    "id"         bigserial PRIMARY KEY,
    "created_at" timestamptz NOT NULL DEFAULT (now()),
    "updated_at" timestamptz NOT NULL DEFAULT (now()),
    "name"       varchar     NOT NULL,
    "address"    varchar     NOT NULL,
    "phone"      varchar     NOT NULL,
    "email"      varchar     NOT NULL
);

CREATE TABLE "users"
(
    "id"                  bigserial PRIMARY KEY,
    "password"            varchar        NOT NULL,
    "name"                varchar        NOT NULL,
    "email"               varchar UNIQUE NOT NULL,
    "email_verified"      boolean        NOT NULL DEFAULT false,
    "password_changed_at" timestamptz    NOT NULL DEFAULT ('0001-01-01 00:00:00Z'),
    "created_at"          timestamptz    NOT NULL DEFAULT (now()),
    "roles"               user_roles[]   NOT NULL DEFAULT ('{User}')::user_roles[]
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
