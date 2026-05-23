-- Create "todos" table
CREATE TABLE "public"."todos" (
  "id" bigserial NOT NULL,
  "created_at" timestamptz NULL,
  "updated_at" timestamptz NULL,
  "title" character varying(255) NOT NULL,
  "description" text NOT NULL,
  PRIMARY KEY ("id")
);
