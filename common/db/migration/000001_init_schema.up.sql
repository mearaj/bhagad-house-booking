CREATE TYPE "user_roles" AS ENUM (
    'User',
    'Admin'
    );

CREATE TABLE "bookings"
(
    "id"             BIGSERIAL PRIMARY KEY,
    "created_at"     timestamptz     DEFAULT (now()),
    "updated_at"     timestamptz     DEFAULT (now()),
    "start_date"     timestamptz      NOT NULL,
    "end_date"       timestamptz      NOT NULL,
    "details"    varchar NOT NULL DEFAULT ('')
);

CREATE TABLE "users"
(
    "id"                  BIGSERIAL PRIMARY KEY,
    "password"            varchar        NOT NULL,
    "name"                varchar     NOT NULL   DEFAULT (''),
    "email"               varchar UNIQUE NOT NULL,
    "email_verified"      boolean            NOT NULL     DEFAULT false,
    "password_changed_at" timestamptz NOT NULL            DEFAULT ('0001-01-01 00:00:00Z'),
    "created_at"          timestamptz  NOT NULL  DEFAULT (now()),
    /* "roles" user_roles[] DEFAULT '{user_roles.User}'*/
    "roles"               user_roles[]  NOT NULL DEFAULT ('{User}')::user_roles[]
);

CREATE INDEX ON "bookings" ("start_date");

CREATE INDEX ON "bookings" ("end_date");

CREATE INDEX ON "bookings" ("start_date", "end_date");

ALTER TABLE "bookings"
    ADD CONSTRAINT CheckEndLaterThanStart CHECK (end_date >= start_date);

ALTER TABLE "bookings"
    ADD CONSTRAINT NoOverlappingTimeRanges
    EXCLUDE USING gist (tstzrange(start_date,end_date,'[]') WITH &&);
