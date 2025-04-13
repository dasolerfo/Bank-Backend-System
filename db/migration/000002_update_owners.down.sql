ALTER TABLE "owners" ADD hashed_password VARCHAR NOT NULL;
ALTER TABLE "owners" ADD email VARCHAR NOT NULL UNIQUE;
ALTER TABLE "owners" ADD "created_at" TIMESTAMPTZ NOT NULL DEFAULT (NOW()) ;
ALTER TABLE "owners" ADD "password_changed_at" TIMESTAMPTZ NOT NULL DEFAULT '0001-01-01 00:00:00Z' ;
-- CREATE UNIQUE INDEX ON "accounts" ("owner_id", "currency")
ALTER TABLE "accounts" ADD CONSTRAINT "owner_currency_same" UNIQUE ("owner", "currency"); 