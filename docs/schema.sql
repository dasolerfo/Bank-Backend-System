-- SQL dump generated using DBML (dbml.dbdiagram.io)
-- Database: PostgreSQL
-- Generated at: 2025-07-18T14:52:26.708Z

CREATE TYPE "Currency" AS ENUM (
  'USD',
  'EUR',
  'KRW',
  'JPY'
);

CREATE TABLE "owners" (
  "id" BIGINT PRIMARY KEY,
  "first_name" VARCHAR NOT NULL,
  "first_surname" VARCHAR NOT NULL,
  "second_surname" VARCHAR NOT NULL,
  "born_at" DATE NOT NULL,
  "nationality" INT NOT NULL,
  "hashed_password" VARCHAR NOT NULL DEFAULT '12345678',
  "email" VARCHAR UNIQUE NOT NULL,
  "created_at" "TIMESTAMPTZ" NOT NULL DEFAULT (NOW()),
  "password_changed_at" "TIMESTAMPTZ" NOT NULL DEFAULT '0001-01-01 00:00:00Z'
);

CREATE TABLE "accounts" (
  "id" BIGINT PRIMARY KEY,
  "owner_id" BIGINT NOT NULL,
  "currency" "Currency" NOT NULL,
  "created_at" "TIMESTAMPTZ" DEFAULT (NOW()),
  "money" FLOAT NOT NULL,
  "country_code" INT NOT NULL
);

CREATE TABLE "entries" (
  "id" BIGINT PRIMARY KEY,
  "account_id" BIGINT,
  "amount" FLOAT NOT NULL,
  "created_at" "TIMESTAMPTZ" DEFAULT (NOW())
);

CREATE TABLE "transfers" (
  "id" BIGINT PRIMARY KEY,
  "from_account_id" BIGINT NOT NULL,
  "to_account_id" BIGINT NOT NULL,
  "amount" FLOAT NOT NULL,
  "created_at" "TIMESTAMPTZ" DEFAULT (NOW())
);

CREATE TABLE "sessions" (
  "id" UUID PRIMARY KEY,
  "owner_id" BIGINT NOT NULL,
  "email" VARCHAR NOT NULL,
  "refresh_token" VARCHAR NOT NULL,
  "user_agent" VARCHAR NOT NULL,
  "client_ip" VARCHAR NOT NULL,
  "is_blocked" BOOLEAN NOT NULL DEFAULT false,
  "created_at" "TIMESTAMPTZ" NOT NULL DEFAULT (NOW()),
  "expires_at" "TIMESTAMPTZ" NOT NULL
);

CREATE INDEX ON "accounts" ("owner_id");

CREATE INDEX ON "entries" ("account_id");

CREATE INDEX ON "transfers" ("from_account_id");

CREATE INDEX ON "transfers" ("to_account_id");

CREATE INDEX ON "transfers" ("from_account_id", "to_account_id");

COMMENT ON TABLE "accounts" IS 'UNIQUE (owner_id, currency) constraint';

COMMENT ON COLUMN "entries"."amount" IS 'Can be negative or positive';

COMMENT ON COLUMN "transfers"."amount" IS 'Must be positive';

ALTER TABLE "accounts" ADD CONSTRAINT "fk_owner" FOREIGN KEY ("owner_id") REFERENCES "owners" ("id") ON DELETE CASCADE;

ALTER TABLE "entries" ADD CONSTRAINT "fk_account" FOREIGN KEY ("account_id") REFERENCES "accounts" ("id") ON DELETE CASCADE;

ALTER TABLE "transfers" ADD CONSTRAINT "fk_from_account" FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id") ON DELETE CASCADE;

ALTER TABLE "transfers" ADD CONSTRAINT "fk_to_account" FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id") ON DELETE CASCADE;

ALTER TABLE "sessions" ADD CONSTRAINT "fk_session_owner" FOREIGN KEY ("owner_id") REFERENCES "owners" ("id") ON DELETE CASCADE;
