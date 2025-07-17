
CREATE TABLE "sessions" (
  "id" uuid PRIMARY KEY,
  "owner_id" BIGSERIAL NOT NULL,
  "email" varchar NOT NULL,
  "refresh_token" varchar NOT NULL,
  "user_agent" varchar NOT NULL,
  "client_ip" varchar NOT NULL,
  "is_blocked" boolean NOT NULL DEFAULT false,


  "created_at" TIMESTAMPTZ NOT NULL DEFAULT (NOW()),
  "expires_at" TIMESTAMPTZ NOT NULL  

);

alter table "sessions" add FOREIGN KEY ("owner_id") references "owners" ("id") on delete cascade;
