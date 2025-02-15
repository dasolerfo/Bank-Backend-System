CREATE TYPE "Currency" AS ENUM (
  'USD',
  'EUR',
  'KRW',
  'JPY'
);

CREATE TABLE "owners" (
  "id" BIGSERIAL PRIMARY KEY,
  "first_name" VARCHAR NOT NULL,
  "first_surname" VARCHAR NOT NULL,
  "second_surname" VARCHAR NOT NULL,
  "born_at" DATE NOT NULL DEFAULT (CURRENT_DATE),
  "nationality" INT NOT NULL
);

CREATE TABLE "accounts" (
  "id" BIGSERIAL PRIMARY KEY,
  "owner_id" BIGSERIAL NOT NULL,
  "currency" "Currency" NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (NOW()) ,
  "money" FLOAT NOT NULL,
  "country_code" INT NOT NULL
);

CREATE TABLE "entries" (
  "id" BIGSERIAL PRIMARY KEY,
  "account_id" BIGSERIAL NOT NULL,
  "amount" FLOAT NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (NOW()) 
);

CREATE TABLE "transfers" (
  "id" BIGSERIAL PRIMARY KEY,
  "from_account_id" BIGSERIAL NOT NULL,
  "to_account_id" BIGSERIAL NOT NULL,
  "amount" FLOAT NOT NULL,
  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (NOW()) 
);

CREATE INDEX ON "accounts" ("owner_id");

CREATE INDEX ON "entries" ("account_id");

CREATE INDEX ON "transfers" ("from_account_id");

CREATE INDEX ON "transfers" ("to_account_id");

CREATE INDEX ON "transfers" ("from_account_id", "to_account_id");

COMMENT ON COLUMN "entries"."amount" IS 'Can be negative or positive';

COMMENT ON COLUMN "transfers"."amount" IS 'Must be positive';

ALTER TABLE "accounts" ADD CONSTRAINT "fk_owner" FOREIGN KEY ("owner_id") REFERENCES "owners" ("id") ON DELETE CASCADE;

ALTER TABLE "entries" ADD CONSTRAINT "fk_account" FOREIGN KEY ("account_id") REFERENCES "accounts" ("id") ON DELETE CASCADE;

ALTER TABLE "transfers" ADD CONSTRAINT "fk_from_account" FOREIGN KEY ("from_account_id") REFERENCES "accounts" ("id") ON DELETE CASCADE;

ALTER TABLE "transfers" ADD CONSTRAINT "fk_to_account" FOREIGN KEY ("to_account_id") REFERENCES "accounts" ("id") ON DELETE CASCADE;

INSERT INTO owners (first_name, first_surname, second_surname, nationality) VALUES ('Daniel', 'Soler', 'Fontanet', 34); 

INSERT INTO owners (first_name, first_surname, second_surname, nationality) VALUES ('Mar', 'Soler', 'Fontanet', 34); 

INSERT INTO owners (first_name, first_surname, second_surname, nationality) VALUES ('Roger', 'Metaute', 'Perez', 34); 