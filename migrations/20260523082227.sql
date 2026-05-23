-- Create "refresh_tokens" table
CREATE TABLE "public"."refresh_tokens" (
  "id" bigserial NOT NULL,
  "user_id" bigint NOT NULL,
  "token" character varying(255) NOT NULL,
  "expires_at" timestamptz NOT NULL,
  "created_at" timestamptz NULL,
  PRIMARY KEY ("id")
);
-- Create index "idx_refresh_tokens_token" to table: "refresh_tokens"
CREATE UNIQUE INDEX "idx_refresh_tokens_token" ON "public"."refresh_tokens" ("token");
-- Create index "idx_refresh_tokens_user_id" to table: "refresh_tokens"
CREATE INDEX "idx_refresh_tokens_user_id" ON "public"."refresh_tokens" ("user_id");
